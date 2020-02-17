package template4app

var (
	SimpleJson = `
package simplejson

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
)

// returns the current implementation version
func Version() string {
	return "0.5.0"
}

type Json struct {
	data interface{}
}

func (j *Json) FromDB(data []byte) error {
	j.data = make(map[string]interface{})

	dec := json.NewDecoder(bytes.NewBuffer(data))
	dec.UseNumber()
	return dec.Decode(&j.data)
}

func (j *Json) ToDB() ([]byte, error) {
	if j == nil || j.data == nil {
		return nil, nil
	}

	return j.Encode()
}

// NewJson returns a pointer to a new `+ "`Json`" + ` object
// after unmarshaling `+ `body` +` bytes
func NewJson(body []byte) (*Json, error) {
	j := new(Json)
	err := j.UnmarshalJSON(body)
	if err != nil {
		return nil, err
	}
	return j, nil
}

// New returns a pointer to a new, empty `+"`json`"+` object
func New() *Json {
	return &Json{
		data: make(map[string]interface{}),
	}
}

// New returns a pointer to a new, empty `+ "`json`"+` object
func NewFromAny(data interface{}) *Json {
	return &Json{data: data}
}

// Interface returns the underlying data
func (j *Json) Interface() interface{} {
	return j.data
}

// Encode returns its marshaled data as ` + `[]byte`+ `
func (j *Json) Encode() ([]byte, error) {
	return j.MarshalJSON()
}

// EncodePretty returns its marshaled data as `+ `[]byte`+ ` with indentation
func (j *Json) EncodePretty() ([]byte, error) {
	return json.MarshalIndent(&j.data, "", "  ")
}

// Implements the json.Marshaler interface.
func (j *Json) MarshalJSON() ([]byte, error) {
	return json.Marshal(&j.data)
}

// Set modifies `+"`Json` map by `key` and `value`"+`
// Useful for changing single key/value in a `+"`json`"+` object easily.
func (j *Json) Set(key string, val interface{}) {
	m, err := j.Map()
	if err != nil {
		return
	}
	m[key] = val
}

// SetPath modifies `+"`json`"+`, recursively checking/creating map keys for the supplied path,
// and then finally writing in the value
func (j *Json) SetPath(branch []string, val interface{}) {
	if len(branch) == 0 {
		j.data = val
		return
	}

	// in order to insert our branch, we need map[string]interface{}
	if _, ok := (j.data).(map[string]interface{}); !ok {
		// have to replace with something suitable
		j.data = make(map[string]interface{})
	}
	curr := j.data.(map[string]interface{})

	for i := 0; i < len(branch)-1; i++ {
		b := branch[i]
		// key exists?
		if _, ok := curr[b]; !ok {
			n := make(map[string]interface{})
			curr[b] = n
			curr = n
			continue
		}

		// make sure the value is the right sort of thing
		if _, ok := curr[b].(map[string]interface{}); !ok {
			// have to replace with something suitable
			n := make(map[string]interface{})
			curr[b] = n
		}

		curr = curr[b].(map[string]interface{})
	}

	// add remaining k/v
	curr[branch[len(branch)-1]] = val
}

// Del modifies `+ "`json`"+` map by deleting `+ "`key`"+` if it is present.
func (j *Json) Del(key string) {
	m, err := j.Map()
	if err != nil {
		return
	}
	delete(m, key)
}

// Get returns a pointer to a new `+ "`json`"+` object
// for `+ "`key`"+` in its `+ "`map`"+` representation
//
// useful for chaining operations (to traverse a nested JSON):
//    js.Get("top_level").Get("dict").Get("value").Int()
func (j *Json) Get(key string) *Json {
	m, err := j.Map()
	if err == nil {
		if val, ok := m[key]; ok {
			return &Json{val}
		}
	}
	return &Json{nil}
}

// GetPath searches for the item as specified by the branch
// without the need to deep dive using Get()'s.
//
//   js.GetPath("top_level", "dict")
func (j *Json) GetPath(branch ...string) *Json {
	jin := j
	for _, p := range branch {
		jin = jin.Get(p)
	}
	return jin
}

//
// this is the analog to Get when accessing elements of
// a json array instead of a json object:
//    js.Get("top_level").Get("array").GetIndex(1).Get("key").Int()
func (j *Json) GetIndex(index int) *Json {
	a, err := j.Array()
	if err == nil {
		if len(a) > index {
			return &Json{a[index]}
		}
	}
	return &Json{nil}
}


//
// useful for chained operations when success is important:
//    if data, ok := js.Get("top_level").CheckGet("inner"); ok {
//        log.Println(data)
//    }
func (j *Json) CheckGet(key string) (*Json, bool) {
	m, err := j.Map()
	if err == nil {
		if val, ok := m[key]; ok {
			return &Json{val}, true
		}
	}
	return nil, false
}

// Map type asserts to `+"`map`"+`
func (j *Json) Map() (map[string]interface{}, error) {
	if m, ok := (j.data).(map[string]interface{}); ok {
		return m, nil
	}
	return nil, errors.New("type assertion to map[string]interface{} failed")
}

func (j *Json) Array() ([]interface{}, error) {
	if a, ok := (j.data).([]interface{}); ok {
		return a, nil
	}
	return nil, errors.New("type assertion to []interface{} failed")
}


func (j *Json) Bool() (bool, error) {
	if s, ok := (j.data).(bool); ok {
		return s, nil
	}
	return false, errors.New("type assertion to bool failed")
}

// String type asserts to `+"`string`"+`
func (j *Json) String() (string, error) {
	if s, ok := (j.data).(string); ok {
		return s, nil
	}
	return "", errors.New("type assertion to string failed")
}

// Bytes type asserts to `+"`[]byte`"+`
func (j *Json) Bytes() ([]byte, error) {
	if s, ok := (j.data).(string); ok {
		return []byte(s), nil
	}
	return nil, errors.New("type assertion to []byte failed")
}

// StringArray type asserts to an `+"`array`"+` of `+"`string`"+`
func (j *Json) StringArray() ([]string, error) {
	arr, err := j.Array()
	if err != nil {
		return nil, err
	}
	retArr := make([]string, 0, len(arr))
	for _, a := range arr {
		if a == nil {
			retArr = append(retArr, "")
			continue
		}
		s, ok := a.(string)
		if !ok {
			return nil, err
		}
		retArr = append(retArr, s)
	}
	return retArr, nil
}

// MustArray guarantees the return of a `+"`[]interface{}`"+` (with optional default)
//
// useful when you want to iterate over array values in a succinct manner:
//		for i, v := range js.Get("results").MustArray() {
//			fmt.Println(i, v)
//		}
func (j *Json) MustArray(args ...[]interface{}) []interface{} {
	var def []interface{}

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("MustArray() received too many arguments %d", len(args))
	}

	a, err := j.Array()
	if err == nil {
		return a
	}

	return def
}

// MustMap guarantees the return of a `+"`map[string]interface{}`"+` (with optional default)
//
// useful when you want to iterate over map values in a succinct manner:
//		for k, v := range js.Get("dictionary").MustMap() {
//			fmt.Println(k, v)
//		}
func (j *Json) MustMap(args ...map[string]interface{}) map[string]interface{} {
	var def map[string]interface{}

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("MustMap() received too many arguments %d", len(args))
	}

	a, err := j.Map()
	if err == nil {
		return a
	}

	return def
}


//
// useful when you explicitly want a `+"`string`"+` in a single value return context:
//     myFunc(js.Get("param1").MustString(), js.Get("optional_param").MustString("my_default"))
func (j *Json) MustString(args ...string) string {
	var def string

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("MustString() received too many arguments %d", len(args))
	}

	s, err := j.String()
	if err == nil {
		return s
	}

	return def
}

// MustStringArray guarantees the return of a `+"`[]string`"+` (with optional default)
//
// useful when you want to iterate over array values in a succinct manner:
//		for i, s := range js.Get("results").MustStringArray() {
//			fmt.Println(i, s)
//		}
func (j *Json) MustStringArray(args ...[]string) []string {
	var def []string

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("MustStringArray() received too many arguments %d", len(args))
	}

	a, err := j.StringArray()
	if err == nil {
		return a
	}

	return def
}


//     myFunc(js.Get("param1").MustInt(), js.Get("optional_param").MustInt(5150))
func (j *Json) MustInt(args ...int) int {
	var def int

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("MustInt() received too many arguments %d", len(args))
	}

	i, err := j.Int()
	if err == nil {
		return i
	}

	return def
}


//     myFunc(js.Get("param1").MustFloat64(), js.Get("optional_param").MustFloat64(5.150))
func (j *Json) MustFloat64(args ...float64) float64 {
	var def float64

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("MustFloat64() received too many arguments %d", len(args))
	}

	f, err := j.Float64()
	if err == nil {
		return f
	}

	return def
}


//     myFunc(js.Get("param1").MustBool(), js.Get("optional_param").MustBool(true))
func (j *Json) MustBool(args ...bool) bool {
	var def bool

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("MustBool() received too many arguments %d", len(args))
	}

	b, err := j.Bool()
	if err == nil {
		return b
	}

	return def
}


//     myFunc(js.Get("param1").MustInt64(), js.Get("optional_param").MustInt64(5150))
func (j *Json) MustInt64(args ...int64) int64 {
	var def int64

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("MustInt64() received too many arguments %d", len(args))
	}

	i, err := j.Int64()
	if err == nil {
		return i
	}

	return def
}


//     myFunc(js.Get("param1").MustUint64(), js.Get("optional_param").MustUint64(5150))
func (j *Json) MustUint64(args ...uint64) uint64 {
	var def uint64

	switch len(args) {
	case 0:
	case 1:
		def = args[0]
	default:
		log.Panicf("MustUint64() received too many arguments %d", len(args))
	}

	i, err := j.Uint64()
	if err == nil {
		return i
	}

	return def
}

`
	SimpleJsonGo11 = `
// +build go1.1

package simplejson

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"reflect"
	"strconv"
)

// Implements the json.Unmarshaler interface.
func (j *Json) UnmarshalJSON(p []byte) error {
	dec := json.NewDecoder(bytes.NewBuffer(p))
	dec.UseNumber()
	return dec.Decode(&j.data)
}

// NewFromReader returns a *Json by decoding from an io.Reader
func NewFromReader(r io.Reader) (*Json, error) {
	j := new(Json)
	dec := json.NewDecoder(r)
	dec.UseNumber()
	err := dec.Decode(&j.data)
	return j, err
}

// Float64 coerces into a float64
func (j *Json) Float64() (float64, error) {
	switch j.data.(type) {
	case json.Number:
		return j.data.(json.Number).Float64()
	case float32, float64:
		return reflect.ValueOf(j.data).Float(), nil
	case int, int8, int16, int32, int64:
		return float64(reflect.ValueOf(j.data).Int()), nil
	case uint, uint8, uint16, uint32, uint64:
		return float64(reflect.ValueOf(j.data).Uint()), nil
	}
	return 0, errors.New("invalid value type")
}

// Int coerces into an int
func (j *Json) Int() (int, error) {
	switch j.data.(type) {
	case json.Number:
		i, err := j.data.(json.Number).Int64()
		return int(i), err
	case float32, float64:
		return int(reflect.ValueOf(j.data).Float()), nil
	case int, int8, int16, int32, int64:
		return int(reflect.ValueOf(j.data).Int()), nil
	case uint, uint8, uint16, uint32, uint64:
		return int(reflect.ValueOf(j.data).Uint()), nil
	}
	return 0, errors.New("invalid value type")
}

// Int64 coerces into an int64
func (j *Json) Int64() (int64, error) {
	switch j.data.(type) {
	case json.Number:
		return j.data.(json.Number).Int64()
	case float32, float64:
		return int64(reflect.ValueOf(j.data).Float()), nil
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(j.data).Int(), nil
	case uint, uint8, uint16, uint32, uint64:
		return int64(reflect.ValueOf(j.data).Uint()), nil
	}
	return 0, errors.New("invalid value type")
}

// Uint64 coerces into an uint64
func (j *Json) Uint64() (uint64, error) {
	switch j.data.(type) {
	case json.Number:
		return strconv.ParseUint(j.data.(json.Number).String(), 10, 64)
	case float32, float64:
		return uint64(reflect.ValueOf(j.data).Float()), nil
	case int, int8, int16, int32, int64:
		return uint64(reflect.ValueOf(j.data).Int()), nil
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(j.data).Uint(), nil
	}
	return 0, errors.New("invalid value type")
}

`
	SimpleJsonTest = `
package simplejson

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimplejson(t *testing.T) {
	var ok bool
	var err error

	js, err := NewJson([]byte(` +"`"+`{
"test": {
"string_array": ["asdf", "ghjk", "zxcv"],
"string_array_null": ["abc", null, "efg"],
"array": [1, "2", 3],
"arraywithsubs": [{"subkeyone": 1},
{"subkeytwo": 2, "subkeythree": 3}],
"int": 10,
"float": 5.150,
"string": "simplejson",
"bool": true,
"sub_obj": {"a": 1}
}
}` +"`"+`))

	assert.NotEqual(t, nil, js)
	assert.Equal(t, nil, err)

	_, ok = js.CheckGet("test")
	assert.Equal(t, true, ok)

	_, ok = js.CheckGet("missing_key")
	assert.Equal(t, false, ok)

	aws := js.Get("test").Get("arraywithsubs")
	assert.NotEqual(t, nil, aws)
	var awsval int
	awsval, _ = aws.GetIndex(0).Get("subkeyone").Int()
	assert.Equal(t, 1, awsval)
	awsval, _ = aws.GetIndex(1).Get("subkeytwo").Int()
	assert.Equal(t, 2, awsval)
	awsval, _ = aws.GetIndex(1).Get("subkeythree").Int()
	assert.Equal(t, 3, awsval)

	i, _ := js.Get("test").Get("int").Int()
	assert.Equal(t, 10, i)

	f, _ := js.Get("test").Get("float").Float64()
	assert.Equal(t, 5.150, f)

	s, _ := js.Get("test").Get("string").String()
	assert.Equal(t, "simplejson", s)

	b, _ := js.Get("test").Get("bool").Bool()
	assert.Equal(t, true, b)

	mi := js.Get("test").Get("int").MustInt()
	assert.Equal(t, 10, mi)

	mi2 := js.Get("test").Get("missing_int").MustInt(5150)
	assert.Equal(t, 5150, mi2)

	ms := js.Get("test").Get("string").MustString()
	assert.Equal(t, "simplejson", ms)

	ms2 := js.Get("test").Get("missing_string").MustString("fyea")
	assert.Equal(t, "fyea", ms2)

	ma2 := js.Get("test").Get("missing_array").MustArray([]interface{}{"1", 2, "3"})
	assert.Equal(t, ma2, []interface{}{"1", 2, "3"})

	msa := js.Get("test").Get("string_array").MustStringArray()
	assert.Equal(t, msa[0], "asdf")
	assert.Equal(t, msa[1], "ghjk")
	assert.Equal(t, msa[2], "zxcv")

	msa2 := js.Get("test").Get("string_array").MustStringArray([]string{"1", "2", "3"})
	assert.Equal(t, msa2[0], "asdf")
	assert.Equal(t, msa2[1], "ghjk")
	assert.Equal(t, msa2[2], "zxcv")

	msa3 := js.Get("test").Get("missing_array").MustStringArray([]string{"1", "2", "3"})
	assert.Equal(t, msa3, []string{"1", "2", "3"})

	mm2 := js.Get("test").Get("missing_map").MustMap(map[string]interface{}{"found": false})
	assert.Equal(t, mm2, map[string]interface{}{"found": false})

	strs, err := js.Get("test").Get("string_array").StringArray()
	assert.Equal(t, err, nil)
	assert.Equal(t, strs[0], "asdf")
	assert.Equal(t, strs[1], "ghjk")
	assert.Equal(t, strs[2], "zxcv")

	strs2, err := js.Get("test").Get("string_array_null").StringArray()
	assert.Equal(t, err, nil)
	assert.Equal(t, strs2[0], "abc")
	assert.Equal(t, strs2[1], "")
	assert.Equal(t, strs2[2], "efg")

	gp, _ := js.GetPath("test", "string").String()
	assert.Equal(t, "simplejson", gp)

	gp2, _ := js.GetPath("test", "int").Int()
	assert.Equal(t, 10, gp2)

	assert.Equal(t, js.Get("test").Get("bool").MustBool(), true)

	js.Set("float2", 300.0)
	assert.Equal(t, js.Get("float2").MustFloat64(), 300.0)

	js.Set("test2", "setTest")
	assert.Equal(t, "setTest", js.Get("test2").MustString())

	js.Del("test2")
	assert.NotEqual(t, "setTest", js.Get("test2").MustString())

	js.Get("test").Get("sub_obj").Set("a", 2)
	assert.Equal(t, 2, js.Get("test").Get("sub_obj").Get("a").MustInt())

	js.GetPath("test", "sub_obj").Set("a", 3)
	assert.Equal(t, 3, js.GetPath("test", "sub_obj", "a").MustInt())
}

func TestStdlibInterfaces(t *testing.T) {
	val := new(struct {
		Name   string ` +"`json:\"name\"`"+`
		Params *Json  ` +"`json:\"params\"`"+`
	})
	val2 := new(struct {
		Name   string ` +"`json:\"name\"`"+`
		Params *Json  ` +"`json:\"params\"`"+`
	})

	raw := ` +"`{\"name\":\"myobject\",\"params\":{\"string\":\"simplejson\"}}`"+`

	assert.Equal(t, nil, json.Unmarshal([]byte(raw), val))

	assert.Equal(t, "myobject", val.Name)
	assert.NotEqual(t, nil, val.Params.data)
	s, _ := val.Params.Get("string").String()
	assert.Equal(t, "simplejson", s)

	p, err := json.Marshal(val)
	assert.Equal(t, nil, err)
	assert.Equal(t, nil, json.Unmarshal(p, val2))
	assert.Equal(t, val, val2) // stable
}

func TestSet(t *testing.T) {
	js, err := NewJson([]byte(` +"`{}`"+`))
	assert.Equal(t, nil, err)

	js.Set("baz", "bing")

	s, err := js.GetPath("baz").String()
	assert.Equal(t, nil, err)
	assert.Equal(t, "bing", s)
}

func TestReplace(t *testing.T) {
	js, err := NewJson([]byte(` +"`{}`"+`))
	assert.Equal(t, nil, err)

	err = js.UnmarshalJSON([]byte(` +"`{\"baz\":\"bing\"}`"+`))
	assert.Equal(t, nil, err)

	s, err := js.GetPath("baz").String()
	assert.Equal(t, nil, err)
	assert.Equal(t, "bing", s)
}

func TestSetPath(t *testing.T) {
	js, err := NewJson([]byte(` +"`{}`"+`))
	assert.Equal(t, nil, err)

	js.SetPath([]string{"foo", "bar"}, "baz")

	s, err := js.GetPath("foo", "bar").String()
	assert.Equal(t, nil, err)
	assert.Equal(t, "baz", s)
}

func TestSetPathNoPath(t *testing.T) {
	js, err := NewJson([]byte(` +"`{\"some\":\"data\",\"some_number\":1.0,\"some_bool\":false}`"+`))
	assert.Equal(t, nil, err)

	f := js.GetPath("some_number").MustFloat64(99.0)
	assert.Equal(t, f, 1.0)

	js.SetPath([]string{}, map[string]interface{}{"foo": "bar"})

	s, err := js.GetPath("foo").String()
	assert.Equal(t, nil, err)
	assert.Equal(t, "bar", s)

	f = js.GetPath("some_number").MustFloat64(99.0)
	assert.Equal(t, f, 99.0)
}

func TestPathWillAugmentExisting(t *testing.T) {
	js, err := NewJson([]byte(` +"`{\"this\":{\"a\":\"aa\",\"b\":\"bb\",\"c\":\"cc\"}}`"+`))
	assert.Equal(t, nil, err)

	js.SetPath([]string{"this", "d"}, "dd")

	cases := []struct {
		path    []string
		outcome string
	}{
		{
			path:    []string{"this", "a"},
			outcome: "aa",
		},
		{
			path:    []string{"this", "b"},
			outcome: "bb",
		},
		{
			path:    []string{"this", "c"},
			outcome: "cc",
		},
		{
			path:    []string{"this", "d"},
			outcome: "dd",
		},
	}

	for _, tc := range cases {
		s, err := js.GetPath(tc.path...).String()
		assert.Equal(t, nil, err)
		assert.Equal(t, tc.outcome, s)
	}
}

func TestPathWillOverwriteExisting(t *testing.T) {
	// notice how "a" is 0.1 - but then we'll try to set at path a, foo
	js, err := NewJson([]byte(`+"`{\"this\":{\"a\":0.1,\"b\":\"bb\",\"c\":\"cc\"}}`"+`))
	assert.Equal(t, nil, err)

	js.SetPath([]string{"this", "a", "foo"}, "bar")

	s, err := js.GetPath("this", "a", "foo").String()
	assert.Equal(t, nil, err)
	assert.Equal(t, "bar", s)
}

`

)

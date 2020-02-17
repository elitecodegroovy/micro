package template4app

var (
	ServiceRemotecacheDatabaseStorage = `
package remotecache

import (
	"context"
	"time"

	"{{.Dir}}/pkg/infra/log"
	"{{.Dir}}/pkg/services/sqlstore"
)

var getTime = time.Now

const databaseCacheType = "database"

type databaseCache struct {
	SQLStore *sqlstore.SqlStore
	log      log.Logger
}

func newDatabaseCache(sqlstore *sqlstore.SqlStore) *databaseCache {
	dc := &databaseCache{
		SQLStore: sqlstore,
		log:      log.New("remotecache.database"),
	}

	return dc
}

func (dc *databaseCache) Run(ctx context.Context) error {
	ticker := time.NewTicker(time.Minute * 10)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			dc.internalRunGC()
		}
	}
}

func (dc *databaseCache) internalRunGC() {
	err := dc.SQLStore.WithDbSession(context.Background(), func(session *sqlstore.DBSession) error {
		now := getTime().Unix()
		sql := ` +"`DELETE FROM cache_data WHERE (? - created_at) >= expires AND expires <> 0`"+`

		_, err := session.Exec(sql, now)
		return err
	})

	if err != nil {
		dc.log.Error("failed to run garbage collect", "error", err)
	}
}

func (dc *databaseCache) Get(key string) (interface{}, error) {
	cacheHit := CacheData{}
	session := dc.SQLStore.NewSession()
	defer session.Close()

	exist, err := session.Where("cache_key= ?", key).Get(&cacheHit)

	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, ErrCacheItemNotFound
	}

	if cacheHit.Expires > 0 {
		existedButExpired := getTime().Unix()-cacheHit.CreatedAt >= cacheHit.Expires
		if existedButExpired {
			err = dc.Delete(key) //ignore this error since we will return ` +"`ErrCacheItemNotFound`"+` anyway
			if err != nil {
				dc.log.Debug("Deletion of expired key failed: %v", err)
			}
			return nil, ErrCacheItemNotFound
		}
	}

	item := &cachedItem{}
	if err = decodeGob(cacheHit.Data, item); err != nil {
		return nil, err
	}

	return item.Val, nil
}

func (dc *databaseCache) Set(key string, value interface{}, expire time.Duration) error {
	item := &cachedItem{Val: value}
	data, err := encodeGob(item)
	if err != nil {
		return err
	}

	session := dc.SQLStore.NewSession()
	defer session.Close()

	var expiresInSeconds int64
	if expire != 0 {
		expiresInSeconds = int64(expire) / int64(time.Second)
	}

	// attempt to insert the key
	sql := ` +"`INSERT INTO cache_data (cache_key,data,created_at,expires) VALUES(?,?,?,?)`"+`
	_, err = session.Exec(sql, key, data, getTime().Unix(), expiresInSeconds)
	if err != nil {
		// attempt to update if a unique constrain violation or a deadlock (for MySQL) occurs
		// if the update fails propagate the error
		// which eventually will result in a key that is not finally set
		// but since it's a cache does not harm a lot
		if dc.SQLStore.Dialect.IsUniqueConstraintViolation(err) || dc.SQLStore.Dialect.IsDeadlock(err) {
			sql := ` +"`UPDATE cache_data SET data=?, created_at=?, expires=? WHERE cache_key=?`"+ `
			_, err = session.Exec(sql, data, getTime().Unix(), expiresInSeconds, key)
			if err != nil && dc.SQLStore.Dialect.IsDeadlock(err) {
				// most probably somebody else is upserting the key
				// so it is safe enough not to propagate this error
				return nil
			}
		}
	}

	return err
}

func (dc *databaseCache) Delete(key string) error {
	return dc.SQLStore.WithDbSession(context.Background(), func(session *sqlstore.DBSession) error {
		sql := "DELETE FROM cache_data WHERE cache_key=?"
		_, err := session.Exec(sql, key)

		return err
	})

}

// CacheData is the struct representing the table in the database
type CacheData struct {
	CacheKey  string
	Data      []byte
	Expires   int64
	CreatedAt int64
}

`
	ServiceRemotecacheDatabaseStorageTest = `
package remotecache

import (
	"testing"
	"time"

	"{{.Dir}}/pkg/infra/log"
	"{{.Dir}}/pkg/services/sqlstore"
	"github.com/stretchr/testify/assert"
)

func TestDatabaseStorageGarbageCollection(t *testing.T) {
	sqlstore := sqlstore.InitTestDB(t)

	db := &databaseCache{
		SQLStore: sqlstore,
		log:      log.New("remotecache.database"),
	}

	obj := &CacheableStruct{String: "foolbar"}

	//set time.now to 2 weeks ago
	var err error
	getTime = func() time.Time { return time.Now().AddDate(0, 0, -2) }
	err = db.Set("key1", obj, 1000*time.Second)
	assert.Equal(t, err, nil)

	err = db.Set("key2", obj, 1000*time.Second)
	assert.Equal(t, err, nil)

	err = db.Set("key3", obj, 1000*time.Second)
	assert.Equal(t, err, nil)

	// insert object that should never expire
	db.Set("key4", obj, 0)

	getTime = time.Now
	db.Set("key5", obj, 1000*time.Second)

	//run GC
	db.internalRunGC()

	//try to read values
	_, err = db.Get("key1")
	assert.Equal(t, err, ErrCacheItemNotFound, "expected cache item not found. got: ", err)
	_, err = db.Get("key2")
	assert.Equal(t, err, ErrCacheItemNotFound)
	_, err = db.Get("key3")
	assert.Equal(t, err, ErrCacheItemNotFound)

	_, err = db.Get("key4")
	assert.Equal(t, err, nil)
	_, err = db.Get("key5")
	assert.Equal(t, err, nil)
}

func TestSecondSet(t *testing.T) {
	var err error
	sqlstore := sqlstore.InitTestDB(t)

	db := &databaseCache{
		SQLStore: sqlstore,
		log:      log.New("remotecache.database"),
	}

	obj := &CacheableStruct{String: "hey!"}

	err = db.Set("killa-gorilla", obj, 0)
	assert.Equal(t, err, nil)

	err = db.Set("killa-gorilla", obj, 0)
	assert.Equal(t, err, nil)
}

`
	ServiceRemotecacheMemcacheStorage = `
package remotecache

import (
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"{{.Dir}}/pkg/setting"
)

const memcachedCacheType = "memcached"

type memcachedStorage struct {
	c *memcache.Client
}

func newMemcachedStorage(opts *setting.RemoteCacheOptions) *memcachedStorage {
	return &memcachedStorage{
		c: memcache.New(opts.ConnStr),
	}
}

func newItem(sid string, data []byte, expire int32) *memcache.Item {
	return &memcache.Item{
		Key:        sid,
		Value:      data,
		Expiration: expire,
	}
}

// Set sets value to given key in the cache.
func (s *memcachedStorage) Set(key string, val interface{}, expires time.Duration) error {
	item := &cachedItem{Val: val}
	bytes, err := encodeGob(item)
	if err != nil {
		return err
	}

	var expiresInSeconds int64
	if expires != 0 {
		expiresInSeconds = int64(expires) / int64(time.Second)
	}

	memcachedItem := newItem(key, bytes, int32(expiresInSeconds))
	return s.c.Set(memcachedItem)
}

// Get gets value by given key in the cache.
func (s *memcachedStorage) Get(key string) (interface{}, error) {
	memcachedItem, err := s.c.Get(key)
	if err != nil && err.Error() == "memcache: cache miss" {
		return nil, ErrCacheItemNotFound
	}

	if err != nil {
		return nil, err
	}

	item := &cachedItem{}

	err = decodeGob(memcachedItem.Value, item)
	if err != nil {
		return nil, err
	}

	return item.Val, nil
}

// Delete delete a key from the cache
func (s *memcachedStorage) Delete(key string) error {
	return s.c.Delete(key)
}

`
	ServiceRemotecacheMemcacheStorageTest = `
// +build memcached

package remotecache

import (
	"testing"

	"{{.Dir}}/pkg/setting"
)

func TestMemcachedCacheStorage(t *testing.T) {
	opts := &setting.RemoteCacheOptions{Name: memcachedCacheType, ConnStr: "localhost:11211"}
	client := createTestClient(t, opts, nil)
	runTestsForClient(t, client)
}

`
	ServiceRemotecacheRedisStorage = `
package remotecache

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"{{.Dir}}/pkg/setting"
	"{{.Dir}}/pkg/util/errutil"
	redis "gopkg.in/redis.v2"
)

const redisCacheType = "redis"

type redisStorage struct {
	c *redis.Client
}

// parseRedisConnStr parses k=v pairs in csv and builds a redis Options object
func parseRedisConnStr(connStr string) (*redis.Options, error) {
	keyValueCSV := strings.Split(connStr, ",")
	options := &redis.Options{Network: "tcp"}
	for _, rawKeyValue := range keyValueCSV {
		keyValueTuple := strings.Split(rawKeyValue, "=")
		if len(keyValueTuple) != 2 {
			return nil, fmt.Errorf("incorrect redis connection string format detected for '%v', format is key=value,key=value", rawKeyValue)
		}
		connKey := keyValueTuple[0]
		connVal := keyValueTuple[1]
		switch connKey {
		case "addr":
			options.Addr = connVal
		case "password":
			options.Password = connVal
		case "db":
			i, err := strconv.ParseInt(connVal, 10, 64)
			if err != nil {
				return nil, errutil.Wrap("value for db in redis connection string must be a number", err)
			}
			options.DB = i
		case "pool_size":
			i, err := strconv.Atoi(connVal)
			if err != nil {
				return nil, errutil.Wrap("value for pool_size in redis connection string must be a number", err)
			}
			options.PoolSize = i
		default:
			return nil, fmt.Errorf("unrecorgnized option '%v' in redis connection string", connVal)
		}
	}
	return options, nil
}

func newRedisStorage(opts *setting.RemoteCacheOptions) (*redisStorage, error) {
	opt, err := parseRedisConnStr(opts.ConnStr)
	if err != nil {
		return nil, err
	}
	return &redisStorage{c: redis.NewClient(opt)}, nil
}

// Set sets value to given key in session.
func (s *redisStorage) Set(key string, val interface{}, expires time.Duration) error {
	item := &cachedItem{Val: val}
	value, err := encodeGob(item)
	if err != nil {
		return err
	}
	status := s.c.SetEx(key, expires, string(value))
	return status.Err()
}

// Get gets value by given key in session.
func (s *redisStorage) Get(key string) (interface{}, error) {
	v := s.c.Get(key)

	item := &cachedItem{}
	err := decodeGob([]byte(v.Val()), item)

	if err == nil {
		return item.Val, nil
	}
	if err.Error() == "EOF" {
		return nil, ErrCacheItemNotFound
	}
	return nil, err
}

// Delete delete a key from session.
func (s *redisStorage) Delete(key string) error {
	cmd := s.c.Del(key)
	return cmd.Err()
}

`
	ServiceRemotecacheRedisStorageIntegrationTest = `
// +build redis

package remotecache

import (
	"testing"

	"{{.Dir}}/pkg/setting"
)

func TestRedisCacheStorage(t *testing.T) {

	opts := &setting.RemoteCacheOptions{Name: redisCacheType, ConnStr: "addr=localhost:6379"}
	client := createTestClient(t, opts, nil)
	runTestsForClient(t, client)
}

`
	ServiceRemotecacheRedisStorageTest = `
package remotecache

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	redis "gopkg.in/redis.v2"
)

func Test_parseRedisConnStr(t *testing.T) {
	cases := map[string]struct {
		InputConnStr  string
		OutputOptions *redis.Options
		ShouldErr     bool
	}{
		"all redis options should parse": {
			"addr=127.0.0.1:6379,pool_size=100,db=1,password=grafanaRocks",
			&redis.Options{
				Addr:     "127.0.0.1:6379",
				PoolSize: 100,
				DB:       1,
				Password: "grafanaRocks",
				Network:  "tcp",
			},
			false,
		},
		"subset of redis options should parse": {
			"addr=127.0.0.1:6379,pool_size=100",
			&redis.Options{
				Addr:     "127.0.0.1:6379",
				PoolSize: 100,
				Network:  "tcp",
			},
			false,
		},
		"trailing comma should err": {
			"addr=127.0.0.1:6379,pool_size=100,",
			nil,
			true,
		},
		"invalid key should err": {
			"addr=127.0.0.1:6379,puddle_size=100",
			nil,
			true,
		},
		"empty connection string should err": {
			"",
			nil,
			true,
		},
	}

	for reason, testCase := range cases {
		options, err := parseRedisConnStr(testCase.InputConnStr)
		if testCase.ShouldErr {
			assert.Error(t, err, fmt.Sprintf("error cases should return non-nil error for test case %v", reason))
			assert.Nil(t, options, fmt.Sprintf("error cases should return nil for redis options for test case %v", reason))
			continue
		}
		assert.NoError(t, err, reason)
		assert.EqualValues(t, testCase.OutputOptions, options, reason)

	}
}

`
	ServiceRemotecache = `
package remotecache

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"time"

	"{{.Dir}}/pkg/infra/log"
	"{{.Dir}}/pkg/registry"
	"{{.Dir}}/pkg/services/sqlstore"
	"{{.Dir}}/pkg/setting"
)

var (
	// ErrCacheItemNotFound is returned if cache does not exist
	ErrCacheItemNotFound = errors.New("cache item not found")

	// ErrInvalidCacheType is returned if the type is invalid
	ErrInvalidCacheType = errors.New("invalid remote cache name")

	defaultMaxCacheExpiration = time.Hour * 24
)

func init() {
	registry.RegisterService(&RemoteCache{})
}

// CacheStorage allows the caller to set, get and delete items in the cache.
// Cached items are stored as byte arrays and marshalled using "encoding/gob"
// so any struct added to the cache needs to be registred with ` +"`remotecache.Register`"+`
// ex ` +"`remotecache.Register(CacheableStruct{})`"+`
type CacheStorage interface {
	// Get reads object from Cache
	Get(key string) (interface{}, error)

	// Set sets an object into the cache. if ` +"`expire`"+` is set to zero it will default to 24h
	Set(key string, value interface{}, expire time.Duration) error

	// Delete object from cache
	Delete(key string) error
}

// RemoteCache allows Grafana to cache data outside its own process
type RemoteCache struct {
	log      log.Logger
	client   CacheStorage
	SQLStore *sqlstore.SqlStore ` +"`inject:\"\"`"+`
	Cfg      *setting.Cfg       ` +"`inject:\"\"`"+`
}

// Get reads object from Cache
func (ds *RemoteCache) Get(key string) (interface{}, error) {
	return ds.client.Get(key)
}

// Set sets an object into the cache. if ` + "`expire`"+ ` is set to zero it will default to 24h
func (ds *RemoteCache) Set(key string, value interface{}, expire time.Duration) error {
	if expire == 0 {
		expire = defaultMaxCacheExpiration
	}

	return ds.client.Set(key, value, expire)
}

// Delete object from cache
func (ds *RemoteCache) Delete(key string) error {
	return ds.client.Delete(key)
}

// Init initializes the service
func (ds *RemoteCache) Init() error {
	ds.log = log.New("cache.remote")
	var err error
	ds.client, err = createClient(ds.Cfg.RemoteCacheOptions, ds.SQLStore)
	return err
}

// Run start the backend processes for cache clients
func (ds *RemoteCache) Run(ctx context.Context) error {
	//create new interface if more clients need GC jobs
	backgroundjob, ok := ds.client.(registry.BackgroundService)
	if ok {
		return backgroundjob.Run(ctx)
	}

	<-ctx.Done()
	return ctx.Err()
}

func createClient(opts *setting.RemoteCacheOptions, sqlstore *sqlstore.SqlStore) (CacheStorage, error) {
	if opts.Name == redisCacheType {
		return newRedisStorage(opts)
	}

	if opts.Name == memcachedCacheType {
		return newMemcachedStorage(opts), nil
	}

	if opts.Name == databaseCacheType {
		return newDatabaseCache(sqlstore), nil
	}

	return nil, ErrInvalidCacheType
}

// Register records a type, identified by a value for that type, under its
// internal type name. That name will identify the concrete type of a value
// sent or received as an interface variable. Only types that will be
// transferred as implementations of interface values need to be registered.
// Expecting to be used only during initialization, it panics if the mapping
// between types and names is not a bijection.
func Register(value interface{}) {
	gob.Register(value)
}

type cachedItem struct {
	Val interface{}
}

func encodeGob(item *cachedItem) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buf).Encode(item)
	return buf.Bytes(), err
}

func decodeGob(data []byte, out *cachedItem) error {
	buf := bytes.NewBuffer(data)
	return gob.NewDecoder(buf).Decode(&out)
}

`
	ServiceRemotecacheTest = `
package remotecache

import (
	"testing"
	"time"

	"{{.Dir}}/pkg/services/sqlstore"
	"{{.Dir}}/pkg/setting"
	"github.com/stretchr/testify/assert"
)

type CacheableStruct struct {
	String string
	Int64  int64
}

func init() {
	Register(CacheableStruct{})
}

func createTestClient(t *testing.T, opts *setting.RemoteCacheOptions, sqlstore *sqlstore.SqlStore) CacheStorage {
	t.Helper()

	dc := &RemoteCache{
		SQLStore: sqlstore,
		Cfg: &setting.Cfg{
			RemoteCacheOptions: opts,
		},
	}

	err := dc.Init()
	if err != nil {
		t.Fatalf("failed to init client for test. error: %v", err)
	}

	return dc
}

func TestCachedBasedOnConfig(t *testing.T) {

	cfg := setting.NewCfg()
	cfg.Load(&setting.CommandLineArgs{
		HomePath: "../../../",
	})

	client := createTestClient(t, cfg.RemoteCacheOptions, sqlstore.InitTestDB(t))
	runTestsForClient(t, client)
}

func TestInvalidCacheTypeReturnsError(t *testing.T) {
	_, err := createClient(&setting.RemoteCacheOptions{Name: "invalid"}, nil)
	assert.Equal(t, err, ErrInvalidCacheType)
}

func runTestsForClient(t *testing.T, client CacheStorage) {
	canPutGetAndDeleteCachedObjects(t, client)
	canNotFetchExpiredItems(t, client)
}

func canPutGetAndDeleteCachedObjects(t *testing.T, client CacheStorage) {
	cacheableStruct := CacheableStruct{String: "hej", Int64: 2000}

	err := client.Set("key1", cacheableStruct, 0)
	assert.Equal(t, err, nil, "expected nil. got: ", err)

	data, err := client.Get("key1")
	assert.Equal(t, err, nil)
	s, ok := data.(CacheableStruct)

	assert.Equal(t, ok, true)
	assert.Equal(t, s.String, "hej")
	assert.Equal(t, s.Int64, int64(2000))

	err = client.Delete("key1")
	assert.Equal(t, err, nil)

	_, err = client.Get("key1")
	assert.Equal(t, err, ErrCacheItemNotFound)
}

func canNotFetchExpiredItems(t *testing.T, client CacheStorage) {
	cacheableStruct := CacheableStruct{String: "hej", Int64: 2000}

	err := client.Set("key1", cacheableStruct, time.Second)
	assert.Equal(t, err, nil)

	//not sure how this can be avoided when testing redis/memcached :/
	<-time.After(time.Second + time.Millisecond)

	// should not be able to read that value since its expired
	_, err = client.Get("key1")
	assert.Equal(t, err, ErrCacheItemNotFound)
}

`
	ServiceRemotecacheTesting = `
package remotecache

import (
	"testing"

	"{{.Dir}}/pkg/services/sqlstore"
	"{{.Dir}}/pkg/setting"
)

// NewFakeStore creates store for testing
func NewFakeStore(t *testing.T) *RemoteCache {
	t.Helper()

	opts := &setting.RemoteCacheOptions{
		Name:    "database",
		ConnStr: "",
	}

	SQLStore := sqlstore.InitTestDB(t)

	dc := &RemoteCache{
		SQLStore: SQLStore,
		Cfg: &setting.Cfg{
			RemoteCacheOptions: opts,
		},
	}

	err := dc.Init()
	if err != nil {
		t.Fatalf("failed to init remote cache for test. error: %v", err)
	}

	return dc
}

`

)

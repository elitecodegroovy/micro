package template4app


var (
	Apikeygen = `
package apikeygen

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	"{{.Dir}}/pkg/util"
)

var ErrInvalidApiKey = errors.New("Invalid Api Key")

type KeyGenResult struct {
	HashedKey    string
	ClientSecret string
}

type ApiKeyJson struct {
	Key   string ` + "`json:\"k\"`" + `
	Name  string ` + "`json:\"n\"`" + `
	OrgId int64  ` + "`json:\"id\"`" + `
}

func New(orgId int64, name string) KeyGenResult {
	jsonKey := ApiKeyJson{}

	jsonKey.OrgId = orgId
	jsonKey.Name = name
	jsonKey.Key = util.GetRandomString(32)

	result := KeyGenResult{}
	result.HashedKey = util.EncodePassword(jsonKey.Key, name)

	jsonString, _ := json.Marshal(jsonKey)

	result.ClientSecret = base64.StdEncoding.EncodeToString(jsonString)
	return result
}

func Decode(keyString string) (*ApiKeyJson, error) {
	jsonString, err := base64.StdEncoding.DecodeString(keyString)
	if err != nil {
		return nil, ErrInvalidApiKey
	}

	var keyObj ApiKeyJson
	err = json.Unmarshal(jsonString, &keyObj)
	if err != nil {
		return nil, ErrInvalidApiKey
	}

	return &keyObj, nil
}

func IsValid(key *ApiKeyJson, hashedKey string) bool {
	check := util.EncodePassword(key.Key, key.Name)
	return check == hashedKey
}

`
	ApikeygenTest = `
package apikeygen

import (
	"testing"

	"{{.Dir}}/pkg/util"
	. "github.com/smartystreets/goconvey/convey"
)

func TestApiKeyGen(t *testing.T) {

	Convey("When generating new api key", t, func() {
		result := New(12, "Cool key")

		So(result.ClientSecret, ShouldNotBeEmpty)
		So(result.HashedKey, ShouldNotBeEmpty)

		Convey("can decode key", func() {
			keyInfo, err := Decode(result.ClientSecret)
			So(err, ShouldBeNil)

			keyHashed := util.EncodePassword(keyInfo.Key, keyInfo.Name)
			So(keyHashed, ShouldEqual, result.HashedKey)
		})
	})
}

`
)

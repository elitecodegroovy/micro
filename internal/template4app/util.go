package template4app

var (
	UtilErrors = `
package errutil

import (
	"fmt"

	"golang.org/x/xerrors"
)

// Wrap is a simple wrapper around Errorf that is doing error wrapping. You can read how that works in
// https://godoc.org/golang.org/x/xerrors#Errorf but its API is very implicit which is a reason for this wrapper.
// There is also a discussion (https://github.com/golang/go/issues/29934) where many comments make arguments for such
// wrapper so hopefully it will be added in the standard lib later.
func Wrap(message string, err error) error {
	if err == nil {
		return nil
	}

	return xerrors.Errorf("%v: %w", message, err)
}

// Wrapf is a simple wrapper around Errorf that is doing error wrapping
// Wrapf allows you to send a format and args instead of just a message.
func Wrapf(err error, message string, a ...interface{}) error {
	if err == nil {
		return nil
	}

	return Wrap(fmt.Sprintf(message, a...), err)
}

`

	UtilEncoding = `
package util

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"hash"
	"strings"
)

// GetRandomString generate random string by specify chars.
// source: https://github.com/gogits/gogs/blob/9ee80e3e5426821f03a4e99fad34418f5c736413/modules/base/tool.go#L58
func GetRandomString(n int, alphabets ...byte) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		if len(alphabets) == 0 {
			bytes[i] = alphanum[b%byte(len(alphanum))]
		} else {
			bytes[i] = alphabets[b%byte(len(alphabets))]
		}
	}
	return string(bytes)
}

// EncodePassword encodes a password using PBKDF2.
func EncodePassword(password string, salt string) string {
	newPasswd := PBKDF2([]byte(password), []byte(salt), 10000, 50, sha256.New)
	return hex.EncodeToString(newPasswd)
}

// EncodeMd5 encodes a string to md5 hex value.
func EncodeMd5(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

// PBKDF2 implements Password-Based Key Derivation Function 2), aimed to reduce
// the vulnerability of encrypted keys to brute force attacks.
// http://code.google.com/p/go/source/browse/pbkdf2/pbkdf2.go?repo=crypto
func PBKDF2(password, salt []byte, iter, keyLen int, h func() hash.Hash) []byte {
	prf := hmac.New(h, password)
	hashLen := prf.Size()
	numBlocks := (keyLen + hashLen - 1) / hashLen

	var buf [4]byte
	dk := make([]byte, 0, numBlocks*hashLen)
	U := make([]byte, hashLen)
	for block := 1; block <= numBlocks; block++ {
		// N.B.: || means concatenation, ^ means XOR
		// for each block T_i = U_1 ^ U_2 ^ ... ^ U_iter
		// U_1 = PRF(password, salt || uint(i))
		prf.Reset()
		prf.Write(salt)
		buf[0] = byte(block >> 24)
		buf[1] = byte(block >> 16)
		buf[2] = byte(block >> 8)
		buf[3] = byte(block)
		prf.Write(buf[:4])
		dk = prf.Sum(dk)
		T := dk[len(dk)-hashLen:]
		copy(U, T)

		// U_n = PRF(password, U_(n-1))
		for n := 2; n <= iter; n++ {
			prf.Reset()
			prf.Write(U)
			U = U[:0]
			U = prf.Sum(U)
			for x := range U {
				T[x] ^= U[x]
			}
		}
	}
	return dk[:keyLen]
}

// GetBasicAuthHeader returns a base64 encoded string from user and password.
func GetBasicAuthHeader(user string, password string) string {
	var userAndPass = user + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(userAndPass))
}

// DecodeBasicAuthHeader decodes user and password from a basic auth header.
func DecodeBasicAuthHeader(header string) (string, string, error) {
	var code string
	parts := strings.SplitN(header, " ", 2)
	if len(parts) == 2 && parts[0] == "Basic" {
		code = parts[1]
	}

	decoded, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		return "", "", err
	}

	userAndPass := strings.SplitN(string(decoded), ":", 2)
	if len(userAndPass) != 2 {
		return "", "", errors.New("Invalid basic auth header")
	}

	return userAndPass[0], userAndPass[1], nil
}

// RandomHex returns a random string from a n seed.
func RandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

`
	UtilEncodingTest = `
package util

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEncoding(t *testing.T) {
	Convey("When generating base64 header", t, func() {
		result := GetBasicAuthHeader("grafana", "1234")

		So(result, ShouldEqual, "Basic Z3JhZmFuYToxMjM0")
	})

	Convey("When decoding basic auth header", t, func() {
		header := GetBasicAuthHeader("grafana", "1234")
		username, password, err := DecodeBasicAuthHeader(header)
		So(err, ShouldBeNil)

		So(username, ShouldEqual, "grafana")
		So(password, ShouldEqual, "1234")
	})

	Convey("When encoding password", t, func() {
		encodedPassword := EncodePassword("iamgod", "pepper")
		So(encodedPassword, ShouldEqual, "e59c568621e57756495a468f47c74e07c911b037084dd464bb2ed72410970dc849cabd71b48c394faf08a5405dae53741ce9")
	})
}

`
	UtilEncryption = `
package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
)

const saltLength = 8

// Decrypt decrypts a payload with a given secret.
func Decrypt(payload []byte, secret string) ([]byte, error) {
	salt := payload[:saltLength]
	key := encryptionKeyToBytes(secret, string(salt))

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(payload) < aes.BlockSize {
		return nil, errors.New("payload too short")
	}
	iv := payload[saltLength : saltLength+aes.BlockSize]
	payload = payload[saltLength+aes.BlockSize:]
	payloadDst := make([]byte, len(payload))

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(payloadDst, payload)
	return payloadDst, nil
}

// Encrypt encrypts a payload with a given secret.
func Encrypt(payload []byte, secret string) ([]byte, error) {
	salt := GetRandomString(saltLength)

	key := encryptionKeyToBytes(secret, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, saltLength+aes.BlockSize+len(payload))
	copy(ciphertext[:saltLength], []byte(salt))
	iv := ciphertext[saltLength : saltLength+aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[saltLength+aes.BlockSize:], payload)

	return ciphertext, nil
}

// Key needs to be 32bytes
func encryptionKeyToBytes(secret, salt string) []byte {
	return PBKDF2([]byte(secret), []byte(salt), 10000, 32, sha256.New)
}

`
	UtilEncryptionTest = `
package util

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEncryption(t *testing.T) {

	Convey("When getting encryption key", t, func() {

		key := encryptionKeyToBytes("secret", "salt")
		So(len(key), ShouldEqual, 32)

		key = encryptionKeyToBytes("a very long secret key that is larger then 32bytes", "salt")
		So(len(key), ShouldEqual, 32)
	})

	Convey("When decrypting basic payload", t, func() {
		encrypted, encryptErr := Encrypt([]byte("grafana"), "1234")
		decrypted, decryptErr := Decrypt(encrypted, "1234")

		So(encryptErr, ShouldBeNil)
		So(decryptErr, ShouldBeNil)
		So(string(decrypted), ShouldEqual, "grafana")
	})

}

`
	UtilFilepath = `
package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

//ErrWalkSkipDir is the Error returned when we want to skip descending into a directory
var ErrWalkSkipDir = errors.New("skip this directory")

//WalkFunc is a callback function called for each path as a directory is walked
//If resolvedPath != "", then we are following symbolic links.
type WalkFunc func(resolvedPath string, info os.FileInfo, err error) error

//Walk walks a path, optionally following symbolic links, and for each path,
//it calls the walkFn passed.
//
//It is similar to filepath.Walk, except that it supports symbolic links and
//can detect infinite loops while following sym links.
//It solves the issue where your WalkFunc needs a path relative to the symbolic link
//(resolving links within walkfunc loses the path to the symbolic link for each traversal).
func Walk(path string, followSymlinks bool, detectSymlinkInfiniteLoop bool, walkFn WalkFunc) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	var symlinkPathsFollowed map[string]bool
	var resolvedPath string
	if followSymlinks {
		resolvedPath = path
		if detectSymlinkInfiniteLoop {
			symlinkPathsFollowed = make(map[string]bool, 8)
		}
	}
	return walk(path, info, resolvedPath, symlinkPathsFollowed, walkFn)
}

//walk walks the path. It is a helper/sibling function to Walk.
//It takes a resolvedPath into consideration. This way, paths being walked are
//always relative to the path argument, even if symbolic links were resolved).
//
//If resolvedPath is "", then we are not following symbolic links.
//If symlinkPathsFollowed is not nil, then we need to detect infinite loop.
func walk(path string, info os.FileInfo, resolvedPath string, symlinkPathsFollowed map[string]bool, walkFn WalkFunc) error {
	if info == nil {
		return errors.New("Walk: Nil FileInfo passed")
	}
	err := walkFn(resolvedPath, info, nil)
	if err != nil {
		if info.IsDir() && err == ErrWalkSkipDir {
			err = nil
		}
		return err
	}
	if resolvedPath != "" && info.Mode()&os.ModeSymlink == os.ModeSymlink {
		path2, err := os.Readlink(resolvedPath)
		if err != nil {
			return err
		}
		//vout("SymLink Path: %v, links to: %v", resolvedPath, path2)
		if symlinkPathsFollowed != nil {
			if _, ok := symlinkPathsFollowed[path2]; ok {
				errMsg := "Potential SymLink Infinite Loop. Path: %v, Link To: %v"
				return fmt.Errorf(errMsg, resolvedPath, path2)
			}
			symlinkPathsFollowed[path2] = true
		}
		info2, err := os.Lstat(path2)
		if err != nil {
			return err
		}
		return walk(path, info2, path2, symlinkPathsFollowed, walkFn)
	}
	if info.IsDir() {
		list, err := ioutil.ReadDir(path)
		if err != nil {
			return walkFn(resolvedPath, info, err)
		}
		var subFiles = make([]subFile, 0)
		for _, fileInfo := range list {
			path2 := filepath.Join(path, fileInfo.Name())
			var resolvedPath2 string
			if resolvedPath != "" {
				resolvedPath2 = filepath.Join(resolvedPath, fileInfo.Name())
			}
			subFiles = append(subFiles, subFile{path: path2, resolvedPath: resolvedPath2, fileInfo: fileInfo})
		}

		if containsDistFolder(subFiles) {
			err := walk(
				filepath.Join(path, "dist"),
				info,
				filepath.Join(resolvedPath, "dist"),
				symlinkPathsFollowed,
				walkFn)

			if err != nil {
				return err
			}
		} else {
			for _, p := range subFiles {
				err = walk(p.path, p.fileInfo, p.resolvedPath, symlinkPathsFollowed, walkFn)

				if err != nil {
					return err
				}
			}
		}

		return nil
	}
	return nil
}

type subFile struct {
	path, resolvedPath string
	fileInfo           os.FileInfo
}

func containsDistFolder(subFiles []subFile) bool {
	for _, p := range subFiles {
		if p.fileInfo.IsDir() && p.fileInfo.Name() == "dist" {
			return true
		}
	}

	return false
}

`
	UtilIpAddress = `
package util

import (
	"net"
	"strings"
)

// ParseIPAddress parses an IP address and removes port and/or IPV6 format
func ParseIPAddress(input string) string {
	host, _ := SplitHostPort(input)

	ip := net.ParseIP(host)

	if ip == nil {
		return host
	}

	if ip.IsLoopback() {
		return "127.0.0.1"
	}

	return ip.String()
}

// SplitHostPortDefault splits ip address/hostname string by host and port. Defaults used if no match found
func SplitHostPortDefault(input, defaultHost, defaultPort string) (host string, port string) {
	port = defaultPort
	s := input
	lastIndex := strings.LastIndex(input, ":")

	if lastIndex != -1 {
		if lastIndex > 0 && input[lastIndex-1:lastIndex] != ":" {
			s = input[:lastIndex]
			port = input[lastIndex+1:]
		} else if lastIndex == 0 {
			s = defaultHost
			port = input[lastIndex+1:]
		}
	} else {
		port = defaultPort
	}

	s = strings.Replace(s, "[", "", -1)
	s = strings.Replace(s, "]", "", -1)
	port = strings.Replace(port, "[", "", -1)
	port = strings.Replace(port, "]", "", -1)

	return s, port
}

// SplitHostPort splits ip address/hostname string by host and port
func SplitHostPort(input string) (host string, port string) {
	return SplitHostPortDefault(input, "", "")
}

`
	UtilIpAddressTest = `
package util

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseIPAddress(t *testing.T) {
	Convey("Test parse ip address", t, func() {
		So(ParseIPAddress("192.168.0.140:456"), ShouldEqual, "192.168.0.140")
		So(ParseIPAddress("192.168.0.140"), ShouldEqual, "192.168.0.140")
		So(ParseIPAddress("[::1:456]"), ShouldEqual, "127.0.0.1")
		So(ParseIPAddress("[::1]"), ShouldEqual, "127.0.0.1")
		So(ParseIPAddress("::1"), ShouldEqual, "127.0.0.1")
		So(ParseIPAddress("::1:123"), ShouldEqual, "127.0.0.1")
	})
}

func TestSplitHostPortDefault(t *testing.T) {
	Convey("Test split ip address to host and port", t, func() {
		host, port := SplitHostPortDefault("192.168.0.140:456", "", "")
		So(host, ShouldEqual, "192.168.0.140")
		So(port, ShouldEqual, "456")

		host, port = SplitHostPortDefault("192.168.0.140", "", "123")
		So(host, ShouldEqual, "192.168.0.140")
		So(port, ShouldEqual, "123")

		host, port = SplitHostPortDefault("[::1:456]", "", "")
		So(host, ShouldEqual, "::1")
		So(port, ShouldEqual, "456")

		host, port = SplitHostPortDefault("[::1]", "", "123")
		So(host, ShouldEqual, "::1")
		So(port, ShouldEqual, "123")

		host, port = SplitHostPortDefault("::1:123", "", "")
		So(host, ShouldEqual, "::1")
		So(port, ShouldEqual, "123")

		host, port = SplitHostPortDefault("::1", "", "123")
		So(host, ShouldEqual, "::1")
		So(port, ShouldEqual, "123")

		host, port = SplitHostPortDefault(":456", "1.2.3.4", "")
		So(host, ShouldEqual, "1.2.3.4")
		So(port, ShouldEqual, "456")

		host, port = SplitHostPortDefault("xyz.rds.amazonaws.com", "", "123")
		So(host, ShouldEqual, "xyz.rds.amazonaws.com")
		So(port, ShouldEqual, "123")

		host, port = SplitHostPortDefault("xyz.rds.amazonaws.com:123", "", "")
		So(host, ShouldEqual, "xyz.rds.amazonaws.com")
		So(port, ShouldEqual, "123")
	})
}

func TestSplitHostPort(t *testing.T) {
	Convey("Test split ip address to host and port", t, func() {
		host, port := SplitHostPort("192.168.0.140:456")
		So(host, ShouldEqual, "192.168.0.140")
		So(port, ShouldEqual, "456")

		host, port = SplitHostPort("192.168.0.140")
		So(host, ShouldEqual, "192.168.0.140")
		So(port, ShouldEqual, "")

		host, port = SplitHostPort("[::1:456]")
		So(host, ShouldEqual, "::1")
		So(port, ShouldEqual, "456")

		host, port = SplitHostPort("[::1]")
		So(host, ShouldEqual, "::1")
		So(port, ShouldEqual, "")

		host, port = SplitHostPort("::1:123")
		So(host, ShouldEqual, "::1")
		So(port, ShouldEqual, "123")

		host, port = SplitHostPort("::1")
		So(host, ShouldEqual, "::1")
		So(port, ShouldEqual, "")

		host, port = SplitHostPort(":456")
		So(host, ShouldEqual, "")
		So(port, ShouldEqual, "456")

		host, port = SplitHostPort("xyz.rds.amazonaws.com")
		So(host, ShouldEqual, "xyz.rds.amazonaws.com")
		So(port, ShouldEqual, "")

		host, port = SplitHostPort("xyz.rds.amazonaws.com:123")
		So(host, ShouldEqual, "xyz.rds.amazonaws.com")
		So(port, ShouldEqual, "123")
	})
}

`
	UtilJson = `
package util

// DynMap defines a dynamic map interface.
type DynMap map[string]interface{}

`
	UtilMath = `
package util

// MaxInt returns the larger of x or y.
func MaxInt(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// MinInt returns the smaller of x or y.
func MinInt(x, y int) int {
	if x > y {
		return y
	}
	return x
}

`
	UtilMd5 = `
package util

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"strings"
)

// Md5Sum calculates the md5sum of a stream
func Md5Sum(reader io.Reader) (string, error) {
	var returnMD5String string
	hash := md5.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return returnMD5String, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil
}

// Md5SumString calculates the md5sum of a string
func Md5SumString(input string) (string, error) {
	buffer := strings.NewReader(input)
	return Md5Sum(buffer)
}

`
	UtilShortIdGenerator = `
package util

import (
	"regexp"

	"github.com/teris-io/shortid"
)

var allowedChars = shortid.DefaultABC

var validUIDPattern = regexp.MustCompile(` +"`^[a-zA-Z0-9\\-\\_]*$`"+`).MatchString

func init() {
	gen, _ := shortid.New(1, allowedChars, 1)
	shortid.SetDefault(gen)
}

// IsValidShortUID checks if short unique identifier contains valid characters
func IsValidShortUID(uid string) bool {
	return validUIDPattern(uid)
}

// GenerateShortUID generates a short unique identifier.
func GenerateShortUID() string {
	return shortid.MustGenerate()
}

`
	UtilStrings = `
package util

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"
)

// StringsFallback2 returns the first of two not empty strings.
func StringsFallback2(val1 string, val2 string) string {
	return stringsFallback(val1, val2)
}

// StringsFallback3 returns the first of three not empty strings.
func StringsFallback3(val1 string, val2 string, val3 string) string {
	return stringsFallback(val1, val2, val3)
}

func stringsFallback(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

// SplitString splits a string by commas or empty spaces.
func SplitString(str string) []string {
	if len(str) == 0 {
		return []string{}
	}

	return regexp.MustCompile("[, ]+").Split(str, -1)
}

// GetAgeString returns a string representing certain time from years to minutes.
func GetAgeString(t time.Time) string {
	if t.IsZero() {
		return "?"
	}

	sinceNow := time.Since(t)
	minutes := sinceNow.Minutes()
	years := int(math.Floor(minutes / 525600))
	months := int(math.Floor(minutes / 43800))
	days := int(math.Floor(minutes / 1440))
	hours := int(math.Floor(minutes / 60))

	if years > 0 {
		return fmt.Sprintf("%dy", years)
	}
	if months > 0 {
		return fmt.Sprintf("%dM", months)
	}
	if days > 0 {
		return fmt.Sprintf("%dd", days)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh", hours)
	}
	if int(minutes) > 0 {
		return fmt.Sprintf("%dm", int(minutes))
	}

	return "< 1m"
}

// ToCamelCase changes kebab case, snake case or mixed strings to camel case. See unit test for examples.
func ToCamelCase(str string) string {
	var finalParts []string
	parts := strings.Split(str, "_")

	for _, part := range parts {
		finalParts = append(finalParts, strings.Split(part, "-")...)
	}

	for index, part := range finalParts[1:] {
		finalParts[index+1] = strings.Title(part)
	}

	return strings.Join(finalParts, "")
}

`
	UtilUrl = `
package util

import (
	"net/url"
	"strings"
)

// URLQueryReader is a URL query type.
type URLQueryReader struct {
	values url.Values
}

// NewURLQueryReader parses a raw query and returns it as a URLQueryReader type.
func NewURLQueryReader(urlInfo *url.URL) (*URLQueryReader, error) {
	u, err := url.ParseQuery(urlInfo.RawQuery)
	if err != nil {
		return nil, err
	}

	return &URLQueryReader{
		values: u,
	}, nil
}

// Get parse parameters from an URL. If the parameter does not exist, it returns
// the default value.
func (r *URLQueryReader) Get(name string, def string) string {
	val := r.values[name]
	if len(val) == 0 {
		return def
	}

	return val[0]
}

// JoinURLFragments joins two URL fragments into only one URL string.
func JoinURLFragments(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")

	if len(b) == 0 {
		return a
	}

	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

`
	UtilUrlTest = `
package util

import (
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUrl(t *testing.T) {

	Convey("When joining two urls where right hand side is empty", t, func() {
		result := JoinURLFragments("http://localhost:8080", "")

		So(result, ShouldEqual, "http://localhost:8080")
	})

	Convey("When joining two urls where right hand side is empty and lefthand side has a trailing slash", t, func() {
		result := JoinURLFragments("http://localhost:8080/", "")

		So(result, ShouldEqual, "http://localhost:8080/")
	})

	Convey("When joining two urls where neither has a trailing slash", t, func() {
		result := JoinURLFragments("http://localhost:8080", "api")

		So(result, ShouldEqual, "http://localhost:8080/api")
	})

	Convey("When joining two urls where lefthand side has a trailing slash", t, func() {
		result := JoinURLFragments("http://localhost:8080/", "api")

		So(result, ShouldEqual, "http://localhost:8080/api")
	})

	Convey("When joining two urls where righthand side has preceding slash", t, func() {
		result := JoinURLFragments("http://localhost:8080", "/api")

		So(result, ShouldEqual, "http://localhost:8080/api")
	})

	Convey("When joining two urls where righthand side has trailing slash", t, func() {
		result := JoinURLFragments("http://localhost:8080", "api/")

		So(result, ShouldEqual, "http://localhost:8080/api/")
	})

	Convey("When joining two urls where lefthand side has a trailing slash and righthand side has preceding slash", t, func() {
		result := JoinURLFragments("http://localhost:8080/", "/api/")

		So(result, ShouldEqual, "http://localhost:8080/api/")
	})
}

func TestNewURLQueryReader(t *testing.T) {
	u, _ := url.Parse("http://www.abc.com/foo?bar=baz&bar2=baz2")
	uqr, _ := NewURLQueryReader(u)

	Convey("when trying to retrieve the first query value", t, func() {
		result := uqr.Get("bar", "foodef")
		So(result, ShouldEqual, "baz")
	})

	Convey("when trying to retrieve the second query value", t, func() {
		result := uqr.Get("bar2", "foodef")
		So(result, ShouldEqual, "baz2")
	})

	Convey("when trying to retrieve from a non-existent key, the default value is returned", t, func() {
		result := uqr.Get("bar3", "foodef")
		So(result, ShouldEqual, "foodef")
	})
}

`
	UtilValidation = `
package util

import (
	"regexp"
	"strings"
)

const (
	emailRegexPattern string = "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_` +"`{\\\\|}~]|[\\\\x{00A0}-\\\\x{D7FF}\\\\x{F900}-\\\\x{FDCF}\\\\x{FDF0}-\\\\x{FFEF}])+(\\\\.([a-zA-Z]|\\\\d|[!#\\\\$%&'\\\\*\\\\+\\\\-\\\\/=\\\\?\\\\^_`"+`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
)

var (
	regexEmail = regexp.MustCompile(emailRegexPattern)
)

// IsEmail checks if a string is a valid email address.
func IsEmail(str string) bool {
	return regexEmail.MatchString(strings.ToLower(str))
}

`
	UtilValidationTest = `
package util

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIsEmail(t *testing.T) {

	Convey("When validating a string that is a valid email", t, func() {
		result := IsEmail("abc@def.com")

		So(result, ShouldEqual, true)
	})

	Convey("When validating a string that is not a valid email", t, func() {
		result := IsEmail("abcdef.com")

		So(result, ShouldEqual, false)
	})
}

`
)

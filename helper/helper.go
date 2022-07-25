package helper

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// Md5 Get md5 for string
func Md5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// Sha256 Get sha256 for string
func Sha256(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}

// Sha256WithSecret Get sha256 by source\secret string
func Sha256WithSecret(s, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// RandomNumber Get random number with length n
func RandomNumber(n int) string {
	rand.Seed(time.Now().UnixNano())
	ans := ""
	for i := 0; i < n; i++ {
		ans += strconv.Itoa(rand.Intn(10))
	}
	return ans
}

// RandomString Get random string with length n
func RandomString(n int) string {
	s := "0123456789AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"
	rand.Seed(time.Now().UnixNano())
	ans := make([]byte, 0, n)
	for i := 0; i < n; i++ {
		ans = append(ans, s[rand.Intn(len(s))])
	}
	return string(ans)
}

// If .
func If(condition bool, trueValue, falseValue interface{}) interface{} {
	if condition {
		return trueValue
	}
	return falseValue
}

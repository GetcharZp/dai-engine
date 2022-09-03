package helper

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"dai-engine/define"
	"dai-engine/logger"
	"encoding/hex"
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"net"
	"strconv"
	"strings"
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

// MapToByte .
func MapToByte(m define.M) []byte {
	b, _ := json.Marshal(m)
	return b
}

// IsLocalIp judge ip in local
func IsLocalIp(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)
	ip4 := ip.To4()
	if ip4 == nil {
		return true
	}
	return ip4[0] == 10 ||
		(ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) ||
		(ip4[0] == 169 && ip4[1] == 254) ||
		(ip4[0] == 192 && ip4[1] == 168)
}

// GetLocalIp get local ip
func GetLocalIp() string {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		logger.Error("DAIL ERROR : " + err.Error())
		return ""
	}
	localAddr := conn.LocalAddr()
	return strings.Split(localAddr.String(), ":")[0]
}

// GetUUID .
func GetUUID() string {
	return uuid.NewV4().String()
}

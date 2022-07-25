package helper

import (
	"fmt"
	"testing"
)

func TestMd5(t *testing.T) {
	fmt.Println(Md5("123456"))
}

func TestSha256(t *testing.T) {
	fmt.Println(Sha256("123456"))
}

func TestSha256WithSecret(t *testing.T) {
	fmt.Println(Sha256WithSecret("123456", ""))
}

func TestRandomNumber(t *testing.T) {
	fmt.Println(RandomNumber(6))
}

func TestRandomString(t *testing.T) {
	fmt.Println(RandomString(6))
}

func TestIf(t *testing.T) {
	fmt.Println(If(true, "true", "false"))
	fmt.Println(If(false, "true", "false"))
}

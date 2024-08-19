package crypt

import (
	"testing"
)

func TestAes(t *testing.T) {
	src := "https://gitlab.github.com/lightmen/nami.git"
	key := "1234567890abcdef"
	buf, err := AesEncrypt([]byte(src), key)
	if err != nil {
		t.Fatalf("aes encrypt error: %s\n", err.Error())
	}

	buf, err = AesDecrypt(buf, key)
	if err != nil {
		t.Fatalf("aes decrypt error: %s\n", err.Error())
	}

	dst := string(buf)
	if len(dst) != len(src) {
		t.Fatalf("src dst string length not equal")
	}
	if dst != src {
		t.Fatalf("src dst string not equal")
	}
}

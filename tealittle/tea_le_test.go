package tealittle_test

import (
	"bytes"
	"crypto/cipher"
	hx "encoding/hex"
	"strings"
	"testing"

	"github.com/dolmen-go/legodim/tealittle"
)

func hex(s string) []byte {
	s = strings.Replace(s, " ", "", -1)
	b, err := hx.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

var _ cipher.Block = teaLE{nil}

func TestEnc(t *testing.T) {
	key := hex("CA3DE8C7 011E608C 58C285A0 9DB48B3E")

	cipher, err := tealittle.NewCipher(key)
	if err != nil {
		t.Fatal("Unexpected failure:", err)
	}

	ciphered := hex("7D B6 72 0B 2A F8 58 0D")
	data := make([]byte, tealittle.BlockSize)
	cipher.Decrypt(data, ciphered)
	expected := hex("1B 00 00 00 1B 00 00 00")

	if !bytes.Equal(data, expected) {
		t.Fatalf("Decrypt: got [% X], expected [% X]", data, expected)
	}

	ciphered2 := make([]byte, tealittle.BlockSize)
	cipher.Encrypt(ciphered2, data)

	if !bytes.Equal(ciphered2, ciphered) {
		t.Fatalf("Decrypt: got [% X], expected [% X]", ciphered2, ciphered)
	}
}

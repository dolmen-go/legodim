// Package tag exposes LEGO Dimensions toy tags related functions and data.
package tag

import (
	"crypto/cipher"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/bits"

	"github.com/DizzyDoug/legodim/tealittle"
)

type UID [7]byte

func (uid UID) String() string {
	return hex.EncodeToString(uid[:])
}

func ParseUID(s string) (UID, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return UID{}, err
	}
	if len(b) != 7 {
		return UID{}, errors.New("invalid length")
	}
	var uid UID
	copy(uid[:], b)
	return uid, nil
}

func MustParseUID(s string) UID {
	uid, err := ParseUID(s)
	if err != nil {
		panic(err)
	}
	return uid
}

// Pwd returns a 4 bytes password.
func (uid UID) Pwd() []byte {
	// Reference code: https://github.com/ags131/node-ld/blob/master/src/lib/PWDGen.js

	tmp := []byte("UUUUUUU(c) Copyright LEGO 2014\xAA\xAA")
	copy(tmp, uid[:])

	var n uint32
	for i := 0; i < 32; i += 4 {
		n = binary.LittleEndian.Uint32(tmp[i:]) +
			bits.RotateLeft32(n, -25) +
			bits.RotateLeft32(n, -10) - n
	}
	var pwd [4]byte
	binary.LittleEndian.PutUint32(pwd[:], n)

	return pwd[:]
}

type Key struct {
	key    [16]byte
	cipher cipher.Block
}

func (k *Key) String() string {
	return fmt.Sprintf("[%X %X %X %X]", k.key[0:4], k.key[4:8], k.key[8:12], k.key[12:16])
}

func (uid UID) Key() *Key {
	tmp := []byte{
		0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xb7,
		0xd5, 0xd7, 0xe6, 0xe7,
		0xba, 0x3c, 0xa8, 0xd8,
		0x75, 0x47, 0x68, 0xcf,
		0x23, 0xe9, 0xfe, 0xaa,
	}
	copy(tmp, uid[:])

	round := func(n uint32, i int, final bool) uint32 {
		var m uint32
		if final {
			m = 0xAA000000 | (uint32(tmp[i+2]) << 16) | uint32(tmp[i+1])<<8 | uint32(tmp[i])
		} else {
			m = binary.LittleEndian.Uint32(tmp[i:])
		}
		return m +
			bits.RotateLeft32(n, -25) +
			bits.RotateLeft32(n, -10) - n
	}

	var k Key
	n := round(0, 0, false)
	n = round(n, 4, false)
	binary.LittleEndian.PutUint32(k.key[0:4], round(n, 8, true))
	n = round(n, 8, false)
	binary.LittleEndian.PutUint32(k.key[4:8], round(n, 12, true))
	n = round(n, 12, false)
	binary.LittleEndian.PutUint32(k.key[8:12], round(n, 16, true))
	n = round(n, 16, false)
	binary.LittleEndian.PutUint32(k.key[12:16], round(n, 20, true))

	// k.cipher, _ = tea.NewCipher(k.key[:])
	k.cipher, _ = tealittle.NewCipher(k.key[:])
	return &k
}

// Encrypt encrypts 8 bytes from src into dst.
func (k *Key) Encrypt(dst, src []byte) {
	k.cipher.Encrypt(dst[:8], dst[:8])
}

// Decrypt decrypts 8 bytes from src into dst.
func (k *Key) Decrypt(dst, src []byte) {
	k.cipher.Decrypt(dst[:8], dst[:8])
}

func (k *Key) EncryptCharacter(c Character) []byte {
	var b [8]byte
	binary.LittleEndian.PutUint32(b[0:4], uint32(c))
	binary.LittleEndian.PutUint32(b[4:8], uint32(c))
	k.cipher.Encrypt(b[:], b[:])
	return b[:]
}

func (k *Key) DecryptCharacter(src []byte) Character {
	var b [8]byte
	k.cipher.Decrypt(b[:8], src)
	return Character(binary.LittleEndian.Uint32(b[0:]))
}

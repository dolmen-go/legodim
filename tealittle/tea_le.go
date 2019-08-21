// Package tealittle provides a TEA implementation where bytes are swapped 4 by 4.
//
// Common implementations of TEA is BigEndian based. This implementation wraps
// the golang.org/crypto/tea implementation to use LittleEndian.
//
// This package is provided to allow interoperability with Little Endian devices
// which use non-standard TEA.
package tealittle

import (
	"crypto/cipher"

	"golang.org/x/crypto/tea"
)

const (
	BlockSize = tea.BlockSize
	KeySize   = tea.KeySize
)

func swCopy4(dst, src []byte) {
	// binary.BigEndian.PutUint32(dst, binary.LittleEndian.Uint32(src))
	dst[0], dst[1], dst[2], dst[3] = src[3], src[2], src[1], src[0]
}

func swCopy8(dst, src []byte) {
	swCopy4(dst, src)
	swCopy4(dst[4:], src[4:])
}

func swCopy16(dst, src []byte) {
	swCopy8(dst, src)
	swCopy8(dst[8:], src[8:])
}

func sw4(b []byte) {
	swCopy4(b, b)
}

func sw8(b []byte) {
	sw4(b)
	sw4(b[4:])
}

type teaLE struct {
	cipher.Block
}

// NewCipher creates a new LittleEndian TEA block cipher with 64 rounds.
func NewCipher(key []byte) (cipher.Block, error) {
	if len(key) != KeySize {
		return tea.NewCipher(key)
	}

	var keyLE [tea.KeySize]byte
	swCopy16(keyLE[:], key)

	c, _ := tea.NewCipher(keyLE[:])
	return teaLE{c}, nil
}

func (c teaLE) Encrypt(dst, src []byte) {
	swCopy8(dst, src)
	c.Block.Encrypt(dst[:BlockSize], dst[:BlockSize])
	sw8(dst)
}

func (c teaLE) Decrypt(dst, src []byte) {
	swCopy8(dst, src)
	c.Block.Decrypt(dst[:BlockSize], dst[:BlockSize])
	sw8(dst)
}

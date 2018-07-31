package hashing

import (
	"crypto/sha256"
)

// Sha256Hasher implements the Hasher interface and computes the crypto/sha256
// internal function.
type Sha256Hasher struct{}

func (s Sha256Hasher) Do(data ...[]byte) []byte {
	hash := sha256.New()

	for i := 0; i < len(data); i++ {
		hash.Write(data[i])
	}

	return hash.Sum(nil)[:]
}

func (s Sha256Hasher) Cipher(id []byte, data ...[]byte) []byte {
	b := [][]byte{id}
	for _, n := range data {
		b = append(b, []byte(n))
	}

	return s.Do(b...)
}

func (s Sha256Hasher) Len() uint64 { return uint64(256) }

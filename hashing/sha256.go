package hashing

import (
	"crypto/sha256"
	"hash"
)

// Sha256 implements the Hasher interface and computes the crypto/sha256
// internal function.
type Sha256 struct {
	h hash.Hash
}

func NewSha256() *Sha256 {
	return &Sha256{
		h: sha256.New(),
	}
}

func (s Sha256) Do(data ...[]byte) []byte {
	s.h.Reset()

	for i := 0; i < len(data); i++ {
		s.h.Write(data[i])
	}

	return s.h.Sum(nil)[:]
}

func (s Sha256) Cipher(id []byte, data ...[]byte) []byte {
	b := [][]byte{id}
	for _, n := range data {
		b = append(b, []byte(n))
	}

	return s.Do(b...)
}

func (s Sha256) Len() uint64 { return uint64(256) }

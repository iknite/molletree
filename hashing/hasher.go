package hashing

import (
	"crypto/sha256"
	"fmt"
)

// Hasher is the public interface to be used as placeholder for the concrete
// implementations.
type Hasher interface {
	Cipher([]byte, ...[]byte) []byte
	Do(...[]byte) []byte
	Len() uint64
}

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

// XorHasher implements the Hasher interface and computes a xor function.
// Handy for testing hash tree implementations.
type XorHasher struct{}

func (x XorHasher) Do(data ...[]byte) []byte {
	fmt.Println("*", data)

	var result byte
	for _, elem := range data {
		var sum byte
		for _, b := range elem {
			sum = sum ^ b
		}
		result = result ^ sum
	}
	return []byte{result}
}

func (s XorHasher) Cipher(_ []byte, data ...[]byte) []byte {
	return s.Do(data...)
}

func (s XorHasher) Len() uint64 { return uint64(8) }

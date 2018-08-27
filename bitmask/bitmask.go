package bitmask

import "math/big"

func genMask(size uint) []byte {
	// `m := 1<<size - 1` but for any size of number
	m := new(big.Int).Add(new(big.Int).Lsh(big.NewInt(1), size), big.NewInt(-1))
	return m.Bytes()
}

func filter(digest []byte, size uint64, clear bool) []byte {
	var b []byte
	b = append(b, digest...)

	mask := genMask(uint(size))

	l, m := len(digest), len(mask)

	if clear {
		for i := m - 1; i >= 0; i-- {
			b[l-m+i] &^= mask[i]
		}
	} else {
		for i := m - 1; i >= 0; i-- {
			b[l-m+i] |= mask[i]
		}
	}
	return b
}

func SetLeft(digest []byte, size uint64) []byte {
	return filter(digest, size, false)
}

func ClearLeft(digest []byte, size uint64) []byte {
	return filter(digest, size, true)
}

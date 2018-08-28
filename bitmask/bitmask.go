package bitmask

import "math/big"

func genMask(size uint) *big.Int {

	// `m := 1<<size - 1` but for any number
	m := new(big.Int).Add(new(big.Int).Lsh(big.NewInt(1), size), big.NewInt(-1))

	// it will return the smallest slice of bytes that represents the mask.
	return m

}

func filter(digest []byte, size uint64, clear bool) (b []byte) {

	mask := genMask(uint(size))

	d := new(big.Int).SetBytes(digest)

	if clear {
		// `b = d &^ mask` set mask part to zeros.
		b = new(big.Int).AndNot(d, mask).Bytes()

	} else {
		// `b = d | mask` set mask part to ones.
		b = new(big.Int).Or(d, mask).Bytes()
	}

	return

}

// SetLeft will return a new byte slice will the latest (size) bits as ones.
func SetLeft(digest []byte, size uint64) []byte {
	return filter(digest, size, false)
}

// ClearLeft will return a new byte slice will the latest (size) bits as zeros.
func ClearLeft(digest []byte, size uint64) []byte {
	return filter(digest, size, true)
}

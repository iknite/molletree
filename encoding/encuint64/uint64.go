package encuint64

import "encoding/binary"

func ToBytes(i uint64) []byte {
	// INFO: https://golang.org/ref/spec#Size_and_alignment_guarantees
	a := make([]byte, 8)
	binary.BigEndian.PutUint64(a, i)
	return a
}

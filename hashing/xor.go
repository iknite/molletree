package hashihg

// XorHasher implements the Hasher interface and computes a xor function.
// Handy for testing hash tree implementations.
type XorHasher struct{}

func (x XorHasher) Do(data ...[]byte) []byte {
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

package hashing

// Xor implements the Hasher interface and computes a xor function.
// Handy for testing hash tree implementations.
type Xor struct{}

func NewXor() *Xor {
	return new(Xor)
}

func (x Xor) Do(data ...[]byte) []byte {
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

func (s Xor) Cipher(_ []byte, data ...[]byte) []byte {
	return s.Do(data...)
}

func (s Xor) Len() uint64 { return uint64(8) }

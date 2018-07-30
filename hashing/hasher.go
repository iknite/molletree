package hashing

// Hasher is the public interface to be used as placeholder for the concrete
// implementations.
type Hasher interface {
	Cipher([]byte, ...[]byte) []byte
	Do(...[]byte) []byte
	Len() uint64
}

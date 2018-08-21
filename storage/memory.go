package storage

import "github.com/iknite/molletree/encoding/encbytes"

type MemoryStore map[string][]byte

func NewMemoryStore() MemoryStore {
	return make(MemoryStore)
}

func (s MemoryStore) Get(id []byte) ([]byte, bool) {
	res, ok := s[encbytes.ToPrettyId(id)]
	if ok {
		return res, true
	}

	return make([]byte, 0), false
}

func (s MemoryStore) Set(id []byte, data []byte) {
	s[encbytes.ToPrettyId(id)] = data
}

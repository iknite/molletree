package storage

type Store map[string][]byte

func NewStore() Store {
	return make(Store)
}

func (s *Store) Get(id string) ([]byte, bool) {
	res, ok := (*s)[id]
	if ok {
		return res, true
	}

	return make([]byte, 0), false
}

func (s *Store) Set(id string, data []byte) {
	(*s)[id] = data
}

package storage

type Store struct {
	Std map[string][]byte
}

var std Store

func NewStore() Store {
	return Store{Std: make(map[string][]byte)}
}

func init() {
	std = NewStore()
}

func (s *Store) Get(id string) ([]byte, bool) {
	res, ok := s.Std[id]
	if ok {
		return res, true
	}

	return make([]byte, 0), false
}

func (s *Store) Set(id string, data []byte) {
	s.Std[id] = data
}

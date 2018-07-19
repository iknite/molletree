package storage

var std map[string][]byte

func init() {
	std = make(map[string][]byte)
}

func Get(id string) ([]byte, bool) {
	res, ok := std[id]
	if ok {
		return res, true
	}

	return make([]byte, 0), false
}

func Set(id string, data []byte) {
	std[id] = data
}

package storage

type Storer interface {
	Get(id []byte) ([]byte, bool)
	Set(id []byte, data []byte)
}

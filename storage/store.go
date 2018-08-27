package storage

type readable interface {
	Get(id []byte) ([]byte, bool)
}

type writable interface {
	Set(id, data []byte)
}

type mergeable interface {
	Merge(data MemoryStore)
}

type cacheable interface {
	GetRange(start, end []byte) MemoryStore
	SetAndPrefetch(id, data []byte, level uint64) MemoryStore
}

type Storer interface {
	readable
	writable
}

type Merger interface {
	readable
	writable
	mergeable
}

type Cacher interface {
	readable
	writable
	cacheable
}

package history

import (
	"testing"

	assert "github.com/stretchr/testify/require"

	"github.com/iknite/molletree/encoding/encstring"
	"github.com/iknite/molletree/hashing"
	"github.com/iknite/molletree/storage"
)

func TestVerifyIncremental(t *testing.T) {

	testCases := []struct {
		store                  storage.Store
		start, end             uint64
		startDigest, endDigest []byte
	}{
		{
			storage.Store{"0|1": []uint8{0x1}, "2|0": []uint8{0x2}, "3|0": []uint8{0x3}, "4|1": []uint8{0x1}, "6|0": []uint8{0x6}},
			2, 6, []byte{0x3}, []byte{0x7},
		},
		{
			storage.Store{"0|1": []uint8{0x1}, "2|0": []uint8{0x2}, "3|0": []uint8{0x3}, "4|1": []uint8{0x1}, "6|0": []uint8{0x6}, "7|0": []uint8{0x7}},
			2, 7, []byte{0x3}, []byte{0x0},
		},
		{
			storage.Store{"0|2": []uint8{0x0}, "4|0": []uint8{0x4}, "5|0": []uint8{0x5}, "6|0": []uint8{0x6}},
			4, 6, []byte{0x4}, []byte{0x7},
		},
		{
			storage.Store{"0|2": []uint8{0x0}, "4|0": []uint8{0x4}, "5|0": []uint8{0x5}, "6|0": []uint8{0x6}, "7|0": []uint8{0x7}},
			4, 7, []byte{0x4}, []byte{0x0},
		},
		{
			storage.Store{"2|0": []uint8{0x2}, "3|0": []uint8{0x3}, "4|0": []uint8{0x4}, "0|1": []uint8{0x1}},
			2, 4, []byte{0x3}, []byte{0x4},
		},
		{
			storage.Store{"0|2": []uint8{0x0}, "4|1": []uint8{0x1}, "6|0": []uint8{0x6}, "7|0": []uint8{0x7}},
			6, 7, []byte{0x7}, []byte{0x0},
		},
	}

	h := hashing.NewXor()

	for _, c := range testCases {
		proof := &IncrementalProof{c.start, c.end, c.store, h}
		assert.True(t, proof.Verify(c.startDigest, c.endDigest))
	}
}

func max(x, y int) uint64 {
	if x > y {
		return uint64(x)
	}
	return uint64(y)
}

func TestProveAndVerifyConsecutivelyN(t *testing.T) {
	tree := &Tree{version: 0, hasher: hashing.NewXor(), store: storage.NewStore()}
	digests := make(map[uint64][]byte)

	for i := uint64(0); i < 10; i++ {
		node := &Node{index: i, layer: 0, tree: tree}
		node.tree.store.Set(node.String(), encstring.ToBytes(string(i)))
		node.Commitment()
		tree.version += 1

		start := max(0, int(i-1))

		digests[i] = node.Root(i).Hash(i)

		pf := tree.ProveIncremental(start, i)

		assert.True(t, pf.Verify(digests[start], digests[i]),
			"The proof should verfify correctly")
	}
}

func TestProveAndVerifyNonConsecutively(t *testing.T) {
	tree := &Tree{version: 0, hasher: hashing.NewXor(), store: storage.NewStore()}

	const size uint64 = 10
	digests := make(map[uint64][]byte)

	for i := uint64(0); i < 10; i++ {
		index := uint64(i)
		node := &Node{index: index, layer: 0, tree: tree}
		node.tree.store.Set(node.String(), encstring.ToBytes(string(i)))
		digests[i] = node.Commitment()
		tree.version += 1
	}

	for i := uint64(0); i < size-1; i++ {
		for j := i + 1; j < size; j++ {
			pf := tree.ProveIncremental(i, j)
			assert.True(t, pf.Verify(digests[i], digests[j]))
		}
	}

}

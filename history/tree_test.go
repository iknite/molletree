package history

import (
	"testing"

	"github.com/iknite/bygone-tree/hashing"
	"github.com/iknite/bygone-tree/storage"
	assert "github.com/stretchr/testify/require"
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

	h := &hashing.XorHasher{}

	for _, c := range testCases {
		proof := &IncrementalProof{c.start, c.end, c.store, h}
		assert.True(t, proof.Verify(c.startDigest, c.endDigest))
	}
}

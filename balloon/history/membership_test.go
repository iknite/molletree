package history

import (
	"testing"

	assert "github.com/stretchr/testify/require"

	"github.com/iknite/molletree/encoding/encstring"
	"github.com/iknite/molletree/hashing"
	"github.com/iknite/molletree/storage"
)

func TestProveMembership(t *testing.T) {

	testCases := []struct {
		eventDigest []byte
		auditPath   storage.Storer
	}{
		{[]byte{0x0}, storage.MemoryStore{"0|0": []uint8{0x0}}},
		{[]byte{0x1}, storage.MemoryStore{"0|0": []uint8{0x0}, "1|0": []uint8{0x1}}},
		{[]byte{0x2}, storage.MemoryStore{"0|1": []uint8{0x1}, "2|0": []uint8{0x2}}},
		{[]byte{0x3}, storage.MemoryStore{"0|1": []uint8{0x1}, "2|0": []uint8{0x2}, "3|0": []uint8{0x3}}},
		{[]byte{0x4}, storage.MemoryStore{"0|2": []uint8{0x0}, "4|0": []uint8{0x4}}},
		{[]byte{0x5}, storage.MemoryStore{"0|2": []uint8{0x0}, "4|0": []uint8{0x4}, "5|0": []uint8{0x5}}},
		{[]byte{0x6}, storage.MemoryStore{"0|2": []uint8{0x0}, "4|1": []uint8{0x1}, "6|0": []uint8{0x6}}},
		{[]byte{0x7}, storage.MemoryStore{"0|2": []uint8{0x0}, "4|1": []uint8{0x1}, "6|0": []uint8{0x6}, "7|0": []uint8{0x7}}},
		{[]byte{0x8}, storage.MemoryStore{"0|3": []uint8{0x0}, "8|0": []uint8{0x8}}},
		{[]byte{0x9}, storage.MemoryStore{"0|3": []uint8{0x0}, "8|0": []uint8{0x8}, "9|0": []uint8{0x9}}},
	}

	tree := &Tree{version: 0, hasher: &hashing.Xor{}, store: storage.NewMemoryStore()}

	for i, c := range testCases {
		index := uint64(i)
		node := &Node{index: index, layer: 0, tree: tree}
		node.tree.store.Set(node.Id(), c.eventDigest)
		node.Commitment()

		assert.Equalf(t, c.auditPath, node.AuditPath(node.index), "Incorrect audit path for index %d", i)

		node.tree.version += 1
	}
}

func TestProveMembershipWithInvalidTargetVersion(t *testing.T) {
	tree := &Tree{version: 0, hasher: &hashing.Xor{}, store: storage.NewMemoryStore()}
	tree.Add(encstring.ToBytes("Event1"))

	defer func() {
		if r := recover(); r == nil {
			t.Error("should raise an error")
		}
	}()
	tree.ProveMembership([]byte{0x0}, 1, 0)
}

func TestProveMembershipNonConsecutive(t *testing.T) {
	tree := &Tree{version: 0, hasher: &hashing.Xor{}, store: storage.NewMemoryStore()}

	// add nine events
	for i := uint64(0); i < 9; i++ {
		node := &Node{index: i, layer: 0, tree: tree}
		node.tree.store.Set(node.Id(), encstring.ToBytes(string(i)))
		node.Commitment()
	}

	// query for membership with event 0 and version 8
	pfNode := &Node{index: 0, layer: 0, tree: tree}

	au := storage.MemoryStore{"0|0": []uint8{0x0}, "1|0": []uint8{0x1}, "2|1": []uint8{0x1}, "4|2": []uint8{0x0}, "8|3": []uint8{0x8}}

	assert.Equal(t, au, pfNode.AuditPath(8), "Invalid audit path")
}

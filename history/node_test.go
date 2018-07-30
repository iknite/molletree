package history

import (
	"fmt"
	"testing"

	assert "github.com/stretchr/testify/require"

	"github.com/iknite/molletree/encoding/encstring"
	"github.com/iknite/molletree/hashing"
	"github.com/iknite/molletree/storage"
)

func TestCommitment(t *testing.T) {

	testCases := []struct {
		eventDigest []byte
		commitment  []byte
	}{
		{[]byte{0x0}, []byte{0x0}},
		{[]byte{0x1}, []byte{0x1}},
		{[]byte{0x2}, []byte{0x3}},
		{[]byte{0x3}, []byte{0x0}},
		{[]byte{0x4}, []byte{0x4}},
		{[]byte{0x5}, []byte{0x1}},
		{[]byte{0x6}, []byte{0x7}},
		{[]byte{0x7}, []byte{0x0}},
		{[]byte{0x8}, []byte{0x8}},
		{[]byte{0x9}, []byte{0x1}},
	}

	// Note that we are using fake hashing functions and the index
	// as the value of the event's digest to make predictable hashes
	tree := &Tree{version: 0, hasher: &hashing.XorHasher{}, store: storage.NewStore()}

	for i, c := range testCases {

		// almost like tree.Add except we provide the digest to allow easier
		// tests.
		node := &Node{index: uint64(i), layer: 0, tree: tree}
		node.tree.store.Set(node.String(), c.eventDigest)

		commitment := node.Commitment()
		node.tree.version += 1

		assert.Equalf(t, c.commitment, commitment, "Incorrect commitment for index %d", i)
	}

}

func TestProveMembership(t *testing.T) {

	testCases := []struct {
		eventDigest []byte
		auditPath   storage.Store
	}{
		{[]byte{0x0}, storage.Store{"0|0": []uint8{0x0}}},
		{[]byte{0x1}, storage.Store{"0|0": []uint8{0x0}, "1|0": []uint8{0x1}}},
		{[]byte{0x2}, storage.Store{"0|1": []uint8{0x1}, "2|0": []uint8{0x2}}},
		{[]byte{0x3}, storage.Store{"0|1": []uint8{0x1}, "2|0": []uint8{0x2}, "3|0": []uint8{0x3}}},
		{[]byte{0x4}, storage.Store{"0|2": []uint8{0x0}, "4|0": []uint8{0x4}}},
		{[]byte{0x5}, storage.Store{"0|2": []uint8{0x0}, "4|0": []uint8{0x4}, "5|0": []uint8{0x5}}},
		{[]byte{0x6}, storage.Store{"0|2": []uint8{0x0}, "4|1": []uint8{0x1}, "6|0": []uint8{0x6}}},
		{[]byte{0x7}, storage.Store{"0|2": []uint8{0x0}, "4|1": []uint8{0x1}, "6|0": []uint8{0x6}, "7|0": []uint8{0x7}}},
		{[]byte{0x8}, storage.Store{"0|3": []uint8{0x0}, "8|0": []uint8{0x8}}},
		{[]byte{0x9}, storage.Store{"0|3": []uint8{0x0}, "8|0": []uint8{0x8}, "9|0": []uint8{0x9}}},
	}

	tree := &Tree{version: 0, hasher: &hashing.XorHasher{}, store: storage.NewStore()}

	for i, c := range testCases {
		index := uint64(i)
		node := &Node{index: index, layer: 0, tree: tree}
		node.tree.store.Set(node.String(), c.eventDigest)
		node.Commitment()

		assert.Equalf(t, c.auditPath, node.AuditPath(node.index), "Incorrect audit path for index %d", i)

		node.tree.version += 1
	}
}

func TestProveMembershipWithInvalidTargetVersion(t *testing.T) {
	tree := &Tree{version: 0, hasher: &hashing.XorHasher{}, store: storage.NewStore()}
	tree.Add("Event1")

	defer func() {
		if r := recover(); r == nil {
			t.Error("should raise an error")
		}
	}()
	tree.ProveMembership([]byte{0x0}, 1, 0)
}

func TestProveMembershipNonConsecutive(t *testing.T) {
	tree := &Tree{version: 0, hasher: &hashing.XorHasher{}, store: storage.NewStore()}

	// add nine events
	for i := uint64(0); i < 9; i++ {
		node := &Node{index: i, layer: 0, tree: tree}
		node.tree.store.Set(node.String(), encstring.ToBytes(string(i)))
		node.Commitment()
	}

	// query for membership with event 0 and version 8
	pfNode := &Node{index: 0, layer: 0, tree: tree}

	au := storage.Store{"0|0": []uint8{0x0}, "1|0": []uint8{0x1}, "2|1": []uint8{0x1}, "4|2": []uint8{0x0}, "8|3": []uint8{0x8}}

	assert.Equal(t, au, pfNode.AuditPath(8), "Invalid audit path")
}

func max(x, y int) uint64 {
	if x > y {
		return uint64(x)
	}
	return uint64(y)
}

func TestProveAndVerifyConsecutivelyN(t *testing.T) {
	tree := &Tree{version: 0, hasher: &hashing.XorHasher{}, store: storage.NewStore()}
	digests := make(map[uint64][]byte)

	for i := uint64(0); i < 10; i++ {
		node := &Node{index: i, layer: 0, tree: tree}
		node.tree.store.Set(node.String(), encstring.ToBytes(string(i)))
		node.Commitment()
		tree.version += 1

		start := max(0, int(i-1))

		digests[i] = node.Root(i).Hash(i)

		pf := tree.ProveIncremental(start, i)

		assert.True(t, pf.Verify(digests[start], digests[i]), "The proof should verfify correctly")
	}
}

func TestProveAndVerifyNonConsecutively(t *testing.T) {
	tree := &Tree{version: 0, hasher: &hashing.XorHasher{}, store: storage.NewStore()}

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

			vf := pf.Verify(digests[i], digests[j])

			if !vf {
				fmt.Println(vf, i, j, ">", digests[i], digests[j], ":", pf.store, ",", tree.store)
			}
		}
	}

}

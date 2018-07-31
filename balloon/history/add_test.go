package history

import (
	"testing"

	assert "github.com/stretchr/testify/require"

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

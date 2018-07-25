package history

import (
	"fmt"
	"testing"

	"github.com/iknite/bygone-tree/hashing"
	"github.com/iknite/bygone-tree/storage"
	assert "github.com/stretchr/testify/require"
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

		commitment := node.Commitment(n.index)
		node.tree.version += 1

		assert.Equalf(t, c.commitment, commitment, "Incorrect commitment for index %d", i)
	}

}

func TestProveMembership(t *testing.T) {

	testCases := []struct {
		eventDigest []byte
		auditPath   storage.Store
	}{
		{[]byte{0x0}, storage.Store{Std: map[string][]byte{"0|0": []uint8{0x0}}}},
		{[]byte{0x1}, storage.Store{Std: map[string][]byte{"0|0": []uint8{0x0}, "1|0": []uint8{0x1}}}},
		{[]byte{0x2}, storage.Store{Std: map[string][]byte{"0|1": []uint8{0x1}, "2|0": []uint8{0x2}}}},
		{[]byte{0x3}, storage.Store{Std: map[string][]byte{"0|1": []uint8{0x1}, "2|0": []uint8{0x2}, "3|0": []uint8{0x3}}}},
		{[]byte{0x4}, storage.Store{Std: map[string][]byte{"0|2": []uint8{0x0}, "4|0": []uint8{0x4}}}},
		{[]byte{0x5}, storage.Store{Std: map[string][]byte{"0|2": []uint8{0x0}, "4|0": []uint8{0x4}, "5|0": []uint8{0x5}}}},
		{[]byte{0x6}, storage.Store{Std: map[string][]byte{"0|2": []uint8{0x0}, "4|1": []uint8{0x1}, "6|0": []uint8{0x6}}}},
		{[]byte{0x7}, storage.Store{Std: map[string][]byte{"0|2": []uint8{0x0}, "4|1": []uint8{0x1}, "6|0": []uint8{0x6}, "7|0": []uint8{0x7}}}},
		{[]byte{0x8}, storage.Store{Std: map[string][]byte{"0|3": []uint8{0x0}, "8|0": []uint8{0x8}}}},
		{[]byte{0x9}, storage.Store{Std: map[string][]byte{"0|3": []uint8{0x0}, "8|0": []uint8{0x8}, "9|0": []uint8{0x9}}}},
	}

	tree := &Tree{version: 0, hasher: &hashing.XorHasher{}, store: storage.NewStore()}

	for i, c := range testCases {
		index := uint64(i)
		node := &Node{index: index, layer: 0, tree: tree}
		node.tree.store.Set(node.String(), c.eventDigest)
		node.Commitment(n.index)

		assert.Equalf(t, c.auditPath, node.AuditPath(), "Incorrect audit path for index %d", i)

		node.tree.version += 1
	}
}

func TestProveMembershipWithInvalidTargetVersion(t *testing.T) {
	tree := &Tree{version: 0, hasher: &hashing.XorHasher{}, store: storage.NewStore()}

	tree.Add("Event1")

	pf := tree.MembershipProof([]byte{0x0}, 1, 0)

	fmt.Println(pf)
}

//
// func TestProveMembershipNonConsecutive(t *testing.T) {
// 	frozen, close := openBPlusStorage()
// 	defer close()
//
// 	hasher := new(hashing.XorHasher)
// 	tree := NewTree("treeId", frozen, hasher)
// 	tree.leafHash = fakeLeafHasherCleanF(hasher)
// 	tree.interiorHash = fakeInteriorHasherCleanF(hasher)
// 	// Note that we are using fake hashing functions and the index
// 	// as the value of the event's digest to make predictable hashes
//
// 	// add nine events
// 	for i := uint64(0); i < 9; i++ {
// 		eventDigest := uint64AsBytes(i)
// 		index := uint64AsBytes(i)
// 		_, err := tree.Add(eventDigest, index)
// 		assert.NoError(t, err, "Error while adding to the tree")
// 	}
//
// 	// query for membership with event 0 and version 8
// 	pf, err := tree.ProveMembership([]byte{0x0}, 0, 8)
// 	assert.NoError(t, err, "Error proving membership")
// 	expectedAuditPath := proof.AuditPath{"0|0": []uint8{0x0}, "1|0": []uint8{0x1}, "2|1": []uint8{0x1}, "4|2": []uint8{0x0}, "8|3": []uint8{0x8}}
// 	assert.Equal(t, expectedAuditPath, pf.AuditPath(), "Invalid audit path")
// }

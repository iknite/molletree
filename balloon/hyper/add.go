package hyper

import (
	"bytes"

	"github.com/iknite/molletree/storage"
)

func (t *Tree) Add(digest []byte) []byte {
	s := storage.NewMemoryStore()
	n := &Node{base: digest, layer: 0, tree: t, opStore: s}
	p := n.tree.store.SetAndPrefetch(n.Id(), digest, n.tree.cacheLevel)

	s.Merge(p)

	return n.Root().Hash(digest)
}

func (n *Node) Get() (hash []byte) {
	hash, _ = n.Store().Get(n.Id())
	return
}

func (n *Node) Set(hash []byte) {
	n.Store().Set(n.Id(), hash)
}

func (n *Node) Hash(target []byte) (hash []byte) {
	if n.layer == 0 {
		return n.Get()
	}

	next := n.Next(target)
	left := n.Left()
	right := n.Right()

	if bytes.Equal(left.base, next.base) {
		hash = n.tree.hasher.Cipher(n.Id(), left.Hash(target), right.Get())
	} else {
		hash = n.tree.hasher.Cipher(n.Id(), left.Get(), right.Hash(target))
	}
	n.Set(hash)

	return
}

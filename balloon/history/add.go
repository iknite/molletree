package history

func (t *Tree) Add(message []byte) (commitment []byte, digest []byte) {
	// Add a leaf node
	node := &Node{index: t.version, layer: 0, tree: t}
	t.version += 1

	digest = t.hasher.Cipher(node.Id(), message)
	t.store.Set(node.Id(), digest)

	commitment = node.Commitment()

	return
}

func (n *Node) Commitment() []byte {
	return n.Root(n.index).Hash(n.index)
}

func (n *Node) Hash(version uint64) (hash []byte) {

	if n.index > version {
		return // if you're trying to get the future, return nil
	}

	id := n.Id()

	if n.layer == 0 || n.Capacity() < version {
		var ok bool
		hash, ok = n.tree.store.Get(id)
		if ok {
			return // you're getting the past so it's cached
		}
	}

	hash = n.tree.hasher.Cipher(id, n.Left().Hash(version), n.Right().Hash(version))

	if n.Capacity() == version {
		n.tree.store.Set(id, hash) // is storable when the childrens are complete
	}

	return

}

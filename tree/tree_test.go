package tree

import (
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestHeight(t *testing.T) {
	tree := Tree{}

	tree.length = 0
	assert.Equal(t, tree.height(), 0.0)

	tree.length = 3
	assert.Equal(t, tree.height(), 2.0)

	tree.length = 8
	assert.Equal(t, tree.height(), 4.0)

}

func TestAdd(t *testing.T) {
	tree := NewTree()
	tree.Add("event1")
	assert.Equal(t, storage["0:0"], "event1")
}

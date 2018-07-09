package tree

import (
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestHeight(t *testing.T) {
	tree := Tree{}

	tree.version = 0
	assert.Equal(t, tree.height(), 0.0)

	tree.version = 3
	assert.Equal(t, tree.height(), 2.0)

	tree.version = 8
	assert.Equal(t, tree.height(), 4.0)

}

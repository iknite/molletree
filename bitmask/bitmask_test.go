package bitmask

import (
	"bytes"
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestSetLeft(t *testing.T) {
	digest := []byte{0x68, 0x6f, 0x6c, 0x61}

	assert.True(t, bytes.Equal([]byte{0x68, 0x6f, 0xff, 0xff}, SetLeft(digest, 18)))
	assert.True(t, bytes.Equal([]byte{0x68, 0x6c, 0x00, 0x00}, ClearLeft(digest, 18)))
}

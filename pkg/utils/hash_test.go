package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHasher_HashAndCompare(t *testing.T) {
    h := NewHasher()

    password := "Password123@"
    hash, err := h.Hash(password)
    require.NoError(t, err)
    require.NotEmpty(t, hash)
    assert.True(t, h.Compare(hash, password))
    assert.False(t, h.Compare(hash, "wrong-password"))
}

func TestHasher_CompareWithInvalidHash(t *testing.T) {
    h := NewHasher()
    assert.False(t, h.Compare("not-a-valid-hash", "any"))
}
var _ Hasher = (*hasher)(nil)

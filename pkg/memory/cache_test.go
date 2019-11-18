package memory_test

import (
	"testing"

	"github.com/pmdcosta/exchange/pkg/memory"
	"github.com/stretchr/testify/require"
)

func TestCache_NotFound(t *testing.T) {
	c := memory.New()
	hit, value := c.Load("test")
	require.False(t, hit)
	require.Nil(t, value)
}

func TestCache_Found(t *testing.T) {
	c := memory.New()
	c.Save("test", []byte("found"))
	hit, value := c.Load("test")
	require.True(t, hit)
	require.Equal(t, []byte("found"), value)
}

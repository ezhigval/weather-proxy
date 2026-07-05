package provider

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockFetch(t *testing.T) {
	p := NewMock()
	w, err := p.Fetch(context.Background(), "Moscow")
	require.NoError(t, err)
	assert.Equal(t, "Moscow", w.City)
	assert.NotEmpty(t, w.Description)
}

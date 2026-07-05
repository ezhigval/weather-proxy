package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCities(t *testing.T) {
	got := ParseCities(" Moscow , London,Tokyo, , ")
	assert.Equal(t, []string{"Moscow", "London", "Tokyo"}, got)
}

func TestNormalizeCity(t *testing.T) {
	assert.Equal(t, "Berlin", normalizeCity("  Berlin "))
	assert.Equal(t, "", normalizeCity("   "))
}

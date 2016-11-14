package aedstorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUUID(t *testing.T) {
	u, err := NewUUID()
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, 16, len(u))
	assert.Equal(t, 36, len(u.String()))
}

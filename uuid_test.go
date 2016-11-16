package aedstorm

import (
	"crypto/rand"
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

func TestNewUUIDPanic(t *testing.T) {
	oldReader := rand.Reader
	defer func() {
		rand.Reader = oldReader
	}()
	rand.Reader = nil
	assert.Panics(t, func() {
		NewUUID()
	})
}

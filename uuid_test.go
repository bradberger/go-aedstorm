package aedstorm

import (
	"bytes"
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

func TestNewUUIDWithErr(t *testing.T) {
	oldReader := rand.Reader
	defer func() {
		rand.Reader = oldReader
	}()
	rand.Reader = bytes.NewBuffer(nil)
	_, err := NewUUID()
	assert.Error(t, err)
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

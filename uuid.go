package aedstorm

import (
	"crypto/rand"
	"fmt"
)

// UUID is a v4 Universally unique identifier
type UUID [16]byte

// NewUUID creates a new uuid v4
func NewUUID() (*UUID, error) {
	u := &UUID{}
	if _, err := rand.Read(u[:16]); err != nil {
		return nil, err
	}

	u[8] = (u[8] | 0x80) & 0xBf
	u[6] = (u[6] | 0x40) & 0x4f
	return u, nil
}

func (u *UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[:4], u[4:6], u[6:8], u[8:10], u[10:])
}

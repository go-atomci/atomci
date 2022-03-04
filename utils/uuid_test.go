package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUUID(t *testing.T) {
	uuid := NewUUID()

	assert.NotEmpty(t, uuid)
}

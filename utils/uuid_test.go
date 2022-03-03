package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NEWUUID_SHOULD_NOT_EMPTY(t *testing.T) {
	uuid := NewUUID()

	assert.NotEmpty(t, uuid)
}

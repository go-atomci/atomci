package utils

import (
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestAesCrypto(t *testing.T) {
	encrypted := base64.StdEncoding.EncodeToString(AesEny([]byte("Hello")))
	log.Printf("%s", encrypted)
	assert.NotEmpty(t, encrypted)
}

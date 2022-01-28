package utils

import (
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestAesCrypto(t *testing.T) {
	crypted := base64.StdEncoding.EncodeToString(AesEny([]byte("Hello")))
	log.Printf("%s", crypted)
	assert.NotEmpty(t, crypted)
}

package tunnel

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrivateKeyFile(t *testing.T) {
	assert := assert.New(t)
	cwd, _ := os.Getwd()
	validKey := PrivateKeyFile(cwd + "/test_files/id_rsa")
	assert.NotNil(validKey)
	invalidKey := PrivateKeyFile(cwd + "/test_files/id_rsa.pub")
	assert.Nil(invalidKey)
	nosuchKey := PrivateKeyFile(cwd + "/test_files/something.pub")
	assert.Nil(nosuchKey)
}

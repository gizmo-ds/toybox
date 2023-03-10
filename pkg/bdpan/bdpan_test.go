package bdpan

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFileHash(t *testing.T) {
	testfile := "my-cat.webp"
	file, err := os.Open(testfile)
	assert.NoError(t, err)
	defer file.Close()
	name, md5, size, err := GetFileHash(file)
	assert.NoError(t, err)
	assert.Equal(t, "f4f395137a7443ea55e456d59d11ec56", md5)
	assert.Equal(t, testfile, name)
	assert.Equal(t, 281486, size)
}

func TestRapidUpload(t *testing.T) {
	bduss := os.Getenv("BDUSS")
	if bduss == "" {
		t.Skip("Skip test because BDUSS is not set")
	}
	err := RapidUpload(bduss, "/test/my-cat.webp", "f4f395137a7443ea55e456d59d11ec56", 281486)
	assert.NoError(t, err)
}

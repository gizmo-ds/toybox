package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandInt(t *testing.T) {
	v := RandInt(1, 10)
	assert.True(t, v >= 1 && v <= 10)
}

func TestMD5(t *testing.T) {
	assert.Equal(t, "e10adc3949ba59abbe56e057f20f883e",
		MD5("123456"))
	assert.Equal(t, "e10adc3949ba59abbe56e057f20f883e",
		MD5([]byte{49, 50, 51, 52, 53, 54}))
}

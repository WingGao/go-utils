package utils

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestBinExist(t *testing.T) {
	r := BinExist("vips", true)
	assert.True(t, r != "")
}

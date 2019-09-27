package ucore

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDesensitizePhone(t *testing.T) {
	assert.Equal(t, "186****9999", DesensitizePhone("18611119999"))
	assert.Equal(t, "186*****9999", DesensitizePhone("186111129999"))
	assert.Equal(t, "1861119999", DesensitizePhone("1861119999"))
}

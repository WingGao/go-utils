package utils

import (
	"testing"
	"path/filepath"
)

func TestExtIsImage(t *testing.T) {
	p := "/a/b/c.png"
	t.Log("filepath.Dir", filepath.Dir(p))
	t.Log("filepath.Base", filepath.Base(p))
	t.Log("filepath.Ext", filepath.Ext(p))
}

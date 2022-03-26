package uapp

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/ungerik/go-dry"
	"os"
	"path/filepath"
	"testing"
)

func TestKeepSingleProcess(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "testapp")
	os.Setenv("WING_BIN_PATH", filepath.Join(tempDir, "aa"))
	// 写入当前进程
	cPid := os.Getpid()
	key := "testkeep"
	pidPath := filepath.Join(tempDir, key+".pid")
	dry.FileSetString(pidPath, fmt.Sprintf("%d", cPid))
	assert.Error(t, KeepSingleProcess(key, false, false))
	assert.NoError(t, KeepSingleProcess(key, true, false))
}

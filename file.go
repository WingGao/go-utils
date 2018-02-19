package utils

import (
	"time"
	"github.com/ungerik/go-dry"
	"encoding/base64"
	"path/filepath"
	"os"
	"os/exec"
)

func FileGetBase64(filenameOrURL string, timeout ...time.Duration) (out string, err error) {
	bs, err1 := dry.FileGetBytes(filenameOrURL, timeout...)
	if err1 != nil {
		err = err1
		return
	}
	out = base64.StdEncoding.EncodeToString(bs)
	return
}

func BinPath() string {
	p, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return p
}

// 可执行文件是否存在
func BinExist(binfile string, canPainc bool) bool {
	cmd := exec.Command("which", binfile)
	out, err := cmd.Output()
	if err != nil || len(out) == 0 {
		if canPainc {
			panic("no such bin file: " + binfile)
		}
		return false
	}
	return true
}

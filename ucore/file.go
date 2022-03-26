package ucore

import (
	"encoding/base64"
	"fmt"
	"github.com/ungerik/go-dry"
	"os"
	"os/exec"
	"path/filepath"
	"time"
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

// 当前程序的完整路径
func BinPath() string {
	envp := os.Getenv("WING_BIN_PATH")
	if envp != "" {
		return envp
	}
	p, _ := filepath.Abs(os.Args[0])
	return p
}

// 当前程序所在目录的完整路径
func BinDirPath() string {
	p, _ := filepath.Abs(filepath.Dir(BinPath()))
	return p
}

// 可执行文件是否存在
func BinExist(binfile string, canPainc bool) string {
	ep, err := exec.LookPath(binfile)
	if err != nil || len(ep) == 0 {
		if canPainc {
			fmt.Println(err)
			panic("no such bin file: " + binfile)
		}
		return ""
	}
	return ep
}

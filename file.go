package utils

import (
	"time"
	"github.com/ungerik/go-dry"
	"encoding/base64"
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

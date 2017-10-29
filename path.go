package utils

import (
	"path"
)

func GetFileBaseName(s string) string {
	out := ""
	for i := len(s) - 1; i >= 0 && s[i] != '/'; i-- {
		if s[i] == '.' {
			out = s[:i]
			break
		}
	}
	return path.Base(out)
}

//ext需要带上'.', eg: '.doc'
func FileRenameExt(s, ext string) string {
	for i := len(s) - 1; i >= 0 && s[i] != '/'; i-- {
		if s[i] == '.' {
			return s[:i] + ext
		}
	}
	return s + ext
}

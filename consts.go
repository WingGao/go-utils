package utils

import "os"

const (
	XSESSION_KEY = "xsession"
)

var (
	DEBUG                        = true
	IMAGE_EXT_LIST               = []string{".jpg", ".jpeg", ".png"}
	DEFAULT_FILEMODE os.FileMode = 0776
)

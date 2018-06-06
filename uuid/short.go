package uuid

import "github.com/teris-io/shortid"

// 9位，无序
func Short() (id string, err error) {
	return shortid.Generate()
}

package utils

func init() {
	for _, v := range UtilsErrList {
		AddHandlerIgnoreErrors(v)
	}
}

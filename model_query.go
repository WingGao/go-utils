package utils

import (
	"fmt"
	"github.com/WingGao/go-utils/ucore"
)

func FillModelWithKey(srcItems interface{}, srcKey string, mod IModel, modKey string,
	setFunc func(srcItem interface{}, modItem interface{})) error {
	list := mod.MakePSlice()
	ids := ucore.MapGetColumn(srcItems, srcKey)
	if len(ids) <= 0 {
		return nil
	}
	if err := mod.Find(list, fmt.Sprintf("`%s` IN (?)", modKey), ids); err != nil {
		return err
	} else {
		ucore.ArrJoin(srcItems, list, srcKey, modKey, setFunc)
	}
	return nil
}

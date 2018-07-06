package elasticutil

import (
	"github.com/olivere/elastic"
	"github.com/WingGao/go-utils"
	"github.com/thoas/go-funk"
	"reflect"
	"context"
	"fmt"
)

type EsClient struct {
	elastic.Client
	Prefix string //Index前缀
}

func (c *EsClient) GetFullIndex(name string) string {
	return c.Prefix + "_" + name
}

func (c *EsClient) SyncMySQLModel(mod utils.IModel) (err error) {
	step := 500
	bs := c.Bulk()
	bs.Index(c.GetFullIndex(mod.GetTableName())).Type("doc")
	ctx := context.Background()
	for i := 0; ; i++ {
		items := mod.MakePSlice()
		err = mod.Table().Limit(step).Offset(i * step).Scan(items).Error
		if err != nil {
			return
		}
		funk.ForEach(reflect.ValueOf(items).Elem().Interface(), func(item interface{}) {
			itemi := item.(utils.IModel)
			itemi.SetDB(mod.GetDB())
			itemi.SetParent(item)
			req := elastic.NewBulkUpdateRequest().Id(fmt.Sprintf("%v", itemi.PrimaryKeyValue())).Doc(item).DocAsUpsert(true)
			bs.Add(req)
		})
		_, err1 := bs.Do(ctx)
		//fmt.Println(utils.JsonMarshalIndentString(rep, "", "    "))
		if err1 != nil {
			return err1
		}
		if utils.SizeOf(items) < step {
			break
		}
	}
	return
}

func (EsClient) NewContextBackground() context.Context {
	return context.Background()
}

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

func (c *EsClient) SyncMySQLModelAll(mod utils.IModel) (err error) {
	step := 500
	ctx := context.Background()
	//先清空
	indexName := c.GetFullIndex(mod.GetTableName())
	_, err = c.DeleteIndex(indexName).Do(ctx)
	if err != nil {
		return
	}

	//批量插入
	bs := c.Bulk()
	bs.Index(c.GetFullIndex(mod.GetTableName())).Type("doc")
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

func (c *EsClient) UpsertMySQLModelOne(mod utils.IModel) (err error) {
	ctx := context.Background()
	_, err = c.Update().Index(c.GetFullIndex(mod.GetTableName())).Type("doc").
		Id(fmt.Sprintf("%v", mod.PrimaryKeyValue())).Doc(mod).DocAsUpsert(true).
		Do(ctx)
	return
}

func (c *EsClient) DeleteMySQLModelOne(mod utils.IModel) (err error) {
	ctx := context.Background()
	_, err = c.Delete().Index(c.GetFullIndex(mod.GetTableName())).Type("doc").
		Id(fmt.Sprintf("%v", mod.PrimaryKeyValue())).
		Do(ctx)
	return
}

func (EsClient) NewContextBackground() context.Context {
	return context.Background()
}

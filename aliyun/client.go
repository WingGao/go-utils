package aliyun

import (

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/WingGao/go-utils"
)

var (
	DefaultClient *sdk.Client
)

func Init(cnf utils.MConfig) error {
	client, err := sdk.NewClientWithAccessKey(cnf.Aliyun.RegionId, cnf.Aliyun.Key, cnf.Aliyun.Secret)
	if err != nil {
		return err
	}
	DefaultClient = client
	return nil
}

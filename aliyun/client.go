package aliyun

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"wingao.net/webproj/core"
)

var (
	DefaultClient *sdk.Client
)

func Init(cnf core.MConfig) error {
	client, err := sdk.NewClientWithAccessKey(cnf.Aliyun.RegionId, cnf.Aliyun.Key, cnf.Aliyun.Secret)
	if err != nil {
		return err
	}
	DefaultClient = client
	return nil
}

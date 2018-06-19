package aliyun

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/rs/xid"
)

//人机验证
func AfsCheck() {
	req := requests.NewCommonRequest()
	req.Domain = "jaq.aliyuncs.com"
	req.Version = "2016-11-23"
	req.ApiName = "AfsCheck"
	req.QueryParams["Token"] = xid.New().String()

	DefaultClient.ProcessCommonRequest(req)
}

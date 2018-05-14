package qiniu

import (
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	//"github.com/go-errors/errors"
	"context"
	"fmt"
)

type Client struct {
	mac          *qbox.Mac
	bucket       string
	formUploader *storage.FormUploader
}

func NewClient(accessKey, secretKey, bucket string) (c *Client, err error) {
	client := &Client{
		mac:    qbox.NewMac(accessKey, secretKey),
		bucket: bucket,
	}
	cfg := &storage.Config{}
	client.formUploader = storage.NewFormUploader(cfg)
	c = client
	return
}

func (m *Client) Put(localFile, key string) (err error) {
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{
	}
	putPolicy := storage.PutPolicy{
		Scope: fmt.Sprintf("%s:%s", m.bucket, key), //允许覆盖
	}
	token := putPolicy.UploadToken(m.mac)
	err = m.formUploader.PutFile(context.Background(), &ret, token, key, localFile, &putExtra)
	return
}

func (m *Client) UpToken() string {
	return ""
}

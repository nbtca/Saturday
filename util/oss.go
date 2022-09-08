package util

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func Upload(name string, reader io.Reader) (string, error) {
	// Endpoint以杭州为例，其它Region请按实际情况填写。
	endpoint := "http://oss-cn-hangzhou.aliyuncs.com"
	// 阿里云账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM用户进行API访问或日常运维，请登录RAM控制台创建RAM用户。
	accessKeyId := os.Getenv("ACCESS_KEY_ID")
	accessKeySecret := os.Getenv("ACCESS_KEY_SECRET")
	bucketName := "sunday-res"
	// <yourObjectName>上传文件到OSS时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
	timestamp := time.Now().Unix()

	ext := filepath.Ext(name)
	objectName := fmt.Sprintf("weekend/%v%v", timestamp, ext)

	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		return "", err
	}
	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return "", err
	}
	// 上传文件。
	if err = bucket.PutObject(objectName, reader); err != nil {
		return "", err
	}
	return "https://sunday-res.oss-cn-hangzhou.aliyuncs.com/" + objectName, nil
}

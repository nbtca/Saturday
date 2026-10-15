package util

import (
	"crypto/rand"
	"fmt"
	"io"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/viper"
)

type PhotoStore interface {
	Put(key, contentType string, r io.Reader) error
	SignURL(key string) (string, error)
}

var Photos PhotoStore = ossPhotoStore{}

type ossPhotoStore struct{}

func (ossPhotoStore) bucket() (*oss.Bucket, error) {
	client, err := oss.New("https://oss-cn-hangzhou.aliyuncs.com",
		viper.GetString("aliyun.access_key_id"), viper.GetString("aliyun.access_key_secret"))
	if err != nil {
		return nil, err
	}
	return client.Bucket("sunday-res")
}

func (s ossPhotoStore) Put(key, contentType string, r io.Reader) error {
	bucket, err := s.bucket()
	if err != nil {
		return err
	}
	return bucket.PutObject(key, r, oss.ContentType(contentType), oss.ObjectACL(oss.ACLPrivate))
}

func (s ossPhotoStore) SignURL(key string) (string, error) {
	bucket, err := s.bucket()
	if err != nil {
		return "", err
	}
	return bucket.SignURL(key, oss.HTTPGet, 3600)
}

func PhotoKey(eventId int64, contentType string) string {
	ext := ".jpg"
	if contentType == "image/webp" {
		ext = ".webp"
	}
	return fmt.Sprintf("weekend/events/%d/%s%s", eventId, rand.Text(), ext)
}

package util

import (
	"crypto/rand"
	"fmt"
	"io"
	"sync"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/viper"
)

type PhotoStore interface {
	Put(key, contentType string, r io.Reader) error
	Delete(keys ...string) error
	SignURL(key string) (string, error)
}

var Photos PhotoStore = &ossPhotoStore{}

type ossPhotoStore struct {
	mu          sync.Mutex
	credentials string
	cached      *oss.Bucket
}

func (s *ossPhotoStore) bucket() (*oss.Bucket, error) {
	id, secret := viper.GetString("aliyun.access_key_id"), viper.GetString("aliyun.access_key_secret")
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cached != nil && s.credentials == id+":"+secret {
		return s.cached, nil
	}
	client, err := oss.New("https://oss-cn-hangzhou.aliyuncs.com", id, secret)
	if err != nil {
		return nil, err
	}
	bucket, err := client.Bucket("sunday-res")
	if err != nil {
		return nil, err
	}
	s.credentials, s.cached = id+":"+secret, bucket
	return bucket, nil
}

func (s *ossPhotoStore) Put(key, contentType string, r io.Reader) error {
	bucket, err := s.bucket()
	if err != nil {
		return err
	}
	return bucket.PutObject(key, r, oss.ContentType(contentType), oss.ObjectACL(oss.ACLPrivate))
}

func (s *ossPhotoStore) Delete(keys ...string) error {
	bucket, err := s.bucket()
	if err != nil {
		return err
	}
	_, err = bucket.DeleteObjects(keys, oss.DeleteObjectsQuiet(true))
	return err
}

func (s *ossPhotoStore) SignURL(key string) (string, error) {
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

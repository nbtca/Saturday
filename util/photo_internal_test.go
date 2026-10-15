package util

import (
	"testing"

	"github.com/spf13/viper"
)

func TestOSSPhotoStoreReusesBucketUntilCredentialsChange(t *testing.T) {
	viper.Set("aliyun.access_key_id", "id-a")
	viper.Set("aliyun.access_key_secret", "secret")
	store := &ossPhotoStore{}
	first, _ := store.bucket()
	second, _ := store.bucket()
	if first == nil || first != second {
		t.Fatal("bucket should be reused while credentials are unchanged")
	}
	viper.Set("aliyun.access_key_id", "id-b")
	third, _ := store.bucket()
	if third == first {
		t.Error("bucket should be rebuilt after credentials change")
	}
}

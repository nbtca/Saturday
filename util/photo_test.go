package util_test

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/nbtca/saturday/util"
)

func TestPhotoKey(t *testing.T) {
	webp := util.PhotoKey(42, "image/webp")
	if !regexp.MustCompile(`^weekend/events/42/[A-Z2-7]{26}\.webp$`).MatchString(webp) {
		t.Errorf("unexpected key %q", webp)
	}
	jpeg := util.PhotoKey(42, "image/jpeg")
	if !regexp.MustCompile(`^weekend/events/42/[A-Z2-7]{26}\.jpg$`).MatchString(jpeg) {
		t.Errorf("unexpected key %q", jpeg)
	}
	if util.PhotoKey(42, "image/webp") == webp {
		t.Error("keys must be unique")
	}
}

func TestPhotoErrorCode(t *testing.T) {
	err := fmt.Errorf("failed to store photos: %w", oss.ServiceError{Code: "AccessDenied"})
	if code := util.PhotoErrorCode(err); code != "AccessDenied" {
		t.Errorf("got %q, want AccessDenied", code)
	}
	if code := util.PhotoErrorCode(errors.New("connection reset")); code != "" {
		t.Errorf("got %q, want empty", code)
	}
}

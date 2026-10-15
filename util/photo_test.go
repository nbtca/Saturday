package util_test

import (
	"regexp"
	"testing"

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

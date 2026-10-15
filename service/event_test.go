package service_test

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/service"
	"github.com/nbtca/saturday/util"
)

var (
	jpegBytes = []byte("\xff\xd8\xff\xe0\x00\x10JFIF\x00")
	webpBytes = []byte("RIFF\x00\x00\x00\x00WEBPVP8 ")
	pngBytes  = []byte("\x89PNG\r\n\x1a\n")
)

func upload(b []byte) service.PhotoUpload {
	return service.PhotoUpload{File: bytes.NewReader(b), Size: int64(len(b))}
}

type memoryStore struct {
	puts    []string
	deleted []string
	failOn  int
}

func (m *memoryStore) Put(key, contentType string, r io.Reader) error {
	if len(m.puts)+1 == m.failOn {
		return errors.New("InvalidAccessKeyId")
	}
	m.puts = append(m.puts, key)
	return nil
}

func (m *memoryStore) Delete(keys ...string) error {
	m.deleted = append(m.deleted, keys...)
	return nil
}

func (m *memoryStore) SignURL(key string) (string, error) { return "https://signed/" + key, nil }

func TestValidatePhotoUpload(t *testing.T) {
	member := model.Identity{Id: "2333333333", Role: "member"}
	accepted := model.Event{EventId: 1, MemberId: "2333333333", Status: util.Accepted}

	types, err := service.ValidatePhotoUpload(accepted, member, []service.PhotoUpload{upload(jpegBytes), upload(webpBytes)})
	if err != nil || strings.Join(types, ",") != "image/jpeg,image/webp" {
		t.Fatalf("valid upload: types=%v err=%v", types, err)
	}

	committed := accepted
	committed.Status = util.Committed
	if _, err := service.ValidatePhotoUpload(committed, member, []service.PhotoUpload{upload(jpegBytes)}); err != nil {
		t.Errorf("committed event should accept photos: %v", err)
	}

	admin := model.Identity{Id: "1111111111", Role: "admin"}
	if _, err := service.ValidatePhotoUpload(accepted, admin, []service.PhotoUpload{upload(jpegBytes)}); err != nil {
		t.Errorf("admin should be allowed: %v", err)
	}

	cases := []struct {
		name     string
		event    model.Event
		identity model.Identity
		files    []service.PhotoUpload
	}{
		{"other member", accepted, model.Identity{Id: "1111111111", Role: "member"}, []service.PhotoUpload{upload(jpegBytes)}},
		{"client", accepted, model.Identity{Role: "client"}, []service.PhotoUpload{upload(jpegBytes)}},
		{"open event", model.Event{MemberId: "2333333333", Status: util.Open}, member, []service.PhotoUpload{upload(jpegBytes)}},
		{"closed event", model.Event{MemberId: "2333333333", Status: util.Closed}, member, []service.PhotoUpload{upload(jpegBytes)}},
		{"no files", accepted, member, nil},
		{"png", accepted, member, []service.PhotoUpload{upload(jpegBytes), upload(pngBytes)}},
		{"too large", accepted, member, []service.PhotoUpload{{File: bytes.NewReader(jpegBytes), Size: 5<<20 + 1}}},
		{"over 20", model.Event{MemberId: "2333333333", Status: util.Accepted, Photos: make([]model.EventPhoto, 20)}, member, []service.PhotoUpload{upload(jpegBytes)}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := service.ValidatePhotoUpload(c.event, c.identity, c.files); err == nil {
				t.Error("expected rejection")
			}
		})
	}
}

func TestAddPhotosRejectsWithoutWriting(t *testing.T) {
	store := &memoryStore{}
	util.Photos = store
	event := model.Event{EventId: 1, MemberId: "2333333333", Status: util.Accepted}
	err := service.EventServiceApp.AddPhotos(&event, model.Identity{Id: "2333333333", Role: "member"},
		[]service.PhotoUpload{upload(jpegBytes), upload(pngBytes)})
	if err == nil || len(store.puts) != 0 {
		t.Errorf("err=%v puts=%v", err, store.puts)
	}
}

func TestAddPhotosReportsStorageFailure(t *testing.T) {
	store := &memoryStore{failOn: 2}
	util.Photos = store
	event := model.Event{EventId: 1, MemberId: "2333333333", Status: util.Accepted}
	err := service.EventServiceApp.AddPhotos(&event, model.Identity{Id: "2333333333", Role: "member"},
		[]service.PhotoUpload{upload(jpegBytes), upload(webpBytes)})
	if !errors.Is(err, service.ErrPhotoStorage) {
		t.Errorf("err=%v, want ErrPhotoStorage", err)
	}
	if len(store.puts) != 1 || strings.Join(store.deleted, ",") != store.puts[0] {
		t.Errorf("uploaded %v, deleted %v", store.puts, store.deleted)
	}
}

func TestSignPhotos(t *testing.T) {
	util.Photos = &memoryStore{}
	event := model.Event{Photos: []model.EventPhoto{{ObjectKey: "weekend/events/1/a.webp"}}}
	if err := service.SignPhotos(&event); err != nil || event.Photos[0].Url != "https://signed/weekend/events/1/a.webp" {
		t.Errorf("err=%v url=%q", err, event.Photos[0].Url)
	}
}

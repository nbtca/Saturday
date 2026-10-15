//go:build darwin || linux

package router

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"runtime"
	"runtime/debug"
	"testing"

	"golang.org/x/sys/unix"
)

func countUploadTestDescriptors() int {
	count := 0
	for fd := 0; fd < 4096; fd++ {
		if _, err := unix.FcntlInt(uintptr(fd), unix.F_GETFD, 0); err == nil {
			count++
		}
	}
	return count
}

func TestRejectedUploadsReleaseResourcesWithoutGC(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("TMPDIR", tempDir)
	runtime.GC()
	previousGC := debug.SetGCPercent(-1)
	defer debug.SetGCPercent(previousGC)
	router := SetupRouter()

	for _, tc := range []struct {
		name        string
		path        string
		contentType string
		wantStatus  int
	}{
		{"unauthorized", "/member/events/1/photos", "image/jpeg", http.StatusUnauthorized},
		{"invalid path", "/member/events/invalid/photos", "image/jpeg", http.StatusUnprocessableEntity},
		{"invalid media type", "/member/events/1/photos", "image/png", http.StatusUnprocessableEntity},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var body bytes.Buffer
			form := multipart.NewWriter(&body)
			part, err := form.CreatePart(textproto.MIMEHeader{
				"Content-Disposition": {`form-data; name="photos"; filename="photo.jpg"`},
				"Content-Type":        {tc.contentType},
			})
			if err != nil {
				t.Fatal(err)
			}
			// Exceed Huma's 8 KiB memory threshold so file handles are real OS resources.
			if _, err := part.Write(bytes.Repeat([]byte{0xff}, 16<<10)); err != nil {
				t.Fatal(err)
			}
			if err := form.Close(); err != nil {
				t.Fatal(err)
			}
			before := countUploadTestDescriptors()
			for i := 0; i < 10; i++ {
				req := httptest.NewRequest(http.MethodPost, tc.path, bytes.NewReader(body.Bytes()))
				req.Header.Set("Authorization", "invalid")
				req.Header.Set("Content-Type", form.FormDataContentType())
				res := httptest.NewRecorder()
				router.ServeHTTP(res, req)
				// Also clean up on failure against an implementation lacking middleware cleanup.
				if req.MultipartForm != nil {
					t.Cleanup(func() { req.MultipartForm.RemoveAll() })
				}
				if res.Code != tc.wantStatus {
					t.Fatalf("status=%d body=%s", res.Code, res.Body.String())
				}
			}
			if after := countUploadTestDescriptors(); after != before {
				t.Errorf("open file descriptors: before=%d after=%d", before, after)
			}
			entries, err := os.ReadDir(tempDir)
			if err != nil {
				t.Fatal(err)
			}
			if len(entries) != 0 {
				t.Errorf("retained %d multipart temporary files", len(entries))
			}
		})
	}
}

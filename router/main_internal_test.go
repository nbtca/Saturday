package router

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strings"
	"testing"
)

func TestOversizedUploadIsRejectedBeforeAuth(t *testing.T) {
	body := &bytes.Buffer{}
	form := multipart.NewWriter(body)
	part, _ := form.CreatePart(textproto.MIMEHeader{
		"Content-Disposition": {`form-data; name="photos"; filename="a.jpg"`},
		"Content-Type":        {"image/jpeg"},
	})
	part.Write(bytes.Repeat([]byte{0xff}, maxRequestBytes+1))
	form.Close()

	req := httptest.NewRequest(http.MethodPost, "/member/events/1/photos", body)
	req.Header.Set("Content-Type", form.FormDataContentType())
	res := httptest.NewRecorder()
	SetupRouter().ServeHTTP(res, req)

	if res.Code == http.StatusUnauthorized || !strings.Contains(res.Body.String(), "too large") {
		t.Errorf("status=%d body=%s", res.Code, res.Body.String())
	}
}

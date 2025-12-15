package router

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/nbtca/saturday/model/dto"
	"github.com/nbtca/saturday/util"
)

type CommonRouter struct{}

const (
	MaxUploadSize = 10 << 20 // 10MB
)

func (CommonRouter) Upload(ctx context.Context, input *UploadFileInput) (*util.CommonResponse[dto.FileUploadResponse], error) {
	// Get the uploaded file from the parsed multipart form
	file := input.RawBody.Data().File

	// Validate file size
	if file.Size > MaxUploadSize {
		return nil, huma.Error400BadRequest(fmt.Sprintf("File size exceeds maximum allowed size of %d bytes", MaxUploadSize))
	}

	// Additional validation: check file signature (magic bytes)
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, huma.Error400BadRequest("Failed to read file header: " + err.Error())
	}

	// Detect content type from file content
	contentType := http.DetectContentType(buffer[:n])
	allowedTypes := []string{"image/jpeg", "image/png", "image/webp"}
	isAllowed := false
	for _, allowed := range allowedTypes {
		if contentType == allowed {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		return nil, huma.Error400BadRequest(fmt.Sprintf("Invalid image file type detected: %s. Allowed types: JPEG, PNG, WebP", contentType))
	}

	// Reset file pointer after signature check
	if _, err := file.Seek(0, 0); err != nil {
		return nil, huma.Error500InternalServerError("Failed to reset file pointer: " + err.Error())
	}

	// Upload to Aliyun OSS
	url, err := util.Upload(file.Filename, file)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to upload file: " + err.Error())
	}

	return &util.CommonResponse[dto.FileUploadResponse]{
		Body: dto.FileUploadResponse{
			Url: url,
		},
	}, nil
}

var CommonRouterApp = CommonRouter{}

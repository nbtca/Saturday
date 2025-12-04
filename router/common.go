package router

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/nbtca/saturday/model/dto"
	"github.com/nbtca/saturday/util"
)

type CommonRouter struct{}

const (
	MaxUploadSize     = 10 << 20 // 10MB
	MaxMemory         = 32 << 20 // 32MB for multipart parsing
	AllowedImageTypes = "image/jpeg,image/png,image/webp,image/jpg"
)

func (CommonRouter) Upload(ctx context.Context, input *UploadFileInput) (*util.CommonResponse[dto.FileUploadResponse], error) {
	// Get the underlying HTTP request from context
	req, ok := huma.ContextRequest(ctx)
	if !ok {
		return nil, huma.Error500InternalServerError("Failed to get HTTP request from context")
	}

	// Parse multipart form with size limit
	if err := req.ParseMultipartForm(MaxMemory); err != nil {
		return nil, huma.Error400BadRequest("Failed to parse multipart form: " + err.Error())
	}
	defer req.MultipartForm.RemoveAll()

	// Get the file from form
	file, header, err := req.FormFile("file")
	if err != nil {
		return nil, huma.Error400BadRequest("No file provided or invalid field name. Use 'file' as the field name")
	}
	defer file.Close()

	// Validate file size
	if header.Size > MaxUploadSize {
		return nil, huma.Error400BadRequest(fmt.Sprintf("File size exceeds maximum allowed size of %d bytes", MaxUploadSize))
	}

	// Validate content type
	contentType := header.Header.Get("Content-Type")
	if !isAllowedImageType(contentType) {
		return nil, huma.Error400BadRequest(fmt.Sprintf("Invalid file type. Allowed types: %s", AllowedImageTypes))
	}

	// Additional validation: check file signature (magic bytes)
	if err := validateImageSignature(file); err != nil {
		return nil, huma.Error400BadRequest("Invalid image file: " + err.Error())
	}

	// Reset file pointer after signature check
	if _, err := file.Seek(0, 0); err != nil {
		return nil, huma.Error500InternalServerError("Failed to reset file pointer")
	}

	// Upload to Aliyun OSS
	url, err := util.Upload(header.Filename, file)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to upload file: " + err.Error())
	}

	return &util.CommonResponse[dto.FileUploadResponse]{
		Data: dto.FileUploadResponse{
			Url: url,
		},
	}, nil
}

func isAllowedImageType(contentType string) bool {
	allowedTypes := strings.Split(AllowedImageTypes, ",")
	for _, allowedType := range allowedTypes {
		if strings.TrimSpace(allowedType) == contentType {
			return true
		}
	}
	return false
}

func validateImageSignature(file multipart.File) error {
	// Read first 512 bytes to detect content type
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read file header: %w", err)
	}

	// Detect content type from file content
	contentType := http.DetectContentType(buffer[:n])

	// Validate detected type matches allowed types
	if !isAllowedImageType(contentType) {
		return fmt.Errorf("file content does not match allowed image types (detected: %s)", contentType)
	}

	return nil
}

var CommonRouterApp = CommonRouter{}

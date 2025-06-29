package router

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/nbtca/saturday/model/dto"
	"github.com/nbtca/saturday/util"
)

type CommonRouter struct{}

func (CommonRouter) Upload(ctx context.Context, input *UploadFileInput) (*util.CommonResponse[dto.FileUploadResponse], error) {
	// TODO: Implement multipart file upload with Huma
	// For now, this needs special handling since Huma's multipart support may require custom implementation
	// This endpoint may need to remain as Gin for now until Huma multipart is properly implemented
	return nil, huma.Error501NotImplemented("Upload endpoint migration pending - requires multipart support")
}

var CommonRouterApp = CommonRouter{}

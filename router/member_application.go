package router

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/nbtca/saturday/middleware"
	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/service"
	"github.com/nbtca/saturday/util"
)

type MemberApplicationRouter struct{}

// Response types
type MemberApplicationResponse struct {
	util.CommonResponse[model.MemberApplication]
}

type MemberApplicationsListResponse struct {
	Success    bool                        `json:"success"`
	Result     []model.MemberApplication   `json:"result"`
	TotalCount int                         `json:"totalCount"`
}

// SubmitApplication creates a new member application (public endpoint)
func (MemberApplicationRouter) SubmitApplication(ctx context.Context, input *SubmitMemberApplicationInput) (*util.CommonResponse[model.MemberApplication], error) {
	application := &model.MemberApplication{
		MemberId: input.Body.MemberId,
		Name:     input.Body.Name,
		Phone:    input.Body.Phone,
		Section:  input.Body.Section,
		QQ:       input.Body.QQ,
		Email:    input.Body.Email,
		Major:    input.Body.Major,
		Class:    input.Body.Class,
		Memo:     input.Body.Memo,
	}

	err := service.MemberApplicationServiceApp.SubmitApplication(application)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	return util.MakeCommonResponse(*application), nil
}

// GetApplications retrieves member applications (admin only)
func (MemberApplicationRouter) GetApplications(ctx context.Context, input *GetMemberApplicationsInput) (*MemberApplicationsListResponse, error) {
	_, err := middleware.AuthenticateUser(input.Authorization, "admin")
	if err != nil {
		return nil, err
	}

	applications, count, err := service.MemberApplicationServiceApp.GetApplications(
		input.Offset,
		input.Limit,
		input.Status,
		input.Search,
	)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	return &MemberApplicationsListResponse{
		Success:    true,
		Result:     applications,
		TotalCount: count,
	}, nil
}

// GetApplicationById retrieves a single application by ID (admin only)
func (MemberApplicationRouter) GetApplicationById(ctx context.Context, input *GetMemberApplicationByIdInput) (*util.CommonResponse[model.MemberApplication], error) {
	_, err := middleware.AuthenticateUser(input.Authorization, "admin")
	if err != nil {
		return nil, err
	}

	application, err := service.MemberApplicationServiceApp.GetApplicationById(input.ApplicationId)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	return util.MakeCommonResponse(application), nil
}

// ApproveApplication approves an application and creates the member (admin only)
func (MemberApplicationRouter) ApproveApplication(ctx context.Context, input *ApproveMemberApplicationInput) (*util.CommonResponse[string], error) {
	auth, err := middleware.AuthenticateUser(input.Authorization, "admin")
	if err != nil {
		return nil, err
	}

	err = service.MemberApplicationServiceApp.ApproveApplication(input.ApplicationId, auth.ID)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	return util.MakeCommonResponse("Application approved and member created successfully"), nil
}

// RejectApplication rejects an application (admin only)
func (MemberApplicationRouter) RejectApplication(ctx context.Context, input *RejectMemberApplicationInput) (*util.CommonResponse[string], error) {
	auth, err := middleware.AuthenticateUser(input.Authorization, "admin")
	if err != nil {
		return nil, err
	}

	err = service.MemberApplicationServiceApp.RejectApplication(
		input.ApplicationId,
		auth.ID,
		input.Body.Reason,
	)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	return util.MakeCommonResponse("Application rejected successfully"), nil
}

var MemberApplicationRouterApp = new(MemberApplicationRouter)

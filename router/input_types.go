package router

import (
	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/model/dto"
)

// Base input types for common parameters

type MemberPathInput struct {
	MemberId string `path:"MemberId" maxLength:"10" example:"2333333333" doc:"Member ID"`
}

type EventPathInput struct {
	EventId int64 `path:"EventId" example:"123" doc:"Event ID"`
}

type ClientPathInput struct {
	ClientId int64 `path:"ClientId" example:"456" doc:"Client ID"`
}

// Authentication input types

type AuthenticatedInput struct {
	Authorization string `header:"Authorization" doc:"Bearer token or JWT token"`
}

type MemberAuthInput struct {
	AuthenticatedInput
}

type AdminAuthInput struct {
	AuthenticatedInput
}

type ClientAuthInput struct {
	AuthenticatedInput
}

// Member endpoint inputs

type GetMemberInput struct {
	MemberAuthInput
}

type UpdateMemberInput struct {
	MemberAuthInput
	Body dto.UpdateMemberRequest
}

type UpdateMemberAvatarInput struct {
	MemberAuthInput
	Body struct {
		Avatar string `json:"avatar" doc:"Avatar URL"`
	}
}

type CreateMemberInput struct {
	AdminAuthInput
	MemberPathInput
	Body dto.CreateMemberRequest
}

type CreateManyMembersInput struct {
	AdminAuthInput
	Body []dto.CreateMemberRequest
}

type UpdateMemberBasicInput struct {
	AdminAuthInput
	MemberPathInput
	Body dto.UpdateMemberBasicRequest
}

type GetMemberByPageInput struct {
	AdminAuthInput
	dto.PageRequest
}

type ActivateMemberInput struct {
	MemberAuthInput
	Body dto.ActivateMemberRequest
}

type UpdateNotificationPreferencesInput struct {
	MemberAuthInput
	Body struct {
		Preferences []struct {
			NotificationType model.NotificationType `json:"notificationType"`
			Enabled          bool                   `json:"enabled"`
		} `json:"preferences"`
	}
}

// Event endpoint inputs

type GetMemberEventByPageInput struct {
	MemberAuthInput
	dto.PageRequest
	Status string `query:"status"`
	Order  string `query:"order" default:"ASC"`
}

type GetEventByIdInput struct {
	MemberAuthInput
	EventPathInput
}

type AcceptEventInput struct {
	MemberAuthInput
	EventPathInput
}

type DropEventInput struct {
	MemberAuthInput
	EventPathInput
}

type CommitEventInput struct {
	MemberAuthInput
	EventPathInput
	Body struct {
		Content string `json:"content"`
		Size    string `json:"size" required:"false"`
	}
}

type AlterCommitEventInput struct {
	MemberAuthInput
	EventPathInput
	Body struct {
		Content string `json:"content"`
		Size    string `json:"size" required:"false"`
	}
}

type RejectCommitEventInput struct {
	AdminAuthInput
	EventPathInput
}

type CloseEventInput struct {
	AdminAuthInput
	EventPathInput
}

type ExportEventsToXlsxInput struct {
	AdminAuthInput
	Status    string `query:"status"`
	Order     string `query:"order" default:"ASC"`
	StartTime string `query:"start_time" required:"true"`
	EndTime   string `query:"end_time" required:"true"`
}

// Client event endpoint inputs

type GetClientEventByIdInput struct {
	ClientAuthInput
	EventPathInput
}

type GetClientEventByPageInput struct {
	ClientAuthInput
	dto.PageRequest
	Status string `query:"status"`
	Order  string `query:"order" default:"ASC"`
}

type CreateClientEventInput struct {
	ClientAuthInput
	Body struct {
		Model             string `json:"model" required:"false"`
		Phone             string `json:"phone"`
		QQ                string `json:"qq" required:"false"`
		ContactPreference string `json:"contactPreference" required:"false"`
		Problem           string `json:"problem"`
	}
}

type UpdateClientEventInput struct {
	ClientAuthInput
	EventPathInput
	Body struct {
		Model             string `json:"model" required:"false"`
		Phone             string `json:"phone"`
		QQ                string `json:"qq" required:"false"`
		ContactPreference string `json:"contactPreference" required:"false"`
		Problem           string `json:"problem"`
		Size              string `json:"size" required:"false"`
	}
}

type CancelClientEventInput struct {
	ClientAuthInput
	EventPathInput
}

// Client auth endpoint inputs

type CreateTokenViaLogtoInput struct {
	ClientAuthInput
}

// Upload endpoint input

type UploadFileInput struct {
	AuthenticatedInput
	// File upload will be handled specially
}

// Webhook inputs (these may stay as Gin handlers)

type GithubWebhookInput struct {
	// Special handling for webhooks
}

type LogtoWebhookInput struct {
	// Special handling for webhooks
}

// Member Application endpoint inputs

type SubmitMemberApplicationInput struct {
	Body struct {
		MemberId string `json:"memberId" minLength:"1" maxLength:"10" doc:"Student ID" example:"2333333333"`
		Name     string `json:"name" minLength:"1" maxLength:"100" doc:"Full name" example:"张三"`
		Phone    string `json:"phone" minLength:"1" maxLength:"20" doc:"Phone number" example:"13800138000"`
		Section  string `json:"section" minLength:"1" maxLength:"50" doc:"Department/Section" example:"web"`
		QQ       string `json:"qq" maxLength:"20" required:"false" doc:"QQ number" example:"123456789"`
		Email    string `json:"email" maxLength:"100" required:"false" doc:"Email address" example:"test@example.com"`
		Major    string `json:"major" maxLength:"100" required:"false" doc:"Major" example:"计算机科学与技术"`
		Class    string `json:"class" maxLength:"50" required:"false" doc:"Class" example:"计算机196"`
		Memo     string `json:"memo" required:"false" doc:"Self introduction/Memo"`
	}
}

type GetMemberApplicationsInput struct {
	AdminAuthInput
	dto.PageRequest
	Status string `query:"status" required:"false" doc:"Filter by status (pending, approved, rejected)"`
	Search string `query:"search" required:"false" doc:"Search by name or member ID"`
}

type GetMemberApplicationByIdInput struct {
	AdminAuthInput
	ApplicationId string `path:"ApplicationId" doc:"Application ID"`
}

type ApproveMemberApplicationInput struct {
	AdminAuthInput
	ApplicationId string `path:"ApplicationId" doc:"Application ID"`
}

type RejectMemberApplicationInput struct {
	AdminAuthInput
	ApplicationId string `path:"ApplicationId" doc:"Application ID"`
	Body          struct {
		Reason string `json:"reason" required:"false" doc:"Rejection reason"`
	}
}

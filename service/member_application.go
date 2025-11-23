package service

import (
	"net/http"

	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/repo"
	"github.com/nbtca/saturday/util"
)

type MemberApplicationService struct{}

// SubmitApplication creates a new member application
func (service *MemberApplicationService) SubmitApplication(app *model.MemberApplication) error {
	// Validate required fields
	if app.MemberId == "" {
		return util.MakeServiceError(http.StatusUnprocessableEntity).
			SetMessage("Validation Failed").
			AddDetailError("memberApplication", "memberId", "memberId is required")
	}
	if app.Name == "" {
		return util.MakeServiceError(http.StatusUnprocessableEntity).
			SetMessage("Validation Failed").
			AddDetailError("memberApplication", "name", "name is required")
	}
	if app.Phone == "" {
		return util.MakeServiceError(http.StatusUnprocessableEntity).
			SetMessage("Validation Failed").
			AddDetailError("memberApplication", "phone", "phone is required")
	}
	if app.Section == "" {
		return util.MakeServiceError(http.StatusUnprocessableEntity).
			SetMessage("Validation Failed").
			AddDetailError("memberApplication", "section", "section is required")
	}

	// Check if member already exists
	exist, err := repo.ExistMember(app.MemberId)
	if err != nil {
		return err
	}
	if exist {
		return util.MakeServiceError(http.StatusUnprocessableEntity).
			SetMessage("Validation Failed").
			AddDetailError("memberApplication", "memberId", "member already exists with this ID")
	}

	// Check if there's already a pending application for this member ID
	existingApps, err := repo.GetMemberApplicationByMemberId(app.MemberId)
	if err != nil {
		return err
	}
	for _, existingApp := range existingApps {
		if existingApp.Status == string(model.ApplicationStatusPending) {
			return util.MakeServiceError(http.StatusUnprocessableEntity).
				SetMessage("Validation Failed").
				AddDetailError("memberApplication", "memberId", "pending application already exists for this member ID")
		}
	}

	// Create the application
	if err := repo.CreateMemberApplication(app); err != nil {
		return err
	}

	return nil
}

// GetApplications retrieves member applications with pagination and filters
func (service *MemberApplicationService) GetApplications(offset uint64, limit uint64, status string, search string) ([]model.MemberApplication, int, error) {
	applications, err := repo.GetMemberApplications(offset, limit, status, search)
	if err != nil {
		return nil, 0, err
	}

	count, err := repo.CountMemberApplications(status, search)
	if err != nil {
		return nil, 0, err
	}

	return applications, count, nil
}

// GetApplicationById retrieves a single application
func (service *MemberApplicationService) GetApplicationById(applicationId string) (model.MemberApplication, error) {
	application, err := repo.GetMemberApplicationById(applicationId)
	if err != nil {
		return model.MemberApplication{}, err
	}

	if application.ApplicationId == "" {
		return model.MemberApplication{}, util.MakeServiceError(http.StatusNotFound).
			SetMessage("Application not found")
	}

	return application, nil
}

// ApproveApplication approves an application and creates the member
func (service *MemberApplicationService) ApproveApplication(applicationId string, reviewedBy string) error {
	// Get the application
	application, err := service.GetApplicationById(applicationId)
	if err != nil {
		return err
	}

	// Check if already reviewed
	if application.Status != string(model.ApplicationStatusPending) {
		return util.MakeServiceError(http.StatusUnprocessableEntity).
			SetMessage("Validation Failed").
			AddDetailError("memberApplication", "status", "application already reviewed")
	}

	// Check if member already exists (in case created between submission and approval)
	exist, err := repo.ExistMember(application.MemberId)
	if err != nil {
		return err
	}
	if exist {
		return util.MakeServiceError(http.StatusUnprocessableEntity).
			SetMessage("Validation Failed").
			AddDetailError("memberApplication", "memberId", "member already exists with this ID")
	}

	// Create the member
	member := &model.Member{
		MemberId:  application.MemberId,
		Name:      application.Name,
		Phone:     application.Phone,
		Section:   application.Section,
		QQ:        application.QQ,
		Email:     "",     // Email field doesn't exist in Member model based on the schema
		Alias:     "",     // Will be set later by member
		Avatar:    "",     // Will be set later by member
		Profile:   application.Memo,
		Role:      "member",
		LogtoId:   "",     // Will be set when member activates
		CreatedBy: reviewedBy,
	}

	// Create member using existing service
	if err := MemberServiceApp.CreateMember(member); err != nil {
		return err
	}

	// Approve the application
	if err := repo.ApproveMemberApplication(applicationId, reviewedBy); err != nil {
		return err
	}

	return nil
}

// RejectApplication rejects an application
func (service *MemberApplicationService) RejectApplication(applicationId string, reviewedBy string, reason string) error {
	// Get the application
	application, err := service.GetApplicationById(applicationId)
	if err != nil {
		return err
	}

	// Check if already reviewed
	if application.Status != string(model.ApplicationStatusPending) {
		return util.MakeServiceError(http.StatusUnprocessableEntity).
			SetMessage("Validation Failed").
			AddDetailError("memberApplication", "status", "application already reviewed")
	}

	// Reject the application
	if err := repo.RejectMemberApplication(applicationId, reviewedBy, reason); err != nil {
		return err
	}

	return nil
}

var MemberApplicationServiceApp = new(MemberApplicationService)

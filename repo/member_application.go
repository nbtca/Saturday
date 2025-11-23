package repo

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/util"
)

// CreateMemberApplication creates a new member application
func CreateMemberApplication(app *model.MemberApplication) error {
	// Generate UUID for application ID
	app.ApplicationId = uuid.New().String()
	app.Status = string(model.ApplicationStatusPending)
	app.GmtCreate = util.GetDate()
	app.GmtModified = util.GetDate()

	sqlStmt, args, _ := sq.Insert("member_applications").Columns(
		"application_id", "member_id", "name", "phone", "section",
		"qq", "email", "major", "class", "memo", "status",
		"gmt_create", "gmt_modified").Values(
		app.ApplicationId, app.MemberId, app.Name, app.Phone, app.Section,
		app.QQ, app.Email, app.Major, app.Class, app.Memo, app.Status,
		app.GmtCreate, app.GmtModified).ToSql()

	if _, err := db.Exec(sqlStmt, args...); err != nil {
		return err
	}

	return nil
}

// GetMemberApplications retrieves member applications with optional filters
func GetMemberApplications(offset uint64, limit uint64, status string, search string) ([]model.MemberApplication, error) {
	query := sq.Select("*").From("member_applications")

	// Apply status filter if provided
	if status != "" {
		// Support comma-separated status values
		query = query.Where(squirrel.Eq{"status": status})
	}

	// Apply search filter if provided
	if search != "" {
		query = query.Where(squirrel.Or{
			squirrel.Like{"name": "%" + search + "%"},
			squirrel.Like{"member_id": "%" + search + "%"},
		})
	}

	// Order by creation time descending (newest first)
	query = query.OrderBy("gmt_create DESC").Offset(offset).Limit(limit)

	sqlStmt, args, _ := query.ToSql()

	applications := []model.MemberApplication{}
	if err := db.Select(&applications, sqlStmt, args...); err != nil {
		return []model.MemberApplication{}, err
	}

	return applications, nil
}

// GetMemberApplicationById retrieves a single application by ID
func GetMemberApplicationById(applicationId string) (model.MemberApplication, error) {
	sqlStmt, args, _ := sq.Select("*").From("member_applications").
		Where(squirrel.Eq{"application_id": applicationId}).ToSql()

	application := model.MemberApplication{}
	if err := db.Get(&application, sqlStmt, args...); err != nil {
		if err == sql.ErrNoRows {
			return model.MemberApplication{}, nil
		}
		return model.MemberApplication{}, err
	}

	return application, nil
}

// GetMemberApplicationByMemberId retrieves applications by member ID
func GetMemberApplicationByMemberId(memberId string) ([]model.MemberApplication, error) {
	sqlStmt, args, _ := sq.Select("*").From("member_applications").
		Where(squirrel.Eq{"member_id": memberId}).
		OrderBy("gmt_create DESC").ToSql()

	applications := []model.MemberApplication{}
	if err := db.Select(&applications, sqlStmt, args...); err != nil {
		return []model.MemberApplication{}, err
	}

	return applications, nil
}

// CountMemberApplications counts total applications with optional filters
func CountMemberApplications(status string, search string) (int, error) {
	query := sq.Select("COUNT(*)").From("member_applications")

	// Apply status filter if provided
	if status != "" {
		query = query.Where(squirrel.Eq{"status": status})
	}

	// Apply search filter if provided
	if search != "" {
		query = query.Where(squirrel.Or{
			squirrel.Like{"name": "%" + search + "%"},
			squirrel.Like{"member_id": "%" + search + "%"},
		})
	}

	sqlStmt, args, _ := query.ToSql()

	var count int
	if err := db.Get(&count, sqlStmt, args...); err != nil {
		return 0, err
	}

	return count, nil
}

// ApproveMemberApplication approves an application
func ApproveMemberApplication(applicationId string, reviewedBy string) error {
	sqlStmt, args, _ := sq.Update("member_applications").
		Set("status", string(model.ApplicationStatusApproved)).
		Set("reviewed_by", reviewedBy).
		Set("reviewed_at", util.GetDate()).
		Set("gmt_modified", util.GetDate()).
		Where(squirrel.Eq{"application_id": applicationId}).ToSql()

	if _, err := db.Exec(sqlStmt, args...); err != nil {
		return err
	}

	return nil
}

// RejectMemberApplication rejects an application
func RejectMemberApplication(applicationId string, reviewedBy string, reason string) error {
	sqlStmt, args, _ := sq.Update("member_applications").
		Set("status", string(model.ApplicationStatusRejected)).
		Set("reviewed_by", reviewedBy).
		Set("reviewed_at", util.GetDate()).
		Set("reject_reason", reason).
		Set("gmt_modified", util.GetDate()).
		Where(squirrel.Eq{"application_id": applicationId}).ToSql()

	if _, err := db.Exec(sqlStmt, args...); err != nil {
		return err
	}

	return nil
}

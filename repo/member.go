package repo

import (
	"database/sql"

	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/util"

	"github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
)

var memberFields = []string{"member_id", "alias", "password", "name", "section", "role", "profile",
	"phone", "qq", "avatar", "created_by", "gmt_create", "gmt_modified"}

func getMemberStatement() squirrel.SelectBuilder {
	return sq.Select("*").From("member_view")
}

func ExistMember(id string) (bool, error) {
	var count int
	sql, args, _ := sq.Select("count(*) as count").From("member").Where(squirrel.Eq{"member_id": id}).ToSql()
	err := db.Get(&count, sql, args...)
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func GetMemberById(id string) (model.Member, error) {
	statement, args, _ := getMemberStatement().Where(squirrel.Eq{"member_id": id}).ToSql()
	member := model.Member{}
	if err := db.Get(&member, statement, args...); err != nil {
		if err == sql.ErrNoRows {
			return model.Member{}, nil
		}
		return model.Member{}, err
	}
	return member, nil
}

func GetMemberIdByLogtoId(logtoId string) (sql.NullString, error) {
	var memberId sql.NullString
	s, args, _ := sq.Select("member_id").From("member").Where(squirrel.Eq{"logto_id": logtoId}).ToSql()
	err := db.Get(&memberId, s, args...)
	if err != nil {
		return sql.NullString{}, err
	}
	return memberId, nil
}

func GetMemberByLogtoId(logtoId string) (model.Member, error) {
	statement, args, _ := getMemberStatement().Where(squirrel.Eq{"logto_id": logtoId}).ToSql()
	member := model.Member{}
	if err := db.Get(&member, statement, args...); err != nil {
		if err == sql.ErrNoRows {
			return model.Member{}, nil
		}
		return model.Member{}, err
	}
	return member, nil
}
func GetMemberByGithubId(githubId string) (model.Member, error) {
	statement, args, _ := getMemberStatement().Where(squirrel.Eq{"github_id": githubId}).ToSql()
	member := model.Member{}
	if err := db.Get(&member, statement, args...); err != nil {
		if err == sql.ErrNoRows {
			return model.Member{}, nil
		}
		return model.Member{}, err
	}
	return member, nil
}

func GetMembers(offset uint64, limit uint64) ([]model.Member, error) {
	sql, args, _ := getMemberStatement().Offset(offset).Limit(limit).ToSql()
	members := []model.Member{}
	if err := db.Select(&members, sql, args...); err != nil {
		return []model.Member{}, err
	}
	return members, nil
}

func CreateMember(member *model.Member) error {
	member.GmtCreate = util.GetDate()
	member.GmtModified = util.GetDate()
	sqlMember, argsMember, _ := sq.Insert("member").Columns(
		"member_id", "logto_id", "alias", "name", "section", "profile", "avatar",
		"phone", "qq", "created_by", "gmt_create", "gmt_modified").Values(
		member.MemberId, member.LogtoId, member.Alias, member.Name, member.Section,
		member.Profile, member.Avatar, member.Phone, member.QQ, member.CreatedBy,
		member.GmtCreate, member.GmtModified).ToSql()
	conn, err := db.Beginx()
	if err != nil {
		return err
	}
	defer util.RollbackOnErr(err, conn)
	if _, err = conn.Exec(sqlMember, argsMember...); err != nil {
		return err
	}
	if err = SetMemberRole(member.MemberId, member.Role, conn); err != nil {
		return err
	}
	if err = conn.Commit(); err != nil {
		return err
	}
	return nil
}

func UpdateMember(member model.Member) error {
	sql, args, _ := sq.Update("member").
		Set("logto_id", member.LogtoId).
		Set("github_id", member.GithubId).
		Set("alias", member.Alias).
		Set("name", member.Name).
		Set("section", member.Section).
		Set("password", member.Password).
		Set("profile", member.Profile).
		Set("phone", member.Phone).
		Set("qq", member.QQ).
		Set("avatar", member.Avatar).
		Set("gmt_modified", util.GetDate()).
		Where(squirrel.Eq{"member_id": member.MemberId}).ToSql()
	conn, err := db.Beginx()
	if err != nil {
		return err
	}
	defer util.RollbackOnErr(err, conn)
	if _, err = conn.Exec(sql, args...); err != nil {
		return err
	}
	if err = SetMemberRole(member.MemberId, member.Role, conn); err != nil {
		return err
	}
	if err = conn.Commit(); err != nil {
		return err
	}
	return nil
}

// UpdateNotificationPreferences updates the notification preferences for a member
func UpdateNotificationPreferences(memberId string, preferences model.NotificationPreferences) error {
	prefsJSON, err := preferences.Value()
	if err != nil {
		return err
	}

	sql, args, _ := sq.Update("member").
		Set("notification_preferences", prefsJSON).
		Set("gmt_modified", util.GetDate()).
		Where(squirrel.Eq{"member_id": memberId}).ToSql()

	if _, err = db.Exec(sql, args...); err != nil {
		return err
	}
	return nil
}

// GetMembersWithNotificationEnabled returns all members who have enabled a specific notification type
func GetMembersWithNotificationEnabled(notifType model.NotificationType) ([]model.Member, error) {
	// Query members where the notification_preferences JSONB contains the notification type set to true
	jsonPath := string(notifType)
	query := `
		SELECT * FROM member_view
		WHERE notification_preferences->$1 = 'true'::jsonb
		AND role IN ('member', 'admin')
	`

	members := []model.Member{}
	if err := db.Select(&members, query, jsonPath); err != nil {
		return []model.Member{}, err
	}
	return members, nil
}

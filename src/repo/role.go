package repo

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func ExistRole(role string) (bool, error) {
	var count int
	err := db.Get(&count, "SELECT count(*) as count FROM role where role = ?", role)
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func GetRoleId(role string) (sql.NullInt64, error) {
	var id sql.NullInt64
	err := db.Get(&id, "SELECT role_id FROM role where role = ?", role)
	if err != nil {
		return sql.NullInt64{}, err
	}
	return id, nil
}

func SetMemberRole(memberId string, role string, conn *sql.Tx) {
	sql := `INSERT INTO member_role_relation (member_id, role_id)
	 		VALUES (?, Select role_id from role where role = ?)`
	conn.Exec(sql, memberId, role)
}

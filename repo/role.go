package repo

import (
	"github.com/jmoiron/sqlx"
)

func ExistRole(role string) (bool, error) {
	var count int
	err := db.Get(&count, db.Rebind("SELECT count(*) as count FROM role where role = ?"), role)
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func SetMemberRole(memberId string, role string, conn *sqlx.Tx) error {
	sql := `INSERT INTO member_role_relation (member_id, role_id)
	 		VALUES (?, (Select role_id from role where role = ?))
			 ON DUPLICATE KEY UPDATE role_id=(Select role_id from role where role = ?)`
	if db.DriverName() == "pqHooked" {
		sql = `INSERT INTO member_role_relation (member_id, role_id)
	 		VALUES ($1, (Select role_id from role where role = $2))
			 ON CONFLICT (member_id) DO UPDATE SET role_id=(Select role_id from role where role = $3)`
	}
	_, err := conn.Exec(db.Rebind(sql), memberId, role, role)
	return err
}

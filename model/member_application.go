package model

// MemberApplication represents a new member application
type MemberApplication struct {
	ApplicationId string `json:"applicationId" db:"application_id"`
	MemberId      string `json:"memberId" db:"member_id"`
	Name          string `json:"name" db:"name"`
	Phone         string `json:"phone" db:"phone"`
	Section       string `json:"section" db:"section"`
	QQ            string `json:"qq" db:"qq"`
	Email         string `json:"email" db:"email"`
	Major         string `json:"major" db:"major"`
	Class         string `json:"class" db:"class"`
	Memo          string `json:"memo" db:"memo"`

	Status string `json:"status" db:"status"`

	ReviewedBy   string `json:"reviewedBy" db:"reviewed_by"`
	ReviewedAt   string `json:"reviewedAt" db:"reviewed_at"`
	RejectReason string `json:"rejectReason" db:"reject_reason"`

	GmtCreate   string `json:"gmtCreate" db:"gmt_create"`
	GmtModified string `json:"gmtModified" db:"gmt_modified"`
}

// MemberApplicationStatus represents the status of a member application
type MemberApplicationStatus string

const (
	ApplicationStatusPending  MemberApplicationStatus = "pending"
	ApplicationStatusApproved MemberApplicationStatus = "approved"
	ApplicationStatusRejected MemberApplicationStatus = "rejected"
)

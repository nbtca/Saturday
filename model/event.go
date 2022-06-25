package model

type Event struct {
	EventId           int64      `json:"event_id" db:"event_id"`
	ClientId          int64      `json:"client_id" db:"client_id"`
	Model             string     `json:"model"`
	Phone             string     `json:"phone"`
	QQ                string     `json:"qq"`
	ContactPreference string     `json:"contact_preference" db:"contact_preference" `
	Problem           string     `json:"problem" db:"problem"`
	MemberId          string     `json:"member_id" db:"member_id"`
	Member            *PublicMember    `json:"member" db:"-"`
	ClosedBy          string     `json:"closed_by_id" db:"closed_by"`
	ClosedByMember    *PublicMember    `json:"closed_by" db:"-"`
	Status            string     `json:"status"`
	Logs              []EventLog `json:"logs"`
	GmtCreate         string     `json:"gmt_create" db:"gmt_create"`
	GmtModified       string     `json:"gmt_modified" db:"gmt_modified"`
}

type Status struct {
	StatusId int64  `json:"status_id"`
	Status   string `json:"status"`
}
type EventEventStatusRelation struct {
	EventStatusId int64 `json:"event_status_id"`
	EventId       int64 `json:"event_id"`
}

type EventLog struct {
	EventLogId  int64  `json:"log_id" db:"event_log_id"`
	EventId     int64  `json:"-" db:"-"`
	Description string `json:"description"`
	MemberId    string `json:"member_id" db:"member_id"`
	Action      string `json:"action"`
	GmtCreate   string `json:"gmt_create" db:"gmt_create"`
}

type EventActionRelation struct {
	EventLogId    int64 `json:"event_log_id"`
	EventActionId int64 `json:"event_action_id"`
}

type EventAction struct {
	EventActionId int64  `json:"event_action_id"`
	Action        string `json:"action"`
}

type PublicEvent struct {
	EventId        int64      `json:"event_id" db:"event_id"`
	ClientId       int64      `json:"client_id" db:"client_id"`
	Model          string     `json:"model"`
	Problem        string     `json:"problem" db:"event_description"`
	MemberId       string     `json:"-" db:"member_id"`
	Member         *PublicMember    `json:"member"`
	ClosedBy       string     `json:"-" db:"closed_by"`
	ClosedByMember *PublicMember    `json:"closed_by"`
	Status         string     `json:"status"`
	Logs           []EventLog `json:"logs"`
	GmtCreate      string     `json:"gmt_create" db:"gmt_create"`
	GmtModified    string     `json:"gmt_modified" db:"gmt_modified"`
}

func CreatePublicEvent(e Event) PublicEvent {
	return PublicEvent{
		EventId:        e.EventId,
		ClientId:       e.ClientId,
		Model:          e.Model,
		Problem:        e.Problem,
		MemberId:       e.MemberId,
		Member:         e.Member,
		ClosedBy:       e.ClosedBy,
		ClosedByMember: e.ClosedByMember,
		Status:         e.Status,
		Logs:           e.Logs,
		GmtCreate:      e.GmtCreate,
		GmtModified:    e.GmtModified,
	}
}

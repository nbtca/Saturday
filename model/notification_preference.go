package model

// NotificationType represents the type of notification
type NotificationType string

const (
	NotifNewEventCreated    NotificationType = "new_event_created"
	NotifEventAssignedToMe  NotificationType = "event_assigned_to_me"
	NotifEventStatusChanged NotificationType = "event_status_changed"
)

// NotificationPreferences represents all notification preferences for a member
type NotificationPreferences struct {
	NewEventCreated    bool `json:"new_event_created"`
	EventAssignedToMe  bool `json:"event_assigned_to_me"`
	EventStatusChanged bool `json:"event_status_changed"`
}

// DefaultNotificationPreferences returns the default notification preferences
func DefaultNotificationPreferences() NotificationPreferences {
	return NotificationPreferences{
		NewEventCreated:    false, // Opt-in to avoid spam
		EventAssignedToMe:  true,  // Important notifications
		EventStatusChanged: true,  // Important notifications
	}
}

// NotificationPreferenceItem represents a single notification preference item for API responses
type NotificationPreferenceItem struct {
	NotificationType NotificationType `json:"notificationType"`
	Enabled          bool             `json:"enabled"`
	Description      string           `json:"description"`
}

// GetDescription returns a user-friendly description for each notification type
func GetNotificationDescription(notifType NotificationType) string {
	descriptions := map[NotificationType]string{
		NotifNewEventCreated:    "通知所有新创建的维修工单",
		NotifEventAssignedToMe:  "通知分配给我的工单状态变化",
		NotifEventStatusChanged: "通知我的工单所有状态变化",
	}
	return descriptions[notifType]
}

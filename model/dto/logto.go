package dto

type PatchLogtoUserRequest struct {
	UserName     string                 `json:"username,omitempty"`
	PrimaryEmail string                 `json:"primaryEmail,omitempty"`
	PrimaryPhone string                 `json:"primaryPhone,omitempty"`
	Name         string                 `json:"name,omitempty"`
	Avatar       string                 `json:"avatar,omitempty"`
	CustomData   map[string]interface{} `json:"customData,omitempty"`
}

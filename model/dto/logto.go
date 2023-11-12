package dto

type PatchLogtoUserRequest struct {
	UserName     string                 `json:"username"`
	PrimaryEmail string                 `json:"primaryEmail"`
	PrimaryPhone string                 `json:"primaryPhone"`
	Name         string                 `json:"name"`
	Avatar       string                 `json:"avatar"`
	CustomData   map[string]interface{} `json:"customData"`
}

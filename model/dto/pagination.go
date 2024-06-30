package dto

type Page struct {
	Offset uint64 `json:"-" form:"offset" binding:"min=0"`
	Limit  uint64 `json:"-" form:"limit" binding:"min=0"`
}

type PageRequest struct {
	Offset uint64 `query:"offset" default:"0" example:"0" minimum:"0" doc:"Offset"`
	Limit  uint64 `query:"limit" default:"50" example:"50" minimum:"0" doc:"Limit"`
}

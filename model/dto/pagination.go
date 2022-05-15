package dto

type Page struct {
	Offset uint64 `json:"-" form:"offset" binding:"min=0"`
	Limit  uint64 `json:"-" form:"limit" binding:"min=0"`
}

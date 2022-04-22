package router

import (
	"gin-example/src/service"
	"gin-example/util"
	"log"

	"github.com/gin-gonic/gin"
)

type MemberRouter struct {
}

type CreateMemberTokenReq struct {
	MemberId string `json:"member_id" validate:"required,len=10,numeric"`
	Password string `json:"password" validate:"required"`
}

type Page struct {
	Offset uint64 `json:"-" validate:"min=0"`
	Limit  uint64 `json:"-" validate:"min=0"`
}

func (MemberRouter) GetMemberById(c *gin.Context) {
	member, err := service.MemberServiceApp.GetMemberById(c.Param("MemberId"))
	serviceError, ok := util.IsServiceError(err)
	if ok {
		c.AbortWithStatusJSON(serviceError.Build())
		return
	}
	c.JSON(200, member)
}

func (MemberRouter) GetByPage(c *gin.Context) {
	offset, limit, err := util.GetPaginationQuery(c)
	if err != nil {
		log.Println(err)
	}
	members := service.MemberServiceApp.GetMembers(offset, limit)
	c.JSON(200, members)
}

func (MemberRouter) CreateToken(c *gin.Context) {
	c.JSON(200, "not implemented")
}

func (MemberRouter) Create(c *gin.Context) {
	c.JSON(200, "not implemented")
}

var MemberRouterApp = new(MemberRouter)

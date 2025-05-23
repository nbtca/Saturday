package router

import (
	"context"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-contrib/cors"
	"github.com/nbtca/saturday/middleware"
	"github.com/nbtca/saturday/service"
	"github.com/nbtca/saturday/util"

	"github.com/gin-gonic/gin"
)

type PingResponse struct {
	Pong string `json:"message" example:"ping" doc:"Ping message"`
}

func SetupRouter() *gin.Engine {
	Router := gin.Default()

	Router.Use(middleware.ErrorHandler)
	Router.Use(middleware.Logger)
	Router.Use(gin.Recovery())
	Router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://repair.nbtca.space", "https://nbtca.space", "http://localhost:5173"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			match, _ := regexp.MatchString(`https:\/\/.*\.nbtca\.space`, origin)
			return match
		},
		MaxAge: 12 * time.Hour,
	}))

	hook, _ := service.MakeGithubWebHook(os.Getenv("GITHUB_WEBHOOK_SECRET"))
	Router.Handle("POST", "/webhook", func(ctx *gin.Context) {
		err := hook.Handle(ctx.Request)
		if err != nil {
			util.Logger.Errorf("Error handling github webhook: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
	})

	logtoHook := service.LogtoWebHook{}
	Router.Handle("POST", "/webhook/logto", func(ctx *gin.Context) {
		err := logtoHook.Handle(ctx.Request)
		if err != nil {
			util.Logger.Errorf("Error handling logto webhook: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
	})

	api := humagin.New(Router, huma.DefaultConfig("Saturday API", "1.0.0"))

	huma.Register(api, huma.Operation{
		OperationID: "ping",
		Method:      http.MethodGet,
		Path:        "/ping",
		Summary:     "Ping",
		Tags:        []string{"Common", "Public"},
	}, func(ctx context.Context, input *struct{}) (*util.CommonResponse[PingResponse], error) {
		resp := PingResponse{
			Pong: "Hello",
		}
		return util.MakeCommonResponse(resp), nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-public-member",
		Method:      http.MethodGet,
		Path:        "/members/{MemberId}",
		Summary:     "Get a public member by id",
		Tags:        []string{"Member", "Public"},
	}, MemberRouterApp.GetPublicMemberById)

	huma.Register(api, huma.Operation{
		OperationID: "get-public-member-by-page",
		Method:      http.MethodGet,
		Path:        "/members",
		Summary:     "Get a public member by page",
		Tags:        []string{"Member", "Public"},
	}, MemberRouterApp.GetPublicMemberByPage)

	huma.Register(api, huma.Operation{
		OperationID: "create-token",
		Method:      http.MethodPost,
		Path:        "/members/{MemberId}/token",
		Summary:     "Create token",
		Tags:        []string{"Member", "Public"},
	}, MemberRouterApp.CreateToken)

	huma.Register(api, huma.Operation{
		OperationID: "create-token-via-logto-token",
		Method:      http.MethodGet,
		Path:        "/member/token/logto",
		Summary:     "Create token via logto token",
		Tags:        []string{"Member", "Public"},
	}, MemberRouterApp.CreateTokenViaLogtoToken)

	huma.Register(api, huma.Operation{
		OperationID: "bind-member-logto-id",
		Method:      http.MethodPatch,
		Path:        "/members/{MemberId}/logto_id",
		Summary:     "Bind member logto id",
		Tags:        []string{"Member", "Public"},
	}, MemberRouterApp.BindMemberLogtoId)

	huma.Register(api, huma.Operation{
		OperationID: "create-token-via-wechat",
		Method:      http.MethodPost,
		Path:        "/clients/token/wechat",
		Summary:     "Create token via wechat",
		Tags:        []string{"Client", "Public"},
	}, ClientRouterApp.CreateTokenViaWeChat)

	Router.POST("/clients/token/logto", middleware.Auth("client"), ClientRouterApp.CreateTokenViaLogto)

	huma.Register(api, huma.Operation{
		OperationID: "get-public-event-by-id",
		Method:      http.MethodGet,
		Path:        "/events/{EventId}",
		Summary:     "Get a public event by id",
		Tags:        []string{"Event", "Public"},
	}, EventRouterApp.GetPublicEventById)

	huma.Register(api, huma.Operation{
		OperationID: "get-public-event-by-page",
		Method:      http.MethodGet,
		Path:        "/events",
		Summary:     "Get a public event by page",
		Tags:        []string{"Event", "Public"},
	}, EventRouterApp.GetPublicEventByPage)

	huma.Register(api, huma.Operation{
		OperationID: "create-member-with-logto",
		Method:      http.MethodPost,
		Path:        "/members/:MemberId/logto",
		Summary:     "Create member with logto",
		Tags:        []string{"Member", "Private"},
	}, MemberRouterApp.CreateWithLogto)

	Router.PATCH("member/activate",
		middleware.Auth("member_inactive", "admin_inactive"),
		MemberRouterApp.Activate)

	MemberGroup := Router.Group("/")
	MemberGroup.Use(middleware.Auth("member", "admin"))
	{
		MemberGroup.GET("/member", MemberRouterApp.GetMemberById)
		MemberGroup.PUT("/member", MemberRouterApp.Update)
		MemberGroup.PATCH("/member/avatar", MemberRouterApp.UpdateAvatar)

		MemberGroup.GET("member/events", EventRouterApp.GetMemberEventByPage)
		MemberGroup.GET("member/events/:EventId", EventRouterApp.GetEventById)
		/*
			!!! IMPORTANT !!!
			this middleware is REQUIRED before all handlers that uses event action (except create)
			or there will be panic
		*/
		MemberGroup.Use(middleware.EventActionPreProcess)
		MemberGroup.POST("member/events/:EventId/accept", EventRouterApp.Accept)
		MemberGroup.DELETE("member/events/:EventId/accept", EventRouterApp.Drop)
		MemberGroup.POST("member/events/:EventId/commit", EventRouterApp.Commit)
		MemberGroup.PATCH("member/events/:EventId/commit", EventRouterApp.AlterCommit)

		// MemberGroup.GET("client/:ClientId/events", EventRouterApp.GetEventByClientAndPage)

	}

	AdminGroup := Router.Group("/")
	AdminGroup.Use(middleware.Auth("admin"))
	{
		AdminGroup.POST("/members", MemberRouterApp.CreateMany)
		AdminGroup.POST("/members/:MemberId", MemberRouterApp.Create)
		// TODO change this path
		AdminGroup.GET("/members/full", MemberRouterApp.GetMemberByPage)
		// AdminGroup.PATCH("/members/:MemberId", MemberRouterApp.UpdateBasic)
		AdminGroup.GET("/events/xlsx", EventRouterApp.ExportEventsToXlsx)

		AdminGroup.Use(middleware.EventActionPreProcess)
		AdminGroup.DELETE("/events/:EventId/commit", EventRouterApp.RejectCommit)
		AdminGroup.POST("/events/:EventId/close", EventRouterApp.Close)
	}

	ClientGroup := Router.Group("/")
	ClientGroup.Use(middleware.Auth("client"))
	{
		ClientGroup.GET("/client/events/:EventId", EventRouterApp.GetEventById)
		ClientGroup.GET("/client/events", EventRouterApp.GetClientEventByPage)
		ClientGroup.POST("/client/event", EventRouterApp.Create)

		ClientGroup.Use(middleware.EventActionPreProcess)
		ClientGroup.PATCH("/client/events/:EventId", EventRouterApp.Update)
		ClientGroup.DELETE("/client/events/:EventId", EventRouterApp.Cancel)
	}

	Router.POST("/upload", middleware.Auth("member", "admin", "client"), CommonRouterApp.Upload)

	return Router
}

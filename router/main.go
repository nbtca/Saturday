package router

import (
	"context"
	"net/http"
	"regexp"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/nbtca/saturday/middleware"
	"github.com/nbtca/saturday/service"
	"github.com/nbtca/saturday/util"
	"github.com/spf13/viper"
)

type PingResponse struct {
	Pong string `json:"message" example:"ping" doc:"Ping message"`
}

func SetupRouter() *chi.Mux {
	// Create Chi router
	router := chi.NewRouter()

	// Add CORS middleware
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://repair.nbtca.space", "https://nbtca.space", "http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(r *http.Request, origin string) bool {
			match, _ := regexp.MatchString(`https:\/\/.*\.nbtca\.space`, origin)
			return match
		},
		MaxAge: int((12 * time.Hour).Seconds()),
	}))

	// Create Huma API
	config := huma.DefaultConfig("Saturday API", "1.0.0")
	config.Servers = []*huma.Server{
		{URL: "https://api.nbtca.space", Description: "Production server"},
		{URL: "http://localhost:4000", Description: "Development server"},
	}
	
	api := humachi.New(router, config)

	// Add Huma middleware
	api.UseMiddleware(middleware.HumaLogger())

	// Keep webhooks as raw endpoints since they don't need OpenAPI documentation
	hook, _ := service.MakeGithubWebHook(viper.GetString("github.webhook_secret"))
	router.Post("/webhook", func(w http.ResponseWriter, r *http.Request) {
		err := hook.Handle(r)
		if err != nil {
			util.Logger.Errorf("Error handling github webhook: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		}
	})

	logtoHook := service.LogtoWebHook{}
	router.Post("/webhook/logto", func(w http.ResponseWriter, r *http.Request) {
		err := logtoHook.Handle(r)
		if err != nil {
			util.Logger.Errorf("Error handling logto webhook: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		}
	})

	// Public endpoints (no authentication required)
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

	// Client authenticated endpoints
	huma.Register(api, huma.Operation{
		OperationID: "create-token-via-logto",
		Method:      http.MethodPost,
		Path:        "/clients/token/logto",
		Summary:     "Create token via logto",
		Tags:        []string{"Client", "Private"},
	}, ClientRouterApp.CreateTokenViaLogto)

	huma.Register(api, huma.Operation{
		OperationID: "get-client-event-by-id",
		Method:      http.MethodGet,
		Path:        "/client/events/{EventId}",
		Summary:     "Get client event by id",
		Tags:        []string{"Event", "Client"},
	}, EventRouterApp.GetEventById)

	huma.Register(api, huma.Operation{
		OperationID: "get-client-events",
		Method:      http.MethodGet,
		Path:        "/client/events",
		Summary:     "Get client events by page",
		Tags:        []string{"Event", "Client"},
	}, EventRouterApp.GetClientEventByPage)

	huma.Register(api, huma.Operation{
		OperationID: "create-client-event",
		Method:      http.MethodPost,
		Path:        "/client/event",
		Summary:     "Create client event",
		Tags:        []string{"Event", "Client"},
	}, EventRouterApp.Create)

	huma.Register(api, huma.Operation{
		OperationID: "update-client-event",
		Method:      http.MethodPatch,
		Path:        "/client/events/{EventId}",
		Summary:     "Update client event",
		Tags:        []string{"Event", "Client"},
	}, EventRouterApp.Update)

	huma.Register(api, huma.Operation{
		OperationID: "cancel-client-event",
		Method:      http.MethodDelete,
		Path:        "/client/events/{EventId}",
		Summary:     "Cancel client event",
		Tags:        []string{"Event", "Client"},
	}, EventRouterApp.Cancel)

	// Member authenticated endpoints
	huma.Register(api, huma.Operation{
		OperationID: "activate-member",
		Method:      http.MethodPatch,
		Path:        "/member/activate",
		Summary:     "Activate member",
		Tags:        []string{"Member", "Private"},
	}, MemberRouterApp.Activate)

	huma.Register(api, huma.Operation{
		OperationID: "get-member",
		Method:      http.MethodGet,
		Path:        "/member",
		Summary:     "Get current member",
		Tags:        []string{"Member", "Private"},
	}, MemberRouterApp.GetMemberById)

	huma.Register(api, huma.Operation{
		OperationID: "update-member",
		Method:      http.MethodPut,
		Path:        "/member",
		Summary:     "Update member",
		Tags:        []string{"Member", "Private"},
	}, MemberRouterApp.Update)

	huma.Register(api, huma.Operation{
		OperationID: "update-member-avatar",
		Method:      http.MethodPatch,
		Path:        "/member/avatar",
		Summary:     "Update member avatar",
		Tags:        []string{"Member", "Private"},
	}, MemberRouterApp.UpdateAvatar)

	huma.Register(api, huma.Operation{
		OperationID: "get-member-events",
		Method:      http.MethodGet,
		Path:        "/member/events",
		Summary:     "Get member events",
		Tags:        []string{"Event", "Member"},
	}, EventRouterApp.GetMemberEventByPage)

	huma.Register(api, huma.Operation{
		OperationID: "get-member-event-by-id",
		Method:      http.MethodGet,
		Path:        "/member/events/{EventId}",
		Summary:     "Get member event by id",
		Tags:        []string{"Event", "Member"},
	}, EventRouterApp.GetEventById)

	huma.Register(api, huma.Operation{
		OperationID: "accept-event",
		Method:      http.MethodPost,
		Path:        "/member/events/{EventId}/accept",
		Summary:     "Accept event",
		Tags:        []string{"Event", "Member"},
	}, EventRouterApp.Accept)

	huma.Register(api, huma.Operation{
		OperationID: "drop-event",
		Method:      http.MethodDelete,
		Path:        "/member/events/{EventId}/accept",
		Summary:     "Drop event",
		Tags:        []string{"Event", "Member"},
	}, EventRouterApp.Drop)

	huma.Register(api, huma.Operation{
		OperationID: "commit-event",
		Method:      http.MethodPost,
		Path:        "/member/events/{EventId}/commit",
		Summary:     "Commit event",
		Tags:        []string{"Event", "Member"},
	}, EventRouterApp.Commit)

	huma.Register(api, huma.Operation{
		OperationID: "alter-commit-event",
		Method:      http.MethodPatch,
		Path:        "/member/events/{EventId}/commit",
		Summary:     "Alter commit event",
		Tags:        []string{"Event", "Member"},
	}, EventRouterApp.AlterCommit)

	// Admin authenticated endpoints
	huma.Register(api, huma.Operation{
		OperationID: "create-members",
		Method:      http.MethodPost,
		Path:        "/members",
		Summary:     "Create multiple members",
		Tags:        []string{"Member", "Admin"},
	}, MemberRouterApp.CreateMany)

	huma.Register(api, huma.Operation{
		OperationID: "create-member",
		Method:      http.MethodPost,
		Path:        "/members/{MemberId}",
		Summary:     "Create member",
		Tags:        []string{"Member", "Admin"},
	}, MemberRouterApp.Create)

	huma.Register(api, huma.Operation{
		OperationID: "get-members-full",
		Method:      http.MethodGet,
		Path:        "/members/full",
		Summary:     "Get members with full details",
		Tags:        []string{"Member", "Admin"},
	}, MemberRouterApp.GetMemberByPage)

	huma.Register(api, huma.Operation{
		OperationID: "update-member-basic",
		Method:      http.MethodPatch,
		Path:        "/members/{MemberId}",
		Summary:     "Update member basic info",
		Tags:        []string{"Member", "Admin"},
	}, MemberRouterApp.UpdateBasic)

	huma.Register(api, huma.Operation{
		OperationID: "export-events-xlsx",
		Method:      http.MethodGet,
		Path:        "/events/xlsx",
		Summary:     "Export events to XLSX",
		Tags:        []string{"Event", "Admin"},
	}, EventRouterApp.ExportEventsToXlsx)

	huma.Register(api, huma.Operation{
		OperationID: "reject-commit-event",
		Method:      http.MethodDelete,
		Path:        "/events/{EventId}/commit",
		Summary:     "Reject commit event",
		Tags:        []string{"Event", "Admin"},
	}, EventRouterApp.RejectCommit)

	huma.Register(api, huma.Operation{
		OperationID: "close-event",
		Method:      http.MethodPost,
		Path:        "/events/{EventId}/close",
		Summary:     "Close event",
		Tags:        []string{"Event", "Admin"},
	}, EventRouterApp.Close)

	// TODO: Upload endpoint - needs special multipart handling
	// For now, keep as commented until Huma multipart is implemented
	/*
	huma.Register(api, huma.Operation{
		OperationID: "upload-file",
		Method:      http.MethodPost,
		Path:        "/upload",
		Summary:     "Upload file",
		Tags:        []string{"Common", "Private"},
		Middlewares: huma.Middlewares{middleware.HumaAuth("member", "admin", "client")},
	}, CommonRouterApp.Upload)
	*/

	return router
}
package router

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/nbtca/saturday/middleware"
	"github.com/nbtca/saturday/model"
	"github.com/nbtca/saturday/model/dto"
	"github.com/nbtca/saturday/repo"
	"github.com/nbtca/saturday/service"
	"github.com/nbtca/saturday/util"
)

type EventRouter struct {
	huma huma.API
}

func (EventRouter) GetPublicEventById(c context.Context, input *struct {
	EventID int64 `path:"EventId"`
}) (*util.CommonResponse[model.PublicEvent], error) {
	event, err := service.EventServiceApp.GetPublicEventById(input.EventID)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	return util.MakeCommonResponse(event), nil
}

func (EventRouter) GetPublicEventByPage(c context.Context, input *struct {
	dto.PageRequest
	Status []string `query:"status"`
	Order  string   `query:"order" default:"ASC"`
}) (*util.CommonResponse[[]model.PublicEvent], error) {
	filter := repo.EventFilter{
		Offset: input.Offset,
		Limit:  input.Limit,
		Status: input.Status,
		Order:  input.Order,
	}
	events, count, err := service.EventServiceApp.GetPublicEventsWithCount(filter)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	return util.MakePaginatedResponse(events, count, input.Offset, input.Limit), nil
}

func (er EventRouter) GetEventById(ctx context.Context, input *GetEventByIdInput) (*util.CommonResponse[model.Event], error) {
	auth, err := middleware.AuthenticateUser(input.Authorization, "member", "admin", "client")
	if err != nil {
		return nil, err
	}

	event, err := service.EventServiceApp.GetEventById(input.EventId)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	clientId, err := middleware.GetClientIdFromAuth(auth)
	if err != nil {
		return nil, err
	}

	if event.MemberId != auth.ID && event.ClientId != clientId {
		return nil, huma.Error401Unauthorized("not authorized")
	}

	return util.MakeCommonResponse(event), nil
}

// return events that is accepted by current member
func (EventRouter) GetMemberEventByPage(ctx context.Context, input *GetMemberEventByPageInput) (*util.CommonResponse[[]model.Event], error) {
	auth, err := middleware.AuthenticateUser(input.Authorization, "member", "admin")
	if err != nil {
		return nil, err
	}
	var status []string
	if input.Status != "" {
		status = []string{input.Status}
	}
	filter := repo.EventFilter{
		Offset: input.Offset,
		Limit:  input.Limit,
		Status: status,
		Order:  input.Order,
	}
	events, count, err := service.EventServiceApp.GetMemberEventsWithCount(filter, auth.ID)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	return util.MakePaginatedResponse(events, count, input.Offset, input.Limit), nil
}

// ExportEventsToXlsx exports events to XLSX format
// Note: This endpoint returns raw file data instead of CommonResponse
func (EventRouter) ExportEventsToXlsx(ctx context.Context, input *ExportEventsToXlsxInput) (*huma.StreamResponse, error) {
	_, err := middleware.AuthenticateUser(input.Authorization, "admin")
	if err != nil {
		return nil, err
	}

	var status []string
	if input.Status != "" {
		status = []string{input.Status}
	}
	f, err := service.EventServiceApp.ExportEventToXlsx(repo.EventFilter{
		Offset: 0,
		Limit:  1000000,
		Status: status,
		Order:  input.Order,
	}, input.StartTime, input.EndTime)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	formatDate := func(date string) string {
		t, err := time.Parse("2006-01-02T15:04:05Z", date)
		if err != nil {
			return date
		}
		return t.Format("2006-01-02")
	}

	filename := fmt.Sprintf("events_%s_to_%s.xlsx", formatDate(input.StartTime), formatDate(input.EndTime))

	// Create a buffer to write the Excel file
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, huma.Error500InternalServerError("Failed to generate Excel file")
	}

	return &huma.StreamResponse{
		Body: func(ctx huma.Context) {
			w := ctx.BodyWriter()
			ctx.SetHeader("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
			ctx.SetHeader("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
			ctx.SetHeader("Content-Transfer-Encoding", "binary")
			w.Write(buf.Bytes())
		},
	}, nil
}

func (EventRouter) Accept(ctx context.Context, input *AcceptEventInput) (*util.CommonResponse[model.Event], error) {
	auth, err := middleware.AuthenticateUser(input.Authorization, "member", "admin")
	if err != nil {
		return nil, err
	}

	event, err := middleware.LoadEvent(input.EventId)
	if err != nil {
		return nil, err
	}

	identity := middleware.CreateIdentityFromAuth(auth)
	if err := service.EventServiceApp.Accept(&event, identity); err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	return util.MakeCommonResponse(event), nil
}

func (EventRouter) Drop(ctx context.Context, input *DropEventInput) (*util.CommonResponse[model.Event], error) {
	auth, err := middleware.AuthenticateUser(input.Authorization, "member", "admin")
	if err != nil {
		return nil, err
	}

	event, err := middleware.LoadEvent(input.EventId)
	if err != nil {
		return nil, err
	}

	identity := middleware.CreateIdentityFromAuth(auth)
	if err := service.EventServiceApp.Act(&event, identity, util.Drop); err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	return util.MakeCommonResponse(event), nil
}

func (EventRouter) Commit(ctx context.Context, input *CommitEventInput) (*util.CommonResponse[model.Event], error) {
	auth, err := middleware.AuthenticateUser(input.Authorization, "member", "admin")
	if err != nil {
		return nil, err
	}

	event, err := middleware.LoadEvent(input.EventId)
	if err != nil {
		return nil, err
	}

	identity := middleware.CreateIdentityFromAuth(auth)
	if input.Body.Size != "" {
		event.Size = input.Body.Size
	}
	opts := service.ActOptions{
		Description: input.Body.Content,
		Images:      input.Body.Images,
	}
	if err := service.EventServiceApp.ActWithOptions(&event, identity, util.Commit, opts); err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	return util.MakeCommonResponse(event), nil
}

func (EventRouter) AlterCommit(ctx context.Context, input *AlterCommitEventInput) (*util.CommonResponse[model.Event], error) {
	auth, err := middleware.AuthenticateUser(input.Authorization, "member", "admin")
	if err != nil {
		return nil, err
	}

	event, err := middleware.LoadEvent(input.EventId)
	if err != nil {
		return nil, err
	}

	identity := middleware.CreateIdentityFromAuth(auth)
	if input.Body.Size != "" {
		event.Size = input.Body.Size
	}
	opts := service.ActOptions{
		Description: input.Body.Content,
		Images:      input.Body.Images,
	}
	if err := service.EventServiceApp.ActWithOptions(&event, identity, util.AlterCommit, opts); err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	return util.MakeCommonResponse(event), nil
}

func (EventRouter) RejectCommit(ctx context.Context, input *RejectCommitEventInput) (*util.CommonResponse[model.Event], error) {
	auth, err := middleware.AuthenticateUser(input.Authorization, "admin")
	if err != nil {
		return nil, err
	}

	event, err := middleware.LoadEvent(input.EventId)
	if err != nil {
		return nil, err
	}

	identity := middleware.CreateIdentityFromAuth(auth)
	if err := service.EventServiceApp.Act(&event, identity, util.Reject); err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	return util.MakeCommonResponse(event), nil
}

func (EventRouter) Close(ctx context.Context, input *CloseEventInput) (*util.CommonResponse[model.Event], error) {
	auth, err := middleware.AuthenticateUser(input.Authorization, "admin")
	if err != nil {
		return nil, err
	}

	event, err := middleware.LoadEvent(input.EventId)
	if err != nil {
		return nil, err
	}

	identity := middleware.CreateIdentityFromAuth(auth)
	if err := service.EventServiceApp.Act(&event, identity, util.Close); err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	return util.MakeCommonResponse(event), nil
}

func (er EventRouter) GetClientEventByPage(ctx context.Context, input *GetClientEventByPageInput) (*util.CommonResponse[[]model.Event], error) {
	auth, err := middleware.AuthenticateUser(input.Authorization, "client")
	if err != nil {
		return nil, err
	}

	clientId, err := middleware.GetClientIdFromAuth(auth)
	if err != nil {
		return nil, err
	}
	var status []string
	if input.Status != "" {
		status = []string{input.Status}
	}
	filter := repo.EventFilter{
		Offset: input.Offset,
		Limit:  input.Limit,
		Status: status,
		Order:  input.Order,
	}
	events, count, err := service.EventServiceApp.GetClientEventsWithCount(filter, clientId)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	return util.MakePaginatedResponse(events, count, input.Offset, input.Limit), nil
}

func (EventRouter) Create(ctx context.Context, input *CreateClientEventInput) (*util.CommonResponse[model.Event], error) {
	auth, err := middleware.AuthenticateUser(input.Authorization, "client")
	if err != nil {
		return nil, err
	}

	clientId, err := middleware.GetClientIdFromAuth(auth)
	if err != nil {
		return nil, err
	}

	event := &model.Event{
		ClientId:          clientId,
		Model:             input.Body.Model,
		Phone:             input.Body.Phone,
		QQ:                input.Body.QQ,
		ContactPreference: input.Body.ContactPreference,
		Problem:           input.Body.Problem,
		Images:            model.StringSlice(input.Body.Images),
	}
	err = service.EventServiceApp.CreateEvent(event)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	return util.MakeCommonResponse(*event), nil
}

func (er EventRouter) Update(ctx context.Context, input *UpdateClientEventInput) (*util.CommonResponse[model.Event], error) {
	auth, err := middleware.AuthenticateUser(input.Authorization, "client")
	if err != nil {
		return nil, err
	}

	event, err := middleware.LoadEvent(input.EventId)
	if err != nil {
		return nil, err
	}

	clientId, err := middleware.GetClientIdFromAuth(auth)
	if err != nil {
		return nil, err
	}

	identity := middleware.CreateIdentityFromAuth(auth)
	identity.ClientId = clientId

	if input.Body.Phone != "" {
		event.Phone = input.Body.Phone
	}
	if input.Body.QQ != "" {
		event.QQ = input.Body.QQ
	}
	if input.Body.Problem != "" {
		event.Problem = input.Body.Problem
	}
	if input.Body.Model != "" {
		event.Model = input.Body.Model
	}
	if input.Body.ContactPreference != "" {
		event.ContactPreference = input.Body.ContactPreference
	}
	if input.Body.Size != "" {
		event.Size = input.Body.Size
	}
	if input.Body.Images != nil {
		event.Images = model.StringSlice(input.Body.Images)
	}
	if err := service.EventServiceApp.Act(&event, identity, util.Update); err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	return util.MakeCommonResponse(event), nil
}

func (er EventRouter) Cancel(ctx context.Context, input *CancelClientEventInput) (*util.CommonResponse[model.Event], error) {
	auth, err := middleware.AuthenticateUser(input.Authorization, "client")
	if err != nil {
		return nil, err
	}

	event, err := middleware.LoadEvent(input.EventId)
	if err != nil {
		return nil, err
	}

	clientId, err := middleware.GetClientIdFromAuth(auth)
	if err != nil {
		return nil, err
	}

	identity := middleware.CreateIdentityFromAuth(auth)
	identity.ClientId = clientId

	if err := service.EventServiceApp.Act(&event, identity, util.Cancel); err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	return util.MakeCommonResponse(event), nil
}

var EventRouterApp = EventRouter{}

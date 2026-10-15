package router

import (
	"testing"

	"github.com/nbtca/saturday/middleware"
	"github.com/nbtca/saturday/model"
)

func TestCanViewEvent(t *testing.T) {
	assigned := model.Event{MemberId: "2333333333", ClientId: 7}
	open := model.Event{MemberId: "", ClientId: 7}
	tests := []struct {
		name     string
		event    model.Event
		auth     middleware.AuthContext
		clientId int64
		want     bool
	}{
		{"assigned member", assigned, middleware.AuthContext{ID: "2333333333", Role: "member"}, 1, true},
		{"other member", assigned, middleware.AuthContext{ID: "1111111111", Role: "member"}, 1, false},
		{"admin", assigned, middleware.AuthContext{ID: "1111111111", Role: "admin"}, 1, true},
		{"owning client", assigned, middleware.AuthContext{}, 7, true},
		{"other client on open event", open, middleware.AuthContext{}, 8, false},
		{"other client on assigned event", assigned, middleware.AuthContext{}, 8, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := canViewEvent(tt.event, &tt.auth, tt.clientId); got != tt.want {
				t.Errorf("canViewEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

package service_test

import (
	"testing"

	"github.com/nbtca/saturday/service"
)

func TestMapLogtoUserRole(t *testing.T) {
	tests := []struct {
		name  string
		roles []service.LogtoUserRole
		want  string
	}{
		{
			name: "Single Repair Admin",
			roles: []service.LogtoUserRole{
				{Name: "Repair Admin"},
			},
			want: "admin",
		},
		{
			name: "Single Repair Member",
			roles: []service.LogtoUserRole{
				{Name: "Repair Member"},
			},
			want: "member",
		},
		{
			name: "Repair Admin and Repair Member",
			roles: []service.LogtoUserRole{
				{Name: "Repair Admin"},
				{Name: "Repair Member"},
			},
			want: "admin",
		},
		{
			name:  "No Roles",
			roles: []service.LogtoUserRole{},
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := service.MemberService{}
			if got := service.MapLogtoUserRole(tt.roles); got != tt.want {
				t.Errorf("MapLogtoUserRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

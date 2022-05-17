package router_test

import (
	"testing"
)

func TestGetPublicMemberById(t *testing.T) {
	for _, data := range GetPublicMemberData {
		t.Run(data.Name, func(t *testing.T) {
			err := DataHandler(data)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestGetMemberById(t *testing.T) {
	for _, data := range GetMemberData {
		t.Run(data.Name, func(t *testing.T) {
			err := DataHandler(data)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestMemberActive(t *testing.T) {
	for _, data := range MemberActiveData {
		t.Run(data.Name, func(t *testing.T) {
			err := DataHandler(data)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestUpdateMember(t *testing.T) {
	for _, data := range UpdateMemberData {
		t.Run(data.Name, func(t *testing.T) {
			err := DataHandler(data)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestUpdateMemberBasic(t *testing.T) {
	for _, data := range UpdateMemberBasicInfoData {
		t.Run(data.Name, func(t *testing.T) {
			err := DataHandler(data)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
func TestCreateMember(t *testing.T) {
	for _, data := range CreateMemberData {
		t.Run(data.Name, func(t *testing.T) {
			err := DataHandler(data)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

package util_test

import (
	"saturday/util"
	"testing"
)

func TestCreateToken(t *testing.T) {
	j, _ := util.CreateToken(util.Payload{Who: "123", Role: "member"})
	_, claims, _ := util.ParseToken(j)
	if claims.Who != "123" || claims.Role != "member" {
		t.Error("测试失败")
	}
}

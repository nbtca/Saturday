package util_test

import (
	"saturday/src/util"
	"testing"
)

func TestCreateToken(t *testing.T) {
	j, _ := util.CreateToken(util.Payload{Id: "123", Role: "member"})
	_, claims, _ := util.ParseToken(j)
	if claims.Id != "123" || claims.Role != "member" {
		t.Error("测试失败")
	}
}

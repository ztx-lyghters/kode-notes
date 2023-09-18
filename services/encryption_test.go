package services

import (
	"testing"

	"github.com/ztx-lyghters/kode-notes/core"
	"github.com/ztx-lyghters/kode-notes/repository"
)

const (
	TEST_USER_ID     = 666
	TEST_USERNAME    = "Vasyan1337"
	TEST_USER_PASSWD = "qwertyasdf"
)

func TestValidateToken(t *testing.T) {
	test_user := &core.User{
		Id:       TEST_USER_ID,
		Username: TEST_USERNAME,
		Password: TEST_USER_PASSWD,
	}
	test_token := assembleTokenJWT(test_user)

	mock_auth := NewAuthService(&repository.Repository{})
	id_from_token, err := mock_auth.ValidateToken(test_token)
	if err != nil {
		t.Error("ValidateToken failed: " + err.Error())
	}
	if id_from_token != test_user.Id {
		t.Error("User ID from token doesn't match user ID" +
			"from 'test_user' struct")
	}
}

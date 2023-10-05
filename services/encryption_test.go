package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"testing"
	"time"

	"github.com/ztx-lyghters/kode-notes/core"
)

const (
	TEST_USER_ID     = 666
	TEST_USERNAME    = "Vasyan1337"
	TEST_USER_PASSWD = "qwertyasdf"
)

func mockAssembleTokenJWT(h string, p string, s string, encoder func(string) string) string {
	if encoder == nil {
		encoder = func(s string) string {
			return encodeSegmentJWT([]byte(s))
		}
	}
	header := encoder(h)
	payload := encoder(p)

	hmac := hmac.New(sha256.New, []byte(s))
	hmac.Write([]byte(header + "." + payload))

	signature := encodeSegmentJWT(hmac.Sum(nil))

	token := fmt.Sprintf("%s.%s.%s", header, payload,
		signature)

	return token
}

func TestValidateToken(t *testing.T) {
	test_user := &core.User{
		Id:       TEST_USER_ID,
		Username: TEST_USERNAME,
		Password: TEST_USER_PASSWD,
	}
	mock_auth := &Auth{}

	test_token := assembleTokenJWT(test_user)
	id_from_token, err := mock_auth.ValidateToken(test_token)

	if err != nil {
		t.Error("ValidateToken failed: " + err.Error())
	}
	if id_from_token != test_user.Id {
		t.Error("User ID from token doesn't match user ID" +
			"from 'test_user' struct")
	}

	var cases = []struct {
		id    uint
		token string
		error string
	}{
		{
			token: "iamtokentrustme",
			error: "invalid token format",
		},
		{
			token: mockAssembleTokenJWT(
				`{"alg":"HS256", "typ":"abc"}`,
				fmt.Sprintf(
					`{"user_id":%d,"iat":%d}`,
					1, time.Now().Unix()),
				"trustmeiamasignature",
				nil,
			),
			error: "invalid header: expecting JWT",
		},
		{
			token: mockAssembleTokenJWT(
				`{"alg":"plain", "typ":"JWT"}`,
				fmt.Sprintf(
					`{"user_id":%d,"iat":%d}`,
					1, time.Now().Unix()),
				"trustmeiamasignature",
				nil,
			),
			error: "invalid header: expecting HS256",
		},
	}

	for i, c := range cases {
		id, err := mock_auth.ValidateToken(c.token)
		if err == nil {
			t.Fatalf("ValidateToken failed: expected \"%s\","+
				" got nil", c.error)
		}
		if err.Error() != c.error {
			t.Fatalf("ValidateToken failed: expected the \"%s\" error, but got "+err.Error(), c.error)
		}

		if c.id != 0 {
			if id != c.id {
				t.Fatalf("validateToken failed: user id "+
					" mismatch with following case: %v", cases[i])
			}
		}
	}
}

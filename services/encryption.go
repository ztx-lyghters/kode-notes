package services

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ztx-lyghters/kode-notes/core"
)

const (
	TOKEN_TTL  = 60 * time.Minute
	JWT_HEADER = `{"alg":"HS256","typ":"JWT"}`

	// Quick and dirty security theater.
	// At least the JWT signing key shoud probably
	// be generated at runtime, so it is different
	// every time. Might intoduce "inconvenience"
	// if the server gets down, but I'd take security
	// over convenience. For now it's hardcoded.
	HASH_SALT = "TheSaltiestSaltOfSaltiestSaltsOutThere"
	SIGN_KEY  = "ASuperSecretPrivateSigningKey"
)

func (s *Auth) ValidateToken(token string) (uint, error) {
	var header_map, claims map[string]interface{}

	segments := strings.Split(token, ".")
	if len(segments) != 3 {
		return 0, errors.New("invalid token format")
	}

	header, err := base64.RawURLEncoding.
		DecodeString(segments[0])
	if err != nil {
		return 0, errors.New("invalid header format")
	}
	if !bytes.Contains(header, []byte(`"typ":"JWT"`)) {
		return 0, errors.New("invalid header: expecting JWT")
	}
	if !bytes.Contains(header, []byte(`"alg":"HS256"`)) {
		return 0, errors.New("invalid header: expecting HS256")
	}

	payload, err := base64.RawURLEncoding.
		DecodeString(segments[1])
	if err != nil {
		return 0, errors.New("invalid payload format")
	}

	signature, err := base64.RawURLEncoding.
		DecodeString(segments[2])
	if err != nil {
		return 0, errors.New("invalid signature format")
	}

	err = json.Unmarshal(header, &header_map)
	if err != nil {
		return 0, errors.New("invalid header format")
	}

	alg, ok := header_map["alg"].(string)
	if !ok || alg != "HS256" {
		return 0, errors.New("invalid encoding algorythm")
	}

	hmac := hmac.New(sha256.New, []byte(SIGN_KEY))
	hmac.Write([]byte(segments[0] + "." + segments[1]))
	expected := hmac.Sum(nil)

	invalid := errors.New("invalid signature")

	if len(signature) != len(expected) {
		return 0, invalid
	}
	for i := 0; i < len(signature); i++ {
		if signature[i] != expected[i] {
			return 0, invalid
		}
	}

	if err := json.Unmarshal(payload, &claims); err != nil {
		return 0, errors.New("invalid payload format")
	}

	iat, ok := claims["iat"].(float64)
	if !ok {
		return 0, errors.New("invalid token issue time")
	}
	t_iat := time.Unix(int64(iat), 0)
	t_exp := t_iat.Add(TOKEN_TTL).Unix()
	t_now := time.Now().Unix()
	if t_now >= t_exp {
		return 0, errors.New("token has expired")
	}

	user_id, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf(fmt.Sprintf("invalid user_id field - %f", user_id))
	}

	return uint(user_id), nil
}

// Probably should take a payload argument, in order
// to encode different things, but no time for that :)
func assembleTokenJWT(user *core.User) string {
	payloadJSON := fmt.Sprintf(
		`{"user_id":%d,"iat":%d}`,
		user.Id, time.Now().Unix())

	header := encodeSegmentJWT([]byte(JWT_HEADER))
	payload := encodeSegmentJWT([]byte(payloadJSON))

	hmac := hmac.New(sha256.New, []byte(SIGN_KEY))
	hmac.Write([]byte(header + "." + payload))

	signature := encodeSegmentJWT(hmac.Sum(nil))

	token := fmt.Sprintf("%s.%s.%s", header, payload,
		signature)

	return token
}

func encodeSegmentJWT(segment []byte) string {
	return strings.TrimRight(base64.URLEncoding.
		EncodeToString(segment), "=")
}

func GeneratePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(HASH_SALT)))
}

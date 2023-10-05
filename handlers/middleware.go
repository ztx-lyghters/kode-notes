package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	AUTH_HEADER = "Authorization"
	USER_CTX    = "user_id"
)

func (h *Handler) userIdentify(w http.ResponseWriter,
	r *http.Request) error {

	header := r.Header.Get(AUTH_HEADER)
	if strings.TrimSpace(header) == "" {
		return errors.New("empty auth header")
	}

	header_parts := strings.Split(header, " ")
	if len(header_parts) != 2 {
		return errors.New("invalid auth header")
	}

	user_id, err := h.services.AuthorizationService.
		ValidateToken(header_parts[1])
	if err != nil {
		return err
	}

	if !h.services.AuthorizationService.
		CheckUserExistence(user_id) {
		return errors.New("user does not exist")
	}

	r.Header.Set(USER_CTX, fmt.Sprint(user_id))
	return nil
}

func getUserId(w http.ResponseWriter, r *http.Request) (uint, error) {
	id := r.Header.Get(USER_CTX)
	if id == "" {
		err := errors.New("user id not found")
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return 0, err
	}

	id_int, err := strconv.ParseUint(id, 10, 0)
	if err != nil {
		err := errors.New("user id not found 2")
		NewErrorResponse(w, http.StatusInternalServerError, err.Error())
		return 0, err
	}

	return uint(id_int), nil
}

package handlers

import (
	"fmt"
	"net/http"

	"example.com/kode-notes/core"
	"github.com/sirupsen/logrus"
)

func (h *Handler) SignIn(w http.ResponseWriter,
	r *http.Request) {

	user := &core.User{}
	err := parseJSON(r, user)
	if err != nil {
		NewErrorResponse(w, http.StatusBadRequest,
			"bad credentials form")
		return
	}

	token, err := h.services.AuthorizationService.
		GenerateToken(user)
	if err != nil {
		logrus.Infof("Unsuccessful auth attempt by a user [ID: %d, username: %s]: %s",
			user.Id, user.Username, err.Error())
		NewErrorResponse(w, http.StatusUnauthorized,
			"bad credentials")
		return
	}

	newJSONResponse(w, http.StatusOK, map[string]interface{}{
		"token": token,
	})
	logrus.Infof("Issued a token to user '%s'", user.Username)
}

func (h *Handler) SignUp(w http.ResponseWriter,
	r *http.Request) {

	user := &core.User{}
	err := parseJSON(r, user)
	if err != nil {
		NewErrorResponse(w, http.StatusBadRequest,
			"bad credentials form")
		return
	}

	id, err := h.services.AuthorizationService.
		CreateUser(user)
	if err != nil {
		logrus.Warnf("Failed to create user '%s': %s",
			user.Username, err.Error())
		NewErrorResponse(w, http.StatusInternalServerError,
			"Unable to create user "+user.Username)
		return
	}

	logrus.Infof("User ID %d with username '%s' has been registered", id, user.Username)
	newJSONResponse(w, http.StatusOK, map[string]interface{}{
		"id": fmt.Sprint(id),
	})
}

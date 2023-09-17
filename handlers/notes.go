package handlers

import (
	"net/http"
	"strings"

	"github.com/ztx-lyghters/kode-notes/core"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetAll(w http.ResponseWriter,
	r *http.Request) {

	user_id, err := getUserId(w, r)
	if err != nil || user_id < 1 {
		NewErrorResponse(w, http.StatusInternalServerError,
			"cannot get user id")
		return
	}

	all_notes, err := h.services.NotesService.GetAll(user_id)
	if err != nil {
		NewErrorResponse(w, http.StatusInternalServerError,
			"cannot get note list: "+err.Error())
	}

	newJSONResponse(w, http.StatusOK, all_notes)
	logrus.Infof("Issued a list of notes to a user with ID '%d'",
		user_id)
}

func (h *Handler) Create(w http.ResponseWriter,
	r *http.Request) {

	user_id, err := getUserId(w, r)
	if err != nil {
		return
	}
	note := &core.Note{User_ID: user_id}
	err = parseJSON(r, note)
	if err != nil {
		NewErrorResponse(w, http.StatusBadRequest,
			"bad note form")
		return
	}

	reply, err := h.services.SpellerService.Yandex.Check([]string{
		note.Title,
		note.Description,
	})
	if err != nil {
		logrus.Error("Spellcheck error: " + err.Error())
	} else {
		h.services.SpellerService.Yandex.Fix(reply,
			&note.Title, &note.Description)
	}

	if strings.TrimSpace(note.Title) == "" {
		NewErrorResponse(w, http.StatusBadRequest,
			"empty note title")
		return

	}

	id, err := h.services.NotesService.Create(note)
	if err != nil {
		NewErrorResponse(w, http.StatusInternalServerError,
			err.Error())
		return
	}

	newJSONResponse(w, http.StatusOK, map[string]interface{}{
		"id": id,
	})
	logrus.Infof("User with ID %d has created a note with ID %d",
		user_id, id)
}

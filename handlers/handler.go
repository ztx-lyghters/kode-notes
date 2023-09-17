package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ztx-lyghters/kode-notes/config"
	"github.com/ztx-lyghters/kode-notes/services"
)

type Handler struct {
	services *services.Services
}

func New(s *services.Services, mux *http.ServeMux, app_cfg *config.Config) *Handler {
	s.SpellerService.Enabled = app_cfg.Spellcheck
	return &Handler{
		services: s,
	}
}

func throwNotFound(w http.ResponseWriter, r *http.Request) {
	newJSONResponse(w, http.StatusNotFound,
		map[string]interface{}{
			"message": "not found",
		})
}

func parseJSON(r *http.Request, entity interface{}) error {
	err := json.NewDecoder(r.Body).Decode(entity)
	if err != nil {
		return err
	}

	return nil
}

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

func NewErrorResponse(w http.ResponseWriter, status_code int, message string) {
	logrus.Error(message)
	newJSONResponse(w, status_code, map[string]interface{}{
		"message": message,
	})
}

func newJSONResponse(w http.ResponseWriter, status_code int, data interface{}) {
	w.Header().Set("Content-Type",
		"application/json; charset=utf-8",
	)

	w.WriteHeader(status_code)

	json_bytes, _ := json.Marshal(data)

	_, err := w.Write(json_bytes)
	if err != nil {
		logrus.Errorln("Reply error: " + err.Error())
	}
}

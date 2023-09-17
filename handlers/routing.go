package handlers

import (
	"net/http"
	"regexp"
)

var (
	AUTH_REGEX          = regexp.MustCompile(`^\/auth[\/]*`)
	API_REGEX           = regexp.MustCompile(`^\/api[\/]*`)
	SIGN_IN_REGEX       = regexp.MustCompile(`^\/auth/sign-in[\/]*$`)
	SIGN_UP_REGEX       = regexp.MustCompile(`^\/auth/sign-up[\/]*$`)
	GET_ALL_NOTES_REGEX = regexp.MustCompile(`^\/api/notes[\/]*$`)
	CREATE_NOTE_REGEX   = regexp.MustCompile(`^\/api/notes/new[\/]*$`)
)

func (h *Handler) ServeHTTP(w http.ResponseWriter,
	r *http.Request) {

	w.Header().Set("content-type", "application/json")
	defer r.Body.Close()

	switch {
	case AUTH_REGEX.MatchString(r.URL.Path):
		handleAuth(w, r, h)

	case API_REGEX.MatchString(r.URL.Path):
		handleAPI(w, r, h)

	default:
		throwNotFound(w, r)
	}
}

func handleAuth(w http.ResponseWriter, r *http.Request,
	h *Handler) {

	switch {
	case r.Method == http.MethodGet &&
		SIGN_IN_REGEX.MatchString(r.URL.Path):

		h.SignIn(w, r)

	case r.Method == http.MethodPost &&
		SIGN_UP_REGEX.MatchString(r.URL.Path):

		h.SignUp(w, r)

	default:
		throwNotFound(w, r)
	}
}

func handleAPI(w http.ResponseWriter, r *http.Request,
	h *Handler) {

	switch {
	case r.Method == http.MethodGet &&
		GET_ALL_NOTES_REGEX.MatchString(r.URL.Path):

		err := h.userIdentify(w, r)
		if err != nil {
			NewErrorResponse(w, http.StatusUnauthorized,
				err.Error())
			return
		}

		h.GetAll(w, r)

	case r.Method == http.MethodPost &&
		CREATE_NOTE_REGEX.MatchString(r.URL.Path):

		err := h.userIdentify(w, r)
		if err != nil {
			NewErrorResponse(w, http.StatusUnauthorized,
				err.Error())
			return
		}

		h.Create(w, r)

	default:
		throwNotFound(w, r)
	}
}

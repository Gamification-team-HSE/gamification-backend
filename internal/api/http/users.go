package http

import (
	"net/http"
	"strconv"
)

func (s *Server) postUserEvent(w http.ResponseWriter, r *http.Request) {
	eventID, err := strconv.Atoi(r.FormValue("event_id"))
	if err != nil {
		handleError(r.Context(), w, err)
		return
	}
	userEmail := r.FormValue("user_email")

	err = s.userService.AddEvent(r.Context(), userEmail, eventID)
	if err != nil {
		handleError(r.Context(), w, err)
		return
	}
	handleResponse(r.Context(), w, "ok")
}

func (s *Server) postUserStat(w http.ResponseWriter, r *http.Request) {
	statID, err := strconv.Atoi(r.FormValue("stat_id"))
	if err != nil {
		handleError(r.Context(), w, err)
		return
	}
	userEmail := r.FormValue("user_email")

	err = s.userService.AddStat(r.Context(), userEmail, statID)
	if err != nil {
		handleError(r.Context(), w, err)
		return
	}
	handleResponse(r.Context(), w, "ok")
}

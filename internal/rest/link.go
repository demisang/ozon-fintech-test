package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/demisang/ozon-fintech-test/internal/models"
)

type CreateLinkRequest struct {
	URL *string `json:"url"`
}

func (s *Server) linkCreate(w http.ResponseWriter, r *http.Request) {
	var request CreateLinkRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errResponse(w, r, http.StatusBadRequest, err)
		return
	}

	if request.URL == nil || *request.URL == "" {
		errResponse(w, r, http.StatusBadRequest, errors.New("URL required"))
		return
	}

	link, err := s.app.LinkCreate(r.Context(), models.CreateLinkDTO{URL: *request.URL})
	if err != nil {
		errResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	response(w, r, http.StatusCreated, link)
}

func (s *Server) linkGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errResponse(w, r, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}

	linkCode := strings.Trim(r.URL.Path, "/")
	if !s.app.ValidateLinkCodeLength(linkCode) {
		errResponse(w, r, http.StatusUnprocessableEntity, errors.New("link must be 10 symbols length"))
		return
	}

	if models.CompiledTemplate.MatchString(linkCode) {
		errResponse(w, r, http.StatusUnprocessableEntity, errors.New("link contain restricted symbols"))
		return
	}

	link, err := s.app.LinkGet(r.Context(), linkCode)

	switch {
	case errors.Is(err, models.ErrLinkNotFound):
		errResponse(w, r, http.StatusNotFound, err)
		return
	case err != nil:
		errResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	response(w, r, http.StatusOK, link)
}

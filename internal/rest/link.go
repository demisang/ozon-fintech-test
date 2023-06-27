package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/demisang/ozon-fintech-test/internal/models"
)

const shortLinkInvalidPattern = `[^a-zA-Z\d_]+`

type CreateLinkRequest struct {
	URL *string `json:"url"`
}

func (s *Server) linkCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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

	link, err := s.app.Storage.Create(ctx, models.CreateLinkDto{URL: *request.URL})

	response(w, r, 201, link)
}

func (s *Server) linkGet(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodGet {
		errResponse(w, r, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}

	linkCode := strings.Trim(r.URL.Path, "/")
	if len(linkCode) != s.app.Config.ShortLinkLength {
		errResponse(w, r, http.StatusUnprocessableEntity, errors.New("link must be 10 symbols length"))
		return
	}

	matchedInvalidSymbol, err := regexp.MatchString(shortLinkInvalidPattern, linkCode)
	if err != nil {
		errResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	if matchedInvalidSymbol {
		errResponse(w, r, http.StatusUnprocessableEntity, errors.New("link contain restricted symbols"))
		return
	}

	link, err := s.app.Storage.GetByCode(ctx, linkCode)

	switch {
	case errors.Is(err, models.ErrLinkNotFound):
		errResponse(w, r, http.StatusNotFound, err)
		return
	case err != nil:
		errResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	response(w, r, 200, link)
}

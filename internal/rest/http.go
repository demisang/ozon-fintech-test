package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"ozon-fintech-test/internal/application"
)

type Server struct {
	ctx    context.Context
	log    *logrus.Entry
	app    *application.App
	host   string
	port   int
	server *http.Server
}

func NewServer(log *logrus.Logger, app *application.App, host string, port int) *Server {
	return &Server{
		ctx:  context.Background(),
		log:  log.WithField("module", "rest"),
		app:  app,
		host: host,
		port: port,
	}
}

func (s *Server) Run(_ context.Context) error {
	http.HandleFunc("/create", s.linkCreate)
	http.HandleFunc("/", s.linkGet)

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", s.host, s.port), nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}

	return nil
}

func errResponse(w http.ResponseWriter, _ *http.Request, statusCode int, err error) {
	http.Error(w, err.Error(), statusCode)
}

func response(w http.ResponseWriter, r *http.Request, statusCode int, content any) {
	body, err := json.Marshal(content)
	if err != nil {
		errResponse(w, r, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(statusCode)
	_, err = io.WriteString(w, string(body))
	if err != nil {
		errResponse(w, r, http.StatusInternalServerError, err)
		return
	}
}

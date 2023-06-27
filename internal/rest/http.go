package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/demisang/ozon-fintech-test/internal/application"
	"github.com/sirupsen/logrus"
)

type Server struct {
	log    *logrus.Entry
	app    *application.App
	server *http.Server
}

func NewServer(log *logrus.Logger, app *application.App, host string, port int) *Server {
	s := Server{
		log: log.WithField("module", "rest"),
		app: app,
	}
	addr := fmt.Sprintf("%s:%d", host, port)
	mux := http.NewServeMux()
	mux.HandleFunc("/create", s.exampleMiddleware(http.HandlerFunc(s.linkCreate)).ServeHTTP)
	mux.HandleFunc("/", s.exampleMiddleware(http.HandlerFunc(s.linkGet)).ServeHTTP)

	s.server = &http.Server{Addr: addr, Handler: mux}

	return &s
}

func (s *Server) exampleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.log.Infof("request %s", r.RequestURI)

		next.ServeHTTP(w, r)
	})
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := s.server.Shutdown(ctx)
		if err != nil {
			s.log.Errorf("server shutdown: %v", err)
		}
	}()

	s.log.Infof("starting server %s", s.server.Addr)
	err := s.server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		s.log.Info("server closed")
	} else if err != nil {
		s.log.Infof("error starting server: %v", err)
	}

	return err
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

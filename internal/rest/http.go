package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

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

const (
	defaultReadHeaderTimeout = 30 * time.Second
	gracefulShutdownTimeout  = 10 * time.Second
)

func NewServer(log *logrus.Logger, app *application.App, host string, port int) *Server {
	server := Server{
		ctx:  context.Background(),
		log:  log.WithField("module", "rest"),
		app:  app,
		host: host,
		port: port,
	}

	http.HandleFunc("/create", server.linkCreate)
	http.HandleFunc("/", server.linkGet)

	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
	server.server = &http.Server{
		Addr:              fmt.Sprintf("%s:%d", host, port),
		ReadHeaderTimeout: defaultReadHeaderTimeout,
	}

	return &server
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
		defer cancel()

		if err := s.server.Shutdown(ctx); err != nil {
			s.log.Warnf("closing server: %v", err)
		}
	}()

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

func errResponse(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
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

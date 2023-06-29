package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/demisang/ozon-fintech-test/internal/models"
	"github.com/demisang/ozon-fintech-test/pkg/grpc/gen"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Клиент
// Сервер
// Интерфейс
//go:generate protoc -I=. --go_out=. --go-grpc_out=. links.proto

var _ gen.LinksServer = (*LinksServer)(nil)

type app interface {
	LinkGet(ctx context.Context, code string) (models.Link, error)
	LinkCreate(ctx context.Context, createDto models.CreateLinkDTO) (models.Link, error)
	ValidateLinkCodeLength(code string) bool
}

type LinksServer struct {
	gen.UnimplementedLinksServer
	app  app
	addr string
	log  *logrus.Entry
}

func NewLinksServer(app app, log *logrus.Logger, addr string) *LinksServer {
	return &LinksServer{
		app:  app,
		addr: addr,
		log:  log.WithField("module", "grpc_server"),
	}
}

func (s *LinksServer) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("grpc server run: %w", err)
	}

	grpcServer := grpc.NewServer(
	// grpc.UnaryInterceptor()
	)
	gen.RegisterLinksServer(grpcServer, s)

	go func() {
		<-ctx.Done()

		grpcServer.GracefulStop()
	}()

	s.log.Infof("listening grpc server %s", s.addr)

	return grpcServer.Serve(listener)
}

func (s *LinksServer) GetLink(ctx context.Context, params *gen.GetLinkParams) (*gen.Link, error) {
	linkCode := strings.Trim(params.Code, "/")
	if !s.app.ValidateLinkCodeLength(linkCode) {
		return nil, status.Errorf(codes.InvalidArgument, "link must be 10 symbols length")
	}

	if models.CompiledTemplate.MatchString(linkCode) {
		return nil, status.Errorf(codes.InvalidArgument, "link contain restricted symbols")
	}

	link, err := s.app.LinkGet(ctx, linkCode)

	switch {
	case errors.Is(err, models.ErrLinkNotFound):
		return nil, status.Errorf(codes.NotFound, "link not found")
	case err != nil:
		return nil, status.Errorf(codes.Internal, "get link: %v", err)
	}

	return responseLinkModel(link), nil
}

func (s *LinksServer) CreateLink(ctx context.Context, params *gen.CreateLinkParams) (*gen.Link, error) {
	if params.Url == "" {
		return nil, status.Errorf(codes.InvalidArgument, "URL required")
	}

	link, err := s.app.LinkCreate(ctx, models.CreateLinkDTO{URL: params.Url})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "link create: %v", err)
	}

	return responseLinkModel(link), nil
}

func responseLinkModel(link models.Link) *gen.Link {
	return &gen.Link{
		Code: link.Code,
		Url:  link.URL,
	}
}

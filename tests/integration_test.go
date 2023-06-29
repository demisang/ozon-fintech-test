package tests

import (
	"context"
	"testing"
	"time"

	"github.com/demisang/ozon-fintech-test/internal/application"
	"github.com/demisang/ozon-fintech-test/internal/config"
	"github.com/demisang/ozon-fintech-test/internal/store/memory"
	"github.com/demisang/ozon-fintech-test/pkg/grpc"
	"github.com/demisang/ozon-fintech-test/pkg/logger"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type GRPCTestSuite struct {
	suite.Suite
	ctx         context.Context
	log         *logrus.Logger
	app         *application.App
	memoryStore *memory.LinkStorage
	grpcClient  *grpc.LinksClient
	grpcServer  *grpc.LinksServer
}

func TestGRPC(t *testing.T) {
	suite.Run(t, new(GRPCTestSuite))
}

const (
	grpcAddr = "localhost:9001"
)

func (s *GRPCTestSuite) SetupSuite() {
	var err error

	s.ctx = context.Background()
	s.log = logger.GetLogger()
	s.memoryStore = memory.New(s.log)
	s.app = application.NewApp(s.log, config.Config{ShortLinkLength: 10}, s.memoryStore)
	s.grpcServer = grpc.NewLinksServer(s.app, s.log, grpcAddr)
	s.grpcClient, err = grpc.NewLinksClient(s.ctx, s.log, grpcAddr)
	s.Require().NoError(err)

	go func() {
		err = s.grpcServer.Run(s.ctx)
		s.Require().NoError(err)
	}()
	time.Sleep(100 * time.Millisecond)
}

func (s *GRPCTestSuite) TestGetLink() {
	url, code := "http://google.com/", "4sAtrvfFwX"

	s.Run("Test create link", func() {
		link, err := s.grpcClient.CreateLink(s.ctx, url)
		s.Require().NoError(err)
		s.Require().Equal(url, link.Url)
		s.Require().Equal(code, link.Code)
	})

	s.Run("Test get link by code", func() {
		link, err := s.grpcClient.GetLink(s.ctx, code)
		s.Require().NoError(err)
		s.Require().Equal(url, link.Url)
		s.Require().Equal(code, link.Code)
	})
}

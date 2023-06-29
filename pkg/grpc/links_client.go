package grpc

import (
	"context"
	"fmt"

	"github.com/demisang/ozon-fintech-test/pkg/grpc/gen"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// gRPC client
//go:generate protoc -I=. --go_out=. --go-grpc_out=. links.proto

// var _ gen.LinksClient = (*LinksClient)(nil)

type LinksClient struct {
	log        *logrus.Entry
	host       string
	opts       []grpc.DialOption
	grpcClient gen.LinksClient
}

func NewLinksClient(ctx context.Context, log *logrus.Logger, addr string) (*LinksClient, error) {
	c := LinksClient{
		log:  log.WithField("module", "grpc_client"),
		host: addr,
		opts: []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		},
	}

	conn, err := grpc.DialContext(ctx, c.host, c.opts...)
	if err != nil {
		return nil, fmt.Errorf("grpc client: %w", err)
	}

	c.grpcClient = gen.NewLinksClient(conn)

	return &c, nil
}

func (c *LinksClient) GetLink(ctx context.Context, code string) (*gen.Link, error) {
	link, err := c.grpcClient.GetLink(ctx, &gen.GetLinkParams{
		Code: code,
	})
	if err != nil {
		return nil, fmt.Errorf("get link: %w", err)
	}

	return link, nil
}

func (c *LinksClient) CreateLink(ctx context.Context, url string) (*gen.Link, error) {
	link, err := c.grpcClient.CreateLink(ctx, &gen.CreateLinkParams{
		Url: url,
	})
	if err != nil {
		return nil, fmt.Errorf("create link: %w", err)
	}

	return link, nil
}

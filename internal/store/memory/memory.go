package memory

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
)

type Memory struct {
	ctx   context.Context
	log   *logrus.Entry
	links map[string]string
}

var ErrNotExists = errors.New("not exists")

func New(ctx context.Context, log *logrus.Logger) (*Memory, error) {
	memory := Memory{
		ctx:   ctx,
		log:   log.WithField("module", "store"),
		links: make(map[string]string),
	}

	return &memory, nil
}

func (m *Memory) get(_ context.Context, key string) (string, error) {
	value, ok := m.links[key]

	if !ok {
		return "", ErrNotExists
	}

	return value, nil
}

func (m *Memory) set(_ context.Context, key string, value string) {
	m.links[key] = value
}

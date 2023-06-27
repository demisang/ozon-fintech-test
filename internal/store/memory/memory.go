package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/sirupsen/logrus"
)

type Memory struct {
	log   *logrus.Entry
	links map[string]string
	mu    sync.RWMutex
}

var ErrNotExists = errors.New("not exists")

func New(log *logrus.Logger) *LinkStorage {
	memory := Memory{
		log:   log.WithField("module", "store"),
		links: make(map[string]string),
	}

	return newLinkStorage(&memory)
}

func (m *Memory) get(_ context.Context, key string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok := m.links[key]

	if !ok {
		return "", ErrNotExists
	}

	return value, nil
}

func (m *Memory) set(_ context.Context, key string, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.links[key] = value
}

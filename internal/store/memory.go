package store

import (
	"sync"

	"github.com/mlawrence427/signalrelay/internal/envelope"
)

type Memory struct {
	mu        sync.RWMutex
	bySubject map[string]envelope.Envelope
}

func NewMemory() *Memory {
	return &Memory{
		bySubject: make(map[string]envelope.Envelope),
	}
}

func (s *Memory) Put(env envelope.Envelope) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.bySubject[env.Subject] = env
	return nil
}

func (s *Memory) Get(subject string) (envelope.Envelope, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	env, ok := s.bySubject[subject]
	return env, ok, nil
}

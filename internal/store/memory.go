package store

import (
	"sync"

	"github.com/mlawrence427/signalrelay/internal/envelope"
)

type Memory struct {
	mu        sync.RWMutex
	bySubject map[string]envelope.Envelope
	events    map[string]string
}

func NewMemory() *Memory {
	return &Memory{
		bySubject: make(map[string]envelope.Envelope),
		events:    make(map[string]string),
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

func (s *Memory) MarkEventSeen(sourceEventID string, subject string) (bool, string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existingSubject, ok := s.events[sourceEventID]
	if ok {
		return true, existingSubject, nil
	}

	s.events[sourceEventID] = subject
	return false, subject, nil
}

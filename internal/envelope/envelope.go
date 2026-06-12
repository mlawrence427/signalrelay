package envelope

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"
)

type Input struct {
	Source         string          `json:"source"`
	Subject        string          `json:"subject"`
	StateType      string          `json:"state_type"`
	ObservedAt     string          `json:"observed_at"`
	StaleAfter     string          `json:"stale_after"`
	SourceEventID  string          `json:"source_event_id"`
	SourceObjectID string          `json:"source_object_id"`
	Payload        json.RawMessage `json:"payload"`
}

type Envelope struct {
	Source         string          `json:"source"`
	Subject        string          `json:"subject"`
	StateType      string          `json:"state_type"`
	ObservedAt     time.Time       `json:"observed_at"`
	StaleAfter     time.Time       `json:"stale_after"`
	Freshness      string          `json:"freshness"`
	SourceEventID  string          `json:"source_event_id"`
	SourceObjectID string          `json:"source_object_id"`
	PayloadHash    string          `json:"payload_hash"`
	Payload        json.RawMessage `json:"payload"`
}

func FromInput(input Input) (Envelope, error) {
	observedAt, err := time.Parse(time.RFC3339, input.ObservedAt)
	if err != nil {
		return Envelope{}, err
	}

	staleAfter, err := time.Parse(time.RFC3339, input.StaleAfter)
	if err != nil {
		return Envelope{}, err
	}

	return Envelope{
		Source:         input.Source,
		Subject:        input.Subject,
		StateType:      input.StateType,
		ObservedAt:     observedAt,
		StaleAfter:     staleAfter,
		SourceEventID:  input.SourceEventID,
		SourceObjectID: input.SourceObjectID,
		PayloadHash:    HashPayload(input.Payload),
		Payload:        input.Payload,
	}, nil
}

func HashPayload(payload json.RawMessage) string {
	sum := sha256.Sum256(payload)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func (e Envelope) WithFreshness(now time.Time) Envelope {
	e.Freshness = "stale"
	if now.Before(e.StaleAfter) {
		e.Freshness = "fresh"
	}
	return e
}

package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/mlawrence427/signalrelay/internal/envelope"
	_ "modernc.org/sqlite"
)

type SQLite struct {
	db *sql.DB
}

func NewSQLite(path string) (*SQLite, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	store := &SQLite{db: db}
	if err := store.init(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}

	return store, nil
}

func (s *SQLite) Close() error {
	return s.db.Close()
}

func (s *SQLite) init(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS subscription_states (
	source TEXT NOT NULL,
	subject TEXT PRIMARY KEY,
	state_type TEXT NOT NULL,
	observed_at TEXT NOT NULL,
	stale_after TEXT NOT NULL,
	source_event_id TEXT NOT NULL,
	source_object_id TEXT NOT NULL,
	payload_hash TEXT NOT NULL,
	payload TEXT NOT NULL
)`)
	return err
}

func (s *SQLite) Put(env envelope.Envelope) error {
	_, err := s.db.ExecContext(
		context.Background(),
		`
INSERT INTO subscription_states (
	source,
	subject,
	state_type,
	observed_at,
	stale_after,
	source_event_id,
	source_object_id,
	payload_hash,
	payload
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(subject) DO UPDATE SET
	source = excluded.source,
	state_type = excluded.state_type,
	observed_at = excluded.observed_at,
	stale_after = excluded.stale_after,
	source_event_id = excluded.source_event_id,
	source_object_id = excluded.source_object_id,
	payload_hash = excluded.payload_hash,
	payload = excluded.payload`,
		env.Source,
		env.Subject,
		env.StateType,
		env.ObservedAt.UTC().Format(time.RFC3339Nano),
		env.StaleAfter.UTC().Format(time.RFC3339Nano),
		env.SourceEventID,
		env.SourceObjectID,
		env.PayloadHash,
		string(env.Payload),
	)
	return err
}

func (s *SQLite) Get(subject string) (envelope.Envelope, bool, error) {
	var env envelope.Envelope
	var observedAt string
	var staleAfter string
	var payload string

	err := s.db.QueryRowContext(
		context.Background(),
		`
SELECT
	source,
	subject,
	state_type,
	observed_at,
	stale_after,
	source_event_id,
	source_object_id,
	payload_hash,
	payload
FROM subscription_states
WHERE subject = ?`,
		subject,
	).Scan(
		&env.Source,
		&env.Subject,
		&env.StateType,
		&observedAt,
		&staleAfter,
		&env.SourceEventID,
		&env.SourceObjectID,
		&env.PayloadHash,
		&payload,
	)
	if err == sql.ErrNoRows {
		return envelope.Envelope{}, false, nil
	}
	if err != nil {
		return envelope.Envelope{}, false, err
	}

	env.ObservedAt, err = time.Parse(time.RFC3339Nano, observedAt)
	if err != nil {
		return envelope.Envelope{}, false, err
	}

	env.StaleAfter, err = time.Parse(time.RFC3339Nano, staleAfter)
	if err != nil {
		return envelope.Envelope{}, false, err
	}

	env.Payload = json.RawMessage(payload)
	return env, true, nil
}

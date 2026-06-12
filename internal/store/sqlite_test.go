package store

import (
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"github.com/mlawrence427/signalrelay/internal/envelope"
)

func TestSQLitePutAndGet(t *testing.T) {
	store := newTestSQLite(t)

	env := testEnvelope()
	if err := store.Put(env); err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	got, ok, err := store.Get(env.Subject)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if !ok {
		t.Fatal("Get() ok = false, want true")
	}

	assertEnvelopeEqual(t, got, env)
}

func TestSQLitePersistsAfterReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "signalrelay.db")

	first, err := NewSQLite(path)
	if err != nil {
		t.Fatalf("NewSQLite() error = %v", err)
	}

	env := testEnvelope()
	if err := first.Put(env); err != nil {
		t.Fatalf("Put() error = %v", err)
	}
	if err := first.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	second, err := NewSQLite(path)
	if err != nil {
		t.Fatalf("NewSQLite() reopen error = %v", err)
	}
	defer second.Close()

	got, ok, err := second.Get(env.Subject)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if !ok {
		t.Fatal("Get() ok = false, want true")
	}

	assertEnvelopeEqual(t, got, env)
}

func TestSQLiteMissingSubject(t *testing.T) {
	store := newTestSQLite(t)

	_, ok, err := store.Get("cus_missing")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if ok {
		t.Fatal("Get() ok = true, want false")
	}
}

func TestSQLitePreservesPayloadAndPayloadHash(t *testing.T) {
	store := newTestSQLite(t)

	env := testEnvelope()
	if err := store.Put(env); err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	got, ok, err := store.Get(env.Subject)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if !ok {
		t.Fatal("Get() ok = false, want true")
	}

	if string(got.Payload) != string(env.Payload) {
		t.Fatalf("Payload = %s, want %s", got.Payload, env.Payload)
	}
	if got.PayloadHash != env.PayloadHash {
		t.Fatalf("PayloadHash = %q, want %q", got.PayloadHash, env.PayloadHash)
	}
}

func TestSQLiteDoesNotPersistFreshnessAsTrustedState(t *testing.T) {
	store := newTestSQLite(t)

	env := testEnvelope()
	env.Freshness = "stale"
	if err := store.Put(env); err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	got, ok, err := store.Get(env.Subject)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if !ok {
		t.Fatal("Get() ok = false, want true")
	}
	if got.Freshness != "" {
		t.Fatalf("Freshness = %q, want empty", got.Freshness)
	}
}

func TestSQLiteMarkEventSeen(t *testing.T) {
	store := newTestSQLite(t)

	duplicate, subject, err := store.MarkEventSeen("evt_123", "cus_123")
	if err != nil {
		t.Fatalf("MarkEventSeen() error = %v", err)
	}
	if duplicate {
		t.Fatal("duplicate = true, want false")
	}
	if subject != "cus_123" {
		t.Fatalf("subject = %q, want %q", subject, "cus_123")
	}

	duplicate, subject, err = store.MarkEventSeen("evt_123", "cus_changed")
	if err != nil {
		t.Fatalf("MarkEventSeen() duplicate error = %v", err)
	}
	if !duplicate {
		t.Fatal("duplicate = false, want true")
	}
	if subject != "cus_123" {
		t.Fatalf("subject = %q, want original %q", subject, "cus_123")
	}
}

func TestSQLiteMarkEventSeenPersistsAfterReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "signalrelay.db")

	first, err := NewSQLite(path)
	if err != nil {
		t.Fatalf("NewSQLite() error = %v", err)
	}
	duplicate, subject, err := first.MarkEventSeen("evt_123", "cus_123")
	if err != nil {
		t.Fatalf("MarkEventSeen() error = %v", err)
	}
	if duplicate {
		t.Fatal("duplicate = true, want false")
	}
	if subject != "cus_123" {
		t.Fatalf("subject = %q, want %q", subject, "cus_123")
	}
	if err := first.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	second, err := NewSQLite(path)
	if err != nil {
		t.Fatalf("NewSQLite() reopen error = %v", err)
	}
	defer second.Close()

	duplicate, subject, err = second.MarkEventSeen("evt_123", "cus_changed")
	if err != nil {
		t.Fatalf("MarkEventSeen() duplicate error = %v", err)
	}
	if !duplicate {
		t.Fatal("duplicate = false, want true")
	}
	if subject != "cus_123" {
		t.Fatalf("subject = %q, want original %q", subject, "cus_123")
	}
}

func newTestSQLite(t *testing.T) *SQLite {
	t.Helper()

	store, err := NewSQLite(filepath.Join(t.TempDir(), "signalrelay.db"))
	if err != nil {
		t.Fatalf("NewSQLite() error = %v", err)
	}
	t.Cleanup(func() {
		if err := store.Close(); err != nil {
			t.Fatalf("Close() error = %v", err)
		}
	})

	return store
}

func testEnvelope() envelope.Envelope {
	payload := json.RawMessage(`{"customer_id":"cus_123","subscription_id":"sub_123","status":"active"}`)
	return envelope.Envelope{
		Source:         "stripe",
		Subject:        "cus_123",
		StateType:      "subscription",
		ObservedAt:     time.Date(2026, 6, 11, 18, 0, 0, 0, time.UTC),
		StaleAfter:     time.Date(2026, 6, 11, 18, 30, 0, 0, time.UTC),
		SourceEventID:  "evt_123",
		SourceObjectID: "sub_123",
		PayloadHash:    envelope.HashPayload(payload),
		Payload:        payload,
	}
}

func assertEnvelopeEqual(t *testing.T, got envelope.Envelope, want envelope.Envelope) {
	t.Helper()

	if got.Source != want.Source {
		t.Fatalf("Source = %q, want %q", got.Source, want.Source)
	}
	if got.Subject != want.Subject {
		t.Fatalf("Subject = %q, want %q", got.Subject, want.Subject)
	}
	if got.StateType != want.StateType {
		t.Fatalf("StateType = %q, want %q", got.StateType, want.StateType)
	}
	if !got.ObservedAt.Equal(want.ObservedAt) {
		t.Fatalf("ObservedAt = %s, want %s", got.ObservedAt, want.ObservedAt)
	}
	if !got.StaleAfter.Equal(want.StaleAfter) {
		t.Fatalf("StaleAfter = %s, want %s", got.StaleAfter, want.StaleAfter)
	}
	if got.SourceEventID != want.SourceEventID {
		t.Fatalf("SourceEventID = %q, want %q", got.SourceEventID, want.SourceEventID)
	}
	if got.SourceObjectID != want.SourceObjectID {
		t.Fatalf("SourceObjectID = %q, want %q", got.SourceObjectID, want.SourceObjectID)
	}
	if got.PayloadHash != want.PayloadHash {
		t.Fatalf("PayloadHash = %q, want %q", got.PayloadHash, want.PayloadHash)
	}
	if string(got.Payload) != string(want.Payload) {
		t.Fatalf("Payload = %s, want %s", got.Payload, want.Payload)
	}
}

package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mlawrence427/signalrelay/internal/envelope"
	"github.com/mlawrence427/signalrelay/internal/store"
)

func TestHealthz(t *testing.T) {
	srv := newTestServer(t, time.Date(2026, 6, 11, 12, 0, 0, 0, time.UTC))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	srv.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	assertJSONContentType(t, rec)
	assertJSONField(t, rec.Body.Bytes(), "ok", true)
	assertNoDecisionFields(t, rec.Body.String())
}

func TestPostStoresEnvelopeAndGetReturnsByCustomerID(t *testing.T) {
	now := time.Date(2026, 6, 11, 12, 0, 0, 0, time.UTC)
	srv := newTestServer(t, now)

	post := request(t, srv, http.MethodPost, "/v1/stripe/subscription-state", validEnvelopeBody("2099-01-01T00:00:00Z"))
	if post.Code != http.StatusCreated {
		t.Fatalf("POST status = %d, want %d: %s", post.Code, http.StatusCreated, post.Body.String())
	}
	assertJSONContentType(t, post)
	assertNoDecisionFields(t, post.Body.String())

	get := request(t, srv, http.MethodGet, "/v1/state/stripe/subscription?customer_id=cus_123", nil)
	if get.Code != http.StatusOK {
		t.Fatalf("GET status = %d, want %d: %s", get.Code, http.StatusOK, get.Body.String())
	}
	assertJSONContentType(t, get)

	assertJSONField(t, get.Body.Bytes(), "subject", "cus_123")
	assertJSONField(t, get.Body.Bytes(), "state_type", "subscription")
	assertJSONField(t, get.Body.Bytes(), "source_event_id", "evt_123")
	assertJSONField(t, get.Body.Bytes(), "source_object_id", "sub_123")
	assertJSONField(t, get.Body.Bytes(), "payload_hash", "sha256:05a74e86ab32a59be13ea1d7b6f9871f0f732eade78eecea260505e09cdb1599")
	assertNoDecisionFields(t, get.Body.String())
}

func TestGetComputesFreshnessFresh(t *testing.T) {
	now := time.Date(2026, 6, 11, 12, 0, 0, 0, time.UTC)
	srv := newTestServer(t, now)

	request(t, srv, http.MethodPost, "/v1/stripe/subscription-state", validEnvelopeBody("2026-06-11T12:01:00Z"))

	get := request(t, srv, http.MethodGet, "/v1/state/stripe/subscription?customer_id=cus_123", nil)
	if get.Code != http.StatusOK {
		t.Fatalf("GET status = %d, want %d: %s", get.Code, http.StatusOK, get.Body.String())
	}

	assertJSONContentType(t, get)
	assertJSONField(t, get.Body.Bytes(), "freshness", "fresh")
	assertNoDecisionFields(t, get.Body.String())
}

func TestGetComputesFreshnessStale(t *testing.T) {
	now := time.Date(2026, 6, 11, 12, 0, 0, 0, time.UTC)
	srv := newTestServer(t, now)

	request(t, srv, http.MethodPost, "/v1/stripe/subscription-state", validEnvelopeBody("2026-06-11T11:59:00Z"))

	get := request(t, srv, http.MethodGet, "/v1/state/stripe/subscription?customer_id=cus_123", nil)
	if get.Code != http.StatusOK {
		t.Fatalf("GET status = %d, want %d: %s", get.Code, http.StatusOK, get.Body.String())
	}

	assertJSONContentType(t, get)
	assertJSONField(t, get.Body.Bytes(), "freshness", "stale")
	assertNoDecisionFields(t, get.Body.String())
}

func TestGetComputesFreshnessFromSQLiteStoredEnvelope(t *testing.T) {
	now := time.Date(2026, 6, 11, 12, 0, 0, 0, time.UTC)

	sqliteStore, err := store.NewSQLite(t.TempDir() + "/signalrelay.db")
	if err != nil {
		t.Fatalf("NewSQLite() error = %v", err)
	}
	t.Cleanup(func() {
		if err := sqliteStore.Close(); err != nil {
			t.Fatalf("Close() error = %v", err)
		}
	})

	payload := json.RawMessage(`{"customer_id":"cus_123","subscription_id":"sub_123","status":"active"}`)
	env := envelope.Envelope{
		Source:         "stripe",
		Subject:        "cus_123",
		StateType:      "subscription",
		ObservedAt:     now.Add(-time.Minute),
		StaleAfter:     now.Add(time.Minute),
		Freshness:      "stale",
		SourceEventID:  "evt_123",
		SourceObjectID: "sub_123",
		PayloadHash:    envelope.HashPayload(payload),
		Payload:        payload,
	}
	if err := sqliteStore.Put(env); err != nil {
		t.Fatalf("Put() error = %v", err)
	}

	srv := New(sqliteStore)
	srv.now = func() time.Time { return now }

	get := request(t, srv, http.MethodGet, "/v1/state/stripe/subscription?customer_id=cus_123", nil)
	if get.Code != http.StatusOK {
		t.Fatalf("GET status = %d, want %d: %s", get.Code, http.StatusOK, get.Body.String())
	}

	assertJSONContentType(t, get)
	assertJSONField(t, get.Body.Bytes(), "freshness", "fresh")
	assertNoDecisionFields(t, get.Body.String())
}

func TestGetMissingCustomerReturns404(t *testing.T) {
	srv := newTestServer(t, time.Date(2026, 6, 11, 12, 0, 0, 0, time.UTC))

	rec := request(t, srv, http.MethodGet, "/v1/state/stripe/subscription?customer_id=cus_missing", nil)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}

	assertJSONContentType(t, rec)
	assertJSONField(t, rec.Body.Bytes(), "error", "subscription_state_missing")
	assertNoDecisionFields(t, rec.Body.String())
}

func TestInvalidPostCasesReturn400(t *testing.T) {
	cases := []struct {
		name      string
		body      string
		wantError string
	}{
		{
			name:      "invalid json",
			body:      `{`,
			wantError: "invalid_json",
		},
		{
			name:      "missing subject",
			body:      `{"source":"stripe","state_type":"subscription","observed_at":"2026-06-11T12:00:00Z","stale_after":"2026-06-11T12:01:00Z","source_event_id":"evt_123","source_object_id":"sub_123","payload":{"status":"active"}}`,
			wantError: "subject_required",
		},
		{
			name:      "missing state type",
			body:      `{"source":"stripe","subject":"cus_123","observed_at":"2026-06-11T12:00:00Z","stale_after":"2026-06-11T12:01:00Z","source_event_id":"evt_123","source_object_id":"sub_123","payload":{"status":"active"}}`,
			wantError: "state_type_required",
		},
		{
			name:      "missing observed at",
			body:      `{"source":"stripe","subject":"cus_123","state_type":"subscription","stale_after":"2026-06-11T12:01:00Z","source_event_id":"evt_123","source_object_id":"sub_123","payload":{"status":"active"}}`,
			wantError: "observed_at_required",
		},
		{
			name:      "invalid observed at",
			body:      `{"source":"stripe","subject":"cus_123","state_type":"subscription","observed_at":"not-a-time","stale_after":"2026-06-11T12:01:00Z","source_event_id":"evt_123","source_object_id":"sub_123","payload":{"status":"active"}}`,
			wantError: "observed_at_invalid",
		},
		{
			name:      "missing stale after",
			body:      `{"source":"stripe","subject":"cus_123","state_type":"subscription","observed_at":"2026-06-11T12:00:00Z","source_event_id":"evt_123","source_object_id":"sub_123","payload":{"status":"active"}}`,
			wantError: "stale_after_required",
		},
		{
			name:      "invalid stale after",
			body:      `{"source":"stripe","subject":"cus_123","state_type":"subscription","observed_at":"2026-06-11T12:00:00Z","stale_after":"not-a-time","source_event_id":"evt_123","source_object_id":"sub_123","payload":{"status":"active"}}`,
			wantError: "stale_after_invalid",
		},
		{
			name:      "missing payload",
			body:      `{"source":"stripe","subject":"cus_123","state_type":"subscription","observed_at":"2026-06-11T12:00:00Z","stale_after":"2026-06-11T12:01:00Z","source_event_id":"evt_123","source_object_id":"sub_123"}`,
			wantError: "payload_required",
		},
		{
			name:      "null payload",
			body:      `{"source":"stripe","subject":"cus_123","state_type":"subscription","observed_at":"2026-06-11T12:00:00Z","stale_after":"2026-06-11T12:01:00Z","source_event_id":"evt_123","source_object_id":"sub_123","payload":null}`,
			wantError: "payload_required",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srv := newTestServer(t, time.Date(2026, 6, 11, 12, 0, 0, 0, time.UTC))

			rec := request(t, srv, http.MethodPost, "/v1/stripe/subscription-state", strings.NewReader(tc.body))
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d: %s", rec.Code, http.StatusBadRequest, rec.Body.String())
			}

			assertJSONContentType(t, rec)
			assertJSONField(t, rec.Body.Bytes(), "error", tc.wantError)
			assertNoDecisionFields(t, rec.Body.String())
		})
	}
}

func newTestServer(t *testing.T, now time.Time) *Server {
	t.Helper()

	srv := New(store.NewMemory())
	srv.now = func() time.Time { return now }
	return srv
}

func request(t *testing.T, srv *Server, method string, target string, body *strings.Reader) *httptest.ResponseRecorder {
	t.Helper()

	var reqBody *strings.Reader
	if body == nil {
		reqBody = strings.NewReader("")
	} else {
		reqBody = body
	}

	req := httptest.NewRequest(method, target, reqBody)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	srv.Routes().ServeHTTP(rec, req)
	return rec
}

func validEnvelopeBody(staleAfter string) *strings.Reader {
	body := `{"source":"stripe","subject":"cus_123","state_type":"subscription","observed_at":"2026-06-11T12:00:00Z","stale_after":"` + staleAfter + `","source_event_id":"evt_123","source_object_id":"sub_123","payload":{"customer_id":"cus_123","subscription_id":"sub_123","status":"active"}}`
	return strings.NewReader(body)
}

func assertJSONField(t *testing.T, body []byte, field string, want any) {
	t.Helper()

	var got map[string]any
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("response is not JSON: %v: %s", err, string(body))
	}

	if got[field] != want {
		t.Fatalf("%s = %v, want %v", field, got[field], want)
	}
}

func assertJSONContentType(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()

	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q, want %q", got, "application/json")
	}
}

func assertNoDecisionFields(t *testing.T, body string) {
	t.Helper()

	if bytes.Contains([]byte(body), []byte("allowed")) || bytes.Contains([]byte(body), []byte("denied")) {
		t.Fatalf("response included an access decision: %s", body)
	}
}

package server

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mlawrence427/signalrelay/internal/envelope"
)

type Store interface {
	Put(envelope.Envelope) error
	Get(subject string) (envelope.Envelope, bool, error)
	MarkEventSeen(sourceEventID string, subject string) (bool, string, error)
}

type Server struct {
	store                    Store
	now                      func() time.Time
	stripeStaleAfterWindow   time.Duration
	stripeWebhookSecret      string
	stripeSignatureTolerance time.Duration
}

func New(store Store) *Server {
	return NewWithStripeStaleAfter(store, 300*time.Second)
}

func NewWithStripeStaleAfter(store Store, staleAfterWindow time.Duration) *Server {
	return NewWithStripeConfig(store, staleAfterWindow, "", 300*time.Second)
}

func NewWithStripeConfig(store Store, staleAfterWindow time.Duration, webhookSecret string, signatureTolerance time.Duration) *Server {
	return &Server{
		store:                    store,
		now:                      time.Now,
		stripeStaleAfterWindow:   staleAfterWindow,
		stripeWebhookSecret:      webhookSecret,
		stripeSignatureTolerance: signatureTolerance,
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.handleHealthz)
	mux.HandleFunc("POST /v1/stripe/subscription-state", s.handlePostSubscriptionState)
	mux.HandleFunc("POST /v1/stripe/events", s.handlePostStripeEvent)
	mux.HandleFunc("POST /v1/stripe/webhook", s.handlePostStripeWebhook)
	mux.HandleFunc("GET /v1/state/stripe/subscription", s.handleGetSubscriptionState)
	return mux
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handlePostSubscriptionState(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var input envelope.Input
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}

	if err := validateInput(input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	env, err := envelope.FromInput(input)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_envelope")
		return
	}

	if err := s.store.Put(env); err != nil {
		writeError(w, http.StatusInternalServerError, "store_write_failed")
		return
	}

	writeJSON(w, http.StatusCreated, env.WithFreshness(s.now()))
}

type stripeEvent struct {
	ID      string          `json:"id"`
	Type    string          `json:"type"`
	Created json.RawMessage `json:"created"`
	Data    struct {
		Object json.RawMessage `json:"object"`
	} `json:"data"`
}

type stripeSubscriptionObject struct {
	ID       string `json:"id"`
	Object   string `json:"object"`
	Customer string `json:"customer"`
	Status   string `json:"status"`
}

func (s *Server) handlePostStripeEvent(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}

	s.ingestStripeEvent(w, rawBody)
}

func (s *Server) handlePostStripeWebhook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if s.stripeWebhookSecret == "" {
		writeError(w, http.StatusBadRequest, "stripe_webhook_secret_not_configured")
		return
	}

	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}

	if err := s.verifyStripeSignature(r.Header.Get("Stripe-Signature"), rawBody); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	s.ingestStripeEvent(w, rawBody)
}

func (s *Server) ingestStripeEvent(w http.ResponseWriter, rawBody []byte) {
	var event stripeEvent
	if err := json.Unmarshal(rawBody, &event); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}

	env, err := s.envelopeFromStripeEvent(event)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	duplicate, subject, err := s.store.MarkEventSeen(env.SourceEventID, env.Subject)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "store_event_failed")
		return
	}
	if duplicate {
		writeJSON(w, http.StatusOK, map[string]any{
			"duplicate":       true,
			"source_event_id": env.SourceEventID,
			"subject":         subject,
		})
		return
	}

	if err := s.store.Put(env); err != nil {
		writeError(w, http.StatusInternalServerError, "store_write_failed")
		return
	}

	writeJSON(w, http.StatusCreated, env.WithFreshness(s.now()))
}

func (s *Server) verifyStripeSignature(header string, rawBody []byte) error {
	if header == "" {
		return errors.New("stripe_signature_header_required")
	}

	timestamp, signatures, err := parseStripeSignatureHeader(header)
	if err != nil {
		return err
	}

	signedAt := time.Unix(timestamp, 0)
	if s.now().Sub(signedAt) > s.stripeSignatureTolerance || signedAt.Sub(s.now()) > s.stripeSignatureTolerance {
		return errors.New("stripe_signature_timestamp_outside_tolerance")
	}

	if len(signatures) == 0 {
		return errors.New("stripe_signature_v1_required")
	}

	expected := computeStripeSignature(s.stripeWebhookSecret, timestamp, rawBody)
	for _, signature := range signatures {
		if signatureBytes, err := hex.DecodeString(signature); err == nil && hmac.Equal(signatureBytes, expected) {
			return nil
		}
	}

	return errors.New("stripe_signature_mismatch")
}

func parseStripeSignatureHeader(header string) (int64, []string, error) {
	var timestamp int64
	var signatures []string

	for _, part := range strings.Split(header, ",") {
		name, value, ok := strings.Cut(strings.TrimSpace(part), "=")
		if !ok {
			continue
		}

		switch name {
		case "t":
			parsed, err := strconv.ParseInt(value, 10, 64)
			if err != nil || parsed <= 0 {
				return 0, nil, errors.New("stripe_signature_timestamp_invalid")
			}
			timestamp = parsed
		case "v1":
			if value != "" {
				signatures = append(signatures, value)
			}
		}
	}

	if timestamp == 0 {
		return 0, nil, errors.New("stripe_signature_timestamp_invalid")
	}

	return timestamp, signatures, nil
}

func computeStripeSignature(secret string, timestamp int64, rawBody []byte) []byte {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(fmt.Sprintf("%d.", timestamp)))
	_, _ = mac.Write(rawBody)
	return mac.Sum(nil)
}

func (s *Server) envelopeFromStripeEvent(event stripeEvent) (envelope.Envelope, error) {
	if err := validateStripeEvent(event); err != nil {
		return envelope.Envelope{}, err
	}

	created, err := parseStripeEventCreated(event.Created)
	if err != nil {
		return envelope.Envelope{}, err
	}

	var subscription stripeSubscriptionObject
	if err := json.Unmarshal(event.Data.Object, &subscription); err != nil {
		return envelope.Envelope{}, errors.New("stripe_event_object_invalid")
	}
	if subscription.Object == "" {
		return envelope.Envelope{}, errors.New("stripe_subscription_object_required")
	}
	if subscription.Object != "subscription" {
		return envelope.Envelope{}, errors.New("stripe_subscription_object_invalid")
	}
	if subscription.Customer == "" {
		return envelope.Envelope{}, errors.New("stripe_subscription_customer_required")
	}
	if subscription.ID == "" {
		return envelope.Envelope{}, errors.New("stripe_subscription_id_required")
	}
	if subscription.Status == "" {
		return envelope.Envelope{}, errors.New("stripe_subscription_status_required")
	}

	observedAt := time.Unix(created, 0).UTC()
	return envelope.Envelope{
		Source:         "stripe",
		Subject:        subscription.Customer,
		StateType:      "stripe.subscription",
		ObservedAt:     observedAt,
		StaleAfter:     observedAt.Add(s.stripeStaleAfterWindow),
		SourceEventID:  event.ID,
		SourceObjectID: subscription.ID,
		PayloadHash:    envelope.HashPayload(event.Data.Object),
		Payload:        event.Data.Object,
	}, nil
}

func (s *Server) handleGetSubscriptionState(w http.ResponseWriter, r *http.Request) {
	customerID := r.URL.Query().Get("customer_id")
	if customerID == "" {
		writeError(w, http.StatusBadRequest, "customer_id_required")
		return
	}

	env, ok, err := s.store.Get(customerID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "store_read_failed")
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "subscription_state_missing")
		return
	}

	writeJSON(w, http.StatusOK, env.WithFreshness(s.now()))
}

func validateInput(input envelope.Input) error {
	switch {
	case input.Source == "":
		return errors.New("source_required")
	case input.Subject == "":
		return errors.New("subject_required")
	case input.StateType == "":
		return errors.New("state_type_required")
	case input.ObservedAt == "":
		return errors.New("observed_at_required")
	case !validTimestamp(input.ObservedAt):
		return errors.New("observed_at_invalid")
	case input.StaleAfter == "":
		return errors.New("stale_after_required")
	case !validTimestamp(input.StaleAfter):
		return errors.New("stale_after_invalid")
	case input.SourceEventID == "":
		return errors.New("source_event_id_required")
	case input.SourceObjectID == "":
		return errors.New("source_object_id_required")
	case len(input.Payload) == 0 || bytes.Equal(bytes.TrimSpace(input.Payload), []byte("null")):
		return errors.New("payload_required")
	default:
		return nil
	}
}

func validateStripeEvent(event stripeEvent) error {
	switch {
	case event.ID == "":
		return errors.New("stripe_event_id_required")
	case event.Type == "":
		return errors.New("stripe_event_type_required")
	case !supportedStripeSubscriptionEvent(event.Type):
		return errors.New("unsupported_stripe_event_type")
	case len(event.Created) == 0 || bytes.Equal(bytes.TrimSpace(event.Created), []byte("null")):
		return errors.New("stripe_event_created_required")
	case len(event.Data.Object) == 0 || bytes.Equal(bytes.TrimSpace(event.Data.Object), []byte("null")):
		return errors.New("stripe_event_object_required")
	default:
		return nil
	}
}

func parseStripeEventCreated(value json.RawMessage) (int64, error) {
	var created int64
	if err := json.Unmarshal(value, &created); err != nil {
		return 0, errors.New("stripe_event_created_invalid")
	}
	if created <= 0 {
		return 0, errors.New("stripe_event_created_invalid")
	}
	return created, nil
}

func supportedStripeSubscriptionEvent(eventType string) bool {
	switch eventType {
	case "customer.subscription.created",
		"customer.subscription.updated",
		"customer.subscription.deleted":
		return true
	default:
		return false
	}
}

func validTimestamp(value string) bool {
	_, err := time.Parse(time.RFC3339, value)
	return err == nil
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, code string) {
	writeJSON(w, status, map[string]string{"error": code})
}

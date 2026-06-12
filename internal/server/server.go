package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/mlawrence427/signalrelay/internal/envelope"
)

type Store interface {
	Put(envelope.Envelope)
	Get(subject string) (envelope.Envelope, bool)
}

type Server struct {
	store Store
	now   func() time.Time
}

func New(store Store) *Server {
	return &Server{
		store: store,
		now:   time.Now,
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.handleHealthz)
	mux.HandleFunc("POST /v1/stripe/subscription-state", s.handlePostSubscriptionState)
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

	env := envelope.FromInput(input)
	s.store.Put(env)

	writeJSON(w, http.StatusCreated, env.WithFreshness(s.now()))
}

func (s *Server) handleGetSubscriptionState(w http.ResponseWriter, r *http.Request) {
	customerID := r.URL.Query().Get("customer_id")
	if customerID == "" {
		writeError(w, http.StatusBadRequest, "customer_id_required")
		return
	}

	env, ok := s.store.Get(customerID)
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
	case input.ObservedAt.IsZero():
		return errors.New("observed_at_required")
	case input.StaleAfter.IsZero():
		return errors.New("stale_after_required")
	case input.SourceEventID == "":
		return errors.New("source_event_id_required")
	case input.SourceObjectID == "":
		return errors.New("source_object_id_required")
	case len(input.Payload) == 0:
		return errors.New("payload_required")
	default:
		return nil
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, code string) {
	writeJSON(w, status, map[string]string{"error": code})
}

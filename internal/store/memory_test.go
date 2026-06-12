package store

import "testing"

func TestMemoryMarkEventSeen(t *testing.T) {
	store := NewMemory()

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

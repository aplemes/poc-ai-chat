package services

import (
	"testing"
	"time"

	"gin-quickstart/models"
)

// ---------------------------------------------------------------------------
// GetOrCreate
// ---------------------------------------------------------------------------

func TestGetOrCreate_NewSession(t *testing.T) {
	svc := NewConversationService()

	session := svc.GetOrCreate("")
	if session == nil {
		t.Fatal("expected a new session, got nil")
	}
	if session.ID == "" {
		t.Error("expected a non-empty session ID")
	}
	if session.Status != models.SessionStatusCollecting {
		t.Errorf("expected status %q, got %q", models.SessionStatusCollecting, session.Status)
	}
	if len(session.Messages) != 0 {
		t.Errorf("expected empty messages, got %d", len(session.Messages))
	}
}

func TestGetOrCreate_ExistingSession(t *testing.T) {
	svc := NewConversationService()

	first := svc.GetOrCreate("")
	id := first.ID

	second := svc.GetOrCreate(id)
	if second.ID != id {
		t.Errorf("expected same session ID %q, got %q", id, second.ID)
	}
}

func TestGetOrCreate_UnknownIDCreatesNew(t *testing.T) {
	svc := NewConversationService()

	session := svc.GetOrCreate("unknown-id-xyz")
	if session == nil {
		t.Fatal("expected a session for unknown ID, got nil")
	}
	// The requested ID is used as the new session ID.
	if session.ID != "unknown-id-xyz" {
		t.Errorf("expected ID %q, got %q", "unknown-id-xyz", session.ID)
	}
}

// ---------------------------------------------------------------------------
// SessionExists
// ---------------------------------------------------------------------------

func TestSessionExists(t *testing.T) {
	svc := NewConversationService()

	session := svc.GetOrCreate("")
	if !svc.SessionExists(session.ID) {
		t.Error("expected session to exist after GetOrCreate")
	}
	if svc.SessionExists("nonexistent") {
		t.Error("expected nonexistent session to not exist")
	}
}

// ---------------------------------------------------------------------------
// AddMessage / GetMessages
// ---------------------------------------------------------------------------

func TestAddMessage_AppendsToSession(t *testing.T) {
	svc := NewConversationService()
	session := svc.GetOrCreate("")

	svc.AddMessage(session.ID, models.Message{Role: models.RoleUser, Content: "hello"})
	svc.AddMessage(session.ID, models.Message{Role: models.RoleAssistant, Content: "hi"})

	msgs := svc.GetMessages(session.ID)
	if len(msgs) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(msgs))
	}
	if msgs[0].Content != "hello" {
		t.Errorf("expected first message %q, got %q", "hello", msgs[0].Content)
	}
	if msgs[1].Content != "hi" {
		t.Errorf("expected second message %q, got %q", "hi", msgs[1].Content)
	}
}

func TestAddMessage_UnknownSessionIsNoOp(t *testing.T) {
	svc := NewConversationService()
	// Should not panic or error — just silently ignored.
	svc.AddMessage("ghost-session", models.Message{Role: models.RoleUser, Content: "test"})
}

func TestGetMessages_ReturnsNilForUnknownSession(t *testing.T) {
	svc := NewConversationService()
	msgs := svc.GetMessages("does-not-exist")
	if msgs != nil {
		t.Errorf("expected nil for unknown session, got %v", msgs)
	}
}

func TestGetMessages_ReturnsCopy(t *testing.T) {
	svc := NewConversationService()
	session := svc.GetOrCreate("")
	svc.AddMessage(session.ID, models.Message{Role: models.RoleUser, Content: "original"})

	msgs := svc.GetMessages(session.ID)
	msgs[0].Content = "mutated"

	// The internal slice should not have been mutated.
	msgs2 := svc.GetMessages(session.ID)
	if msgs2[0].Content != "original" {
		t.Errorf("GetMessages returned a reference instead of a copy: got %q", msgs2[0].Content)
	}
}

// ---------------------------------------------------------------------------
// PendingFormData
// ---------------------------------------------------------------------------

func TestSetAndTakePendingFormData(t *testing.T) {
	svc := NewConversationService()
	session := svc.GetOrCreate("")

	data := &models.FormFillData{Title: "Test demand"}
	svc.SetPendingFormData(session.ID, data)

	taken, exists := svc.TakeAndClearPendingFormData(session.ID)
	if !exists {
		t.Fatal("expected session to exist")
	}
	if taken == nil {
		t.Fatal("expected pending form data, got nil")
	}
	if taken.Title != "Test demand" {
		t.Errorf("expected title %q, got %q", "Test demand", taken.Title)
	}

	// Second take should return nil (cleared after first take).
	taken2, exists2 := svc.TakeAndClearPendingFormData(session.ID)
	if !exists2 {
		t.Fatal("session should still exist after clearing pending data")
	}
	if taken2 != nil {
		t.Errorf("expected nil after clear, got %v", taken2)
	}
}

func TestTakePendingFormData_UnknownSession(t *testing.T) {
	svc := NewConversationService()
	taken, exists := svc.TakeAndClearPendingFormData("ghost")
	if exists {
		t.Error("expected exists=false for unknown session")
	}
	if taken != nil {
		t.Errorf("expected nil data for unknown session, got %v", taken)
	}
}

func TestSetPendingFormData_UnknownSessionIsNoOp(t *testing.T) {
	svc := NewConversationService()
	// Should not panic.
	svc.SetPendingFormData("ghost", &models.FormFillData{Title: "x"})
}

// ---------------------------------------------------------------------------
// evictExpired
// ---------------------------------------------------------------------------

func TestEvictExpired_RemovesOldSessions(t *testing.T) {
	svc := NewConversationService()
	session := svc.GetOrCreate("")

	// Backdating lastAccess so eviction considers this session expired.
	svc.mu.Lock()
	svc.sessions[session.ID].lastAccess = time.Now().Add(-(sessionTTL + time.Minute))
	svc.mu.Unlock()

	svc.evictExpired()

	if svc.SessionExists(session.ID) {
		t.Error("expected expired session to be evicted")
	}
}

func TestEvictExpired_KeepsRecentSessions(t *testing.T) {
	svc := NewConversationService()
	session := svc.GetOrCreate("")

	svc.evictExpired()

	if !svc.SessionExists(session.ID) {
		t.Error("expected recent session to survive eviction")
	}
}

// ---------------------------------------------------------------------------
// generateID
// ---------------------------------------------------------------------------

func TestGenerateID_Uniqueness(t *testing.T) {
	ids := make(map[string]struct{}, 100)
	for i := 0; i < 100; i++ {
		id := generateID()
		if id == "" {
			t.Fatal("generateID returned empty string")
		}
		if _, dup := ids[id]; dup {
			t.Fatalf("generateID produced duplicate ID: %q", id)
		}
		ids[id] = struct{}{}
	}
}

func TestGenerateID_Length(t *testing.T) {
	id := generateID()
	// 16 bytes hex-encoded = 32 characters
	if len(id) != 32 {
		t.Errorf("expected ID length 32, got %d (%q)", len(id), id)
	}
}

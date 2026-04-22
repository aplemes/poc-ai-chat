package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"gin-quickstart/models"
)

const sessionTTL = 2 * time.Hour
const cleanupInterval = 15 * time.Minute

type sessionEntry struct {
	session    *models.Session
	lastAccess time.Time
}

type ConversationService struct {
	sessions map[string]*sessionEntry
	mu       sync.Mutex
}

func NewConversationService() *ConversationService {
	svc := &ConversationService{
		sessions: make(map[string]*sessionEntry),
	}
	go svc.cleanupLoop()
	return svc
}

// GetOrCreate returns an existing session or creates a new one.
// Callers must not retain the returned pointer beyond the current request;
// use the service methods below for all mutations.
func (s *ConversationService) GetOrCreate(sessionID string) *models.Session {
	s.mu.Lock()
	defer s.mu.Unlock()

	if sessionID != "" {
		if entry, ok := s.sessions[sessionID]; ok {
			entry.lastAccess = time.Now()
			return entry.session
		}
	}

	if sessionID == "" {
		sessionID = generateID()
	}

	session := &models.Session{
		ID:       sessionID,
		Messages: []models.Message{},
		Status:   models.SessionStatusCollecting,
	}
	s.sessions[sessionID] = &sessionEntry{session: session, lastAccess: time.Now()}
	return session
}

// SessionExists reports whether a session with the given ID is known.
func (s *ConversationService) SessionExists(sessionID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.sessions[sessionID]
	return ok
}

// GetMessages returns a shallow copy of the session's message slice.
// Returns nil if the session does not exist.
func (s *ConversationService) GetMessages(sessionID string) []models.Message {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.sessions[sessionID]
	if !ok {
		return nil
	}
	entry.lastAccess = time.Now()
	msgs := make([]models.Message, len(entry.session.Messages))
	copy(msgs, entry.session.Messages)
	return msgs
}

func (s *ConversationService) AddMessage(sessionID string, msg models.Message) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if entry, ok := s.sessions[sessionID]; ok {
		entry.session.Messages = append(entry.session.Messages, msg)
		entry.lastAccess = time.Now()
	}
}

// SetPendingFormData stores proposed form data awaiting user confirmation (BE-01).
func (s *ConversationService) SetPendingFormData(sessionID string, data *models.FormFillData) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if entry, ok := s.sessions[sessionID]; ok {
		entry.session.PendingFormData = data
	}
}

// TakeAndClearPendingFormData atomically reads and removes the pending form data (BE-01/BE-05).
// Returns (nil, false) when the session does not exist; (nil, true) when it exists but has no pending data.
func (s *ConversationService) TakeAndClearPendingFormData(sessionID string) (*models.FormFillData, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.sessions[sessionID]
	if !ok {
		return nil, false
	}
	data := entry.session.PendingFormData
	entry.session.PendingFormData = nil
	return data, true
}

func (s *ConversationService) cleanupLoop() {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()
	for range ticker.C {
		s.evictExpired()
	}
}

func (s *ConversationService) evictExpired() {
	s.mu.Lock()
	defer s.mu.Unlock()
	cutoff := time.Now().Add(-sessionTTL)
	for id, entry := range s.sessions {
		if entry.lastAccess.Before(cutoff) {
			delete(s.sessions, id)
		}
	}
}

func generateID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// crypto/rand failure is unrecoverable in any meaningful way; panic to surface immediately.
		panic(fmt.Sprintf("crypto/rand.Read failed: %v", err))
	}
	return hex.EncodeToString(b)
}

package services

import (
	"crypto/rand"
	"encoding/hex"
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
	mu       sync.RWMutex
}

func NewConversationService() *ConversationService {
	svc := &ConversationService{
		sessions: make(map[string]*sessionEntry),
	}
	go svc.cleanupLoop()
	return svc
}

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
		Status:   "collecting",
	}
	s.sessions[sessionID] = &sessionEntry{session: session, lastAccess: time.Now()}
	return session
}

func (s *ConversationService) AddMessage(sessionID string, msg models.Message) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if entry, ok := s.sessions[sessionID]; ok {
		entry.session.Messages = append(entry.session.Messages, msg)
		entry.lastAccess = time.Now()
	}
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
	rand.Read(b)
	return hex.EncodeToString(b)
}

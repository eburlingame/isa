package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"gopkg.in/olahol/melody.v1"
)

type PersistentSession struct {
	SessionId  string    `json:"sessionId"`
	PlayerName string    `json:"playerName"`
	GameHost   bool      `json:"gameHost"`
	ActiveGame string    `json:"activeGame"`
	UnsubChan  chan bool `json:"-"`
}

func NewSessionId() string {
	uuidWithHyphen := uuid.New()
	uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	return uuid
}

func NewPersistentSession() *PersistentSession {
	sessionId := NewSessionId()
	log.Printf("Creating new sessionId %s", sessionId)

	return &PersistentSession{
		GameHost:   false,
		SessionId:  sessionId,
		PlayerName: "",
		ActiveGame: "",
		UnsubChan:  nil,
	}
}

// GetPersistentSession retrieves the current Session object from the melody session
func GetPersistentSession(session *melody.Session) (*PersistentSession, error) {
	stored, exists := session.Get("persistentSession")
	if !exists {
		return nil, errors.New("Session is not stored")
	}

	persistentSession := stored.(PersistentSession)
	return &persistentSession, nil
}

// FetchPersistentSession retrieves the current session from Redis if it exists, if not it returns a new empty
// session
func FetchPersistentSession(ctx *context.Context, session *melody.Session, rdb *redis.Client, sessionId string) *PersistentSession {
	stored, err := rdb.Get(*ctx, "sessions:"+sessionId).Result()
	if err != nil {
		return NewPersistentSession()
	}

	var persistentSession PersistentSession

	err = json.Unmarshal([]byte(stored), &persistentSession)
	if err != nil {
		log.Printf("Error unmarshalling Redis response, %s\n", err)
		return NewPersistentSession()
	}

	session.Set("persistentSession", persistentSession)

	return &persistentSession
}

// SetPersistentSession update the current melody session value and stores it in Redis
func SetPersistentSession(ctx *context.Context, session *melody.Session, rdb *redis.Client, persistentSession *PersistentSession) {
	session.Set("persistentSession", *persistentSession)

	payload, err := json.Marshal(*persistentSession)
	if err != nil {
		log.Printf("Error marshalling session, %s\n", err)
		return
	}

	rdb.Set(*ctx, "sessions:"+persistentSession.SessionId, payload, 12*time.Hour)
}

package session

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Michaelpalacce/gobi/pkg/client"
	"github.com/Michaelpalacce/gobi/pkg/models"
	"github.com/Michaelpalacce/gobi/pkg/redis"
	"github.com/google/uuid"
)

var ExpirationTime = time.Hour * 24 * 7

type Session struct {
	SessionID string                 `json:"session_id"`
	Client    *client.ClientMetadata `json:"client"`
	User      *models.User           `json:"user"`
}

// NewSession will instantiate a new Session
func NewSession(client *client.ClientMetadata, user *models.User) *Session {
	session := &Session{
		SessionID: uuid.New().String(),
		Client:    client,
		User:      user,
	}

	session.Update()

	return session
}

// Update will update the session in redis
func (s *Session) Update() {
	redis.Set(fmt.Sprintf("%s-%s-%s", s.User.Username, s.Client.VaultName, s.SessionID), s.Encode(), ExpirationTime)
}

// Encode will encode the session into a string
// Use this when you want to store the session data in redis
func (s *Session) Encode() string {
	sessionBytes, err := json.Marshal(s)
	if err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(sessionBytes)
}

// RestoreSession will restore a session from an encoded string
// Use this when you retrieve the session data from redis
func RestoreSession(encoded string) (*Session, error) {
	sessionBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	var session Session
	err = json.Unmarshal(sessionBytes, &session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

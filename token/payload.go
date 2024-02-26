package token

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

// Different types of error returned by the VerifyToken function
var (
	ErrorExpiredToken = errors.New("token has expired")
	ErrInvalidToken   = errors.New("token is invalid")
)

// Payload contains the payload data of the token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issuedAt"`
	ExpiredAt time.Time `json:"expiredAt"`
}

// NewPayload creates a new token payload with a specific username and duration
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

// Valid checks if the token payload is valid or not
func (paylod *Payload) Valid() error {
	if time.Now().After(paylod.ExpiredAt) {
		return ErrorExpiredToken
	}
	return nil
}

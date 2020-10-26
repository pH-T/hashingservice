package service

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type database interface {
	Set(context.Context, string) error
	Exists(context.Context, string) (bool, error)
}

// NewBcryptHashService
func NewBcryptHashService(store database) *bcrypthashservice {
	return &bcrypthashservice{store: store}
}

type bcrypthashservice struct {
	store database
}

// Hash creates bcrypt of the given password
func (hs *bcrypthashservice) Hash(ctx context.Context, pw string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pw), 14)
	if err != nil {
		return "", fmt.Errorf("error creating a bcrypt hash: %w", err)
	}

	hash := string(bytes)

	err = hs.store.Set(ctx, hash)
	if err != nil {
		return "", fmt.Errorf("error storing a bcrypt hash: %w", err)
	}

	return hash, err
}

// Verify checks if the hash matches the given password and if this hash was already seen by this service
func (hs *bcrypthashservice) Verify(ctx context.Context, pw, hash string) (bool, bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
	if err != nil {
		return false, false, fmt.Errorf("error verifying a bcrypt hash: %w", err)
	}

	exists, err := hs.store.Exists(ctx, hash)
	if err != nil {
		return true, false, fmt.Errorf("error checking is bcrypt hash exists: %w", err)
	}

	if exists {
		return true, true, nil
	}
	return true, false, nil
}

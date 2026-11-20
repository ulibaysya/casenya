package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func (s *Service) GenerateUserOpaqueToken(ctx context.Context, id int64) (string, error) {
	tokenBytes := make([]byte, 32)
	rand.Read(tokenBytes)

	tokenHash := s.hasher.Hash(tokenBytes)
	err := s.strg.SaveUserOpaqueToken(ctx, id, tokenHash)
	if err != nil {
		return "", fmt.Errorf("saving token to storage: %w", err)
	}

	tokenRepresentation := base64.RawURLEncoding.EncodeToString(tokenBytes)

	return tokenRepresentation, nil
}

func (s *Service) ExchangeUserOpaqueToken(ctx context.Context, token string) (int64, error) {
	tokenBytes, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return 0, fmt.Errorf("invalid token format")
	}

	tokenHash := s.hasher.Hash(tokenBytes)
	userID, err := s.strg.ExchangeUserOpaqueToken(ctx, tokenHash)
	if err != nil {
		return 0, fmt.Errorf("exchanging token in storage: %w", err)
	}

	return userID, nil
}

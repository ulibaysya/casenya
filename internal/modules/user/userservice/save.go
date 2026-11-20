package userservice

import (
	"context"
	"fmt"

	v1 "github.com/ulibaysya/casenya/api/casenya/hub/user/v1"
)

func (s *UserService) ExchangeUserOpaqueToken(ctx context.Context, token string) (string, error) {
	exchangeResp, err := s.hubClient.ExchangeUserOpaqueToken(ctx, &v1.ExchangeUserOpaqueTokenRequest{
		Token: token,
	})
	if err != nil {
		return "", fmt.Errorf("exchanging token at hub: %w", err)
	}

	session := s.sg.GenerateSession()

	sessionHash := s.hasher.Hash(session)
	if err = s.ustrg.SaveUserSession(ctx, exchangeResp.GetId(), sessionHash); err != nil {
		return "", fmt.Errorf("saving user session: %w", err)
	}

	return session, nil
}

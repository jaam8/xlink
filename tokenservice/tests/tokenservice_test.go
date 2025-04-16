package tests

import (
	"context"
	"fmt"
	"github.com/chempik1234/common-chempik-pkg-golang/pkg/logger"
	"github.com/google/uuid"
	"testing"
	"tokenservice/internal/service"
	"tokenservice/internal/utils"
	"tokenservice/pkg/api/tokenservice"
)

const (
	tokenLength = 32
)

func MakeStringWrong(s string) string {
	return s + "dasodsadsajdsa"
}

type MockTokenRepository struct {
	db          map[string]string
	tokenLength int8
}

func NewTokenRepositoryMock(tokenLength int8) *MockTokenRepository {
	return &MockTokenRepository{db: make(map[string]string), tokenLength: tokenLength}
}

func (r *MockTokenRepository) Check(userId string, token string) (bool, error) {
	return r.db[userId] == token, nil
}

func (r *MockTokenRepository) GetUserIdByToken(token string) (string, error) {
	for key, val := range r.db {
		if val == token {
			return key, nil
		}
	}
	return "", fmt.Errorf("couldn't find token '%s'", token)
}

func (r *MockTokenRepository) Create(userId string) (string, error) {
	newToken := utils.GenerateToken(r.tokenLength)
	r.db[userId] = newToken
	return newToken, nil
}

func (r *MockTokenRepository) Delete(userId string) error {
	if _, ok := r.db[userId]; ok {
		delete(r.db, userId)
	}
	return nil // fmt.Errorf("couldn't find user id '%s'", userId)
}

func TestTokenService(t *testing.T) {
	tokenRepository := NewTokenRepositoryMock(tokenLength)
	tokenService := service.New(tokenRepository)

	loggerCtx, _ := logger.New(context.Background())

	userId := uuid.New().String()

	var ok bool

	var firstCreatedToken string
	var overwrittenToken string

	// create test data
	ok = t.Run("Create test token data", func(t *testing.T) {
		response, err := tokenService.CreateToken(
			loggerCtx,
			&tokenservice.TokenRequest{
				UserId: userId,
			})
		if err != nil {
			t.Fatalf("couldn't create token for user id '%s': %v", userId, err)
		}
		firstCreatedToken = response.Token
	})
	if !ok {
		return
	}

	// check token successfully
	ok = t.Run("check token with correct data", func(t *testing.T) {
		response, err := tokenService.CheckToken(
			loggerCtx,
			&tokenservice.TokenCheckRequest{
				UserId: userId,
				Token:  firstCreatedToken,
			},
		)
		if err != nil {
			t.Fatalf("couldn't check token for user id '%s', even though data is correct: %v", userId, err)
		}
		if !response.Status {
			t.Fatalf("returned 'false' for correct data (userId='%s', token='%s')", userId, firstCreatedToken)
		}
	})

	// check token with wrong token and correct user id
	ok = t.Run("check token with wrong token and correct user id", func(t *testing.T) {
		response, err := tokenService.CheckToken(
			loggerCtx,
			&tokenservice.TokenCheckRequest{
				UserId: userId,
				Token:  MakeStringWrong(firstCreatedToken),
			},
		)
		if err != nil {
			t.Fatalf("couldn't check token for user id '%s': %v", userId, err)
		}
		if response.Status {
			t.Fatalf("returned 'true' for wrong token for user id '%s'", userId)
		}
	})

	// check token with wrong user id and correct token
	ok = t.Run("check token with correct token and wrong user id", func(t *testing.T) {
		response, err := tokenService.CheckToken(
			loggerCtx,
			&tokenservice.TokenCheckRequest{
				UserId: MakeStringWrong(userId),
				Token:  firstCreatedToken,
			})
		if err != nil {
			t.Fatalf("couldn't check token for user id '%s': %v", userId, err)
		}
		if response.Status {
			t.Fatalf("returned 'true' for wrong user id '%s'", userId)
		}
	})

	// recreate token for user id
	ok = t.Run("recreate token for user id", func(t *testing.T) {
		response, err := tokenService.CreateToken(
			loggerCtx,
			&tokenservice.TokenRequest{
				UserId: userId,
			})
		if err != nil {
			t.Fatalf("couldn't recreate token for user id '%s': %v", userId, err)
		}
		overwrittenToken = response.Token
	})

	// check token after rewriting
	ok = t.Run("check token after overwriting", func(t *testing.T) {
		response, err := tokenService.CheckToken(
			loggerCtx,
			&tokenservice.TokenCheckRequest{
				UserId: userId,
				Token:  firstCreatedToken,
			})
		if err != nil {
			t.Fatalf("couldn't check token for user id '%s': %v", userId, err)
		}
		if response.Status {
			t.Fatalf("returned 'true' for deprecated token for user id '%s'", userId)
		}

		response, err = tokenService.CheckToken(
			loggerCtx,
			&tokenservice.TokenCheckRequest{
				UserId: userId,
				Token:  overwrittenToken,
			})
		if err != nil {
			t.Fatalf("couldn't check token for user id '%s': %v", userId, err)
		}
		if !response.Status {
			t.Fatalf("returned 'false' for correct token for user id '%s'", userId)
		}
	})

	// delete token
	ok = t.Run("delete token", func(t *testing.T) {
		response, err := tokenService.DeleteToken(
			loggerCtx,
			&tokenservice.TokenRequest{
				UserId: MakeStringWrong(userId),
			})
		if err != nil {
			t.Fatalf("couldn't delete token for NON-EXISTING user id '%s': %v", userId, err)
		}
		if !response.Status {
			t.Fatalf("returned 'false' for correct delete operation (with NON-EXISTING user id): '%s'", userId)
		}

		response, err = tokenService.DeleteToken(
			loggerCtx,
			&tokenservice.TokenRequest{
				UserId: userId,
			})
		if err != nil {
			t.Fatalf("couldn't delete token for CORRECT user id '%s': %v", userId, err)
		}
		if !response.Status {
			t.Fatalf("returned 'false' for correct delete operation (with CORRECT user id): '%s'", userId)
		}
	})

	// try to read token after deletion
	ok = t.Run("check token after delete", func(t *testing.T) {
		response, err := tokenService.CheckToken(
			loggerCtx,
			&tokenservice.TokenCheckRequest{
				UserId: userId,
			})
		if err != nil {
			t.Fatalf("couldn't check token for user id '%s': %v", userId, err)
		}
		if response.Status {
			t.Fatalf("returned 'true' for deprecated token for user id '%s'", userId)
		}
	})
}

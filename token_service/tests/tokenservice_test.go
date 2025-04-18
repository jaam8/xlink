package tests

import (
	"context"
	"fmt"
	"testing"
	"xlink/common/logger"
	"xlink/token_service/internal/service"
	"xlink/token_service/internal/utils"
	"xlink/token_service/pkg/api/token_service"

	"github.com/google/uuid"
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

func (r *MockTokenRepository) Create(userId string) (string, error) {
	newToken := utils.GenerateToken(r.tokenLength)
	r.db[userId] = newToken
	return newToken, nil
}

func (r *MockTokenRepository) Refresh(userId string, token string) (string, error) {
	status, err := r.Check(userId, token)
	if err != nil {
		if !status {
			return "", fmt.Errorf("Error by checking token, token does not exist:%v", err)
		}
		return "", fmt.Errorf("Error by checking token:%v", err)
	}
	if !status {
		return "", fmt.Errorf("token does not exist")
	}
	return r.Create(userId)
}

func (r *MockTokenRepository) Delete(userId string) (bool, error) {
	if _, ok := r.db[userId]; ok {
		delete(r.db, userId)
		return true, nil
	}
	return false, fmt.Errorf("Error by deleting token")
}

func TestTokenService(t *testing.T) {
	tokenRepository := NewTokenRepositoryMock(tokenLength)
	tokenService := service.New(tokenRepository)

	loggerCtx, err := logger.New(context.Background())
	if err != nil {
		fmt.Errorf("Error while creating logger")
	}

	userId := uuid.New().String()

	var ok bool

	var firstCreatedToken string
	var overwrittenToken string

	// create test data
	ok = t.Run("Create test token data", func(t *testing.T) {
		response, err := tokenService.CreateToken(
			loggerCtx,
			&token_service.TokenRequest{
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
			&token_service.TokenCheckRequest{
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
	if !ok {
		return
	}

	// check token with wrong token and correct user id
	ok = t.Run("check token with wrong token and correct user id", func(t *testing.T) {
		response, err := tokenService.CheckToken(
			loggerCtx,
			&token_service.TokenCheckRequest{
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
	if !ok {
		return
	}

	// check token with wrong user id and correct token
	ok = t.Run("check token with correct token and wrong user id", func(t *testing.T) {
		response, err := tokenService.CheckToken(
			loggerCtx,
			&token_service.TokenCheckRequest{
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
	if !ok {
		return
	}

	//refresh token for user id
	ok = t.Run("refresh token for user id", func(t *testing.T) {
		response, err := tokenService.RefreshToken(
			loggerCtx,
			&token_service.TokenCheckRequest{
				UserId: userId,
				Token:  firstCreatedToken,
			})
		if err != nil {
			t.Fatalf("couldn't refresh token for user id '%s': %v", userId, err)
		}
		overwrittenToken = response.Token
	})
	if !ok {
		return
	}

	// check token after rewriting
	ok = t.Run("check token after overwriting", func(t *testing.T) {
		response, err := tokenService.CheckToken(
			loggerCtx,
			&token_service.TokenCheckRequest{
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
			&token_service.TokenCheckRequest{
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
	if !ok {
		return
	}

	// delete token
	ok = t.Run("delete token", func(t *testing.T) {
		response, err := tokenService.DeleteToken(
			loggerCtx,
			&token_service.TokenRequest{
				UserId: MakeStringWrong(userId),
			})
		if err == nil {
			t.Fatalf("deleted token for NON-EXISTING user id '%s': %v", userId, err)
		}

		response, err = tokenService.DeleteToken(
			loggerCtx,
			&token_service.TokenRequest{
				UserId: userId,
			})
		if err != nil {
			t.Fatalf("couldn't delete token for CORRECT user id '%s': %v", userId, err)
		}
		if !response.Status {
			t.Fatalf("returned 'false' for correct delete operation (with CORRECT user id): '%s'", userId)
		}
	})
	if !ok {
		return
	}

	// try to read token after deletion
	ok = t.Run("check token after delete", func(t *testing.T) {
		_, err := tokenService.CheckToken(
			loggerCtx,
			&token_service.TokenCheckRequest{
				UserId: userId,
			})
		if err == nil {
			t.Fatalf("checked token for user id after deletion'%s': %v", userId, err)
		}
	})
}

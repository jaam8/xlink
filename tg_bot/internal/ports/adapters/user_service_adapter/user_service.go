package user_service_adapter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"time"
	"xlink/common/callers"
	"xlink/common/gen/user_service"
)

type UserServiceAdapter struct {
	Address     string
	BaseURL     string
	DialOptions []grpc.DialOption
	Timeout     time.Duration
}

func NewUserServiceAdapter(address, baseURL string, options []grpc.DialOption, timeout time.Duration) *UserServiceAdapter {
	return &UserServiceAdapter{
		Address:     address,
		BaseURL:     baseURL,
		DialOptions: options,
		Timeout:     timeout,
	}
}

func (us *UserServiceAdapter) CreateUser(tgID *int64) (string, string, error) {
	var reqData struct {
		TgID int64 `json:"tg_id,omitempty"`
	}
	var respData struct {
		UserID string `json:"user_id"`
		Token  string `json:"token"`
	}

	if tgID != nil {
		reqData.TgID = *tgID
	}
	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return "", "", err
	}

	path := us.BaseURL + "create"
	req, err := http.NewRequest(http.MethodPost, path, bytes.NewReader(reqBody))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: us.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	//nolint
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.Status)
		fmt.Println(resp.Body)
		fmt.Println(resp)
		return "", "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", "", err
	}

	return respData.UserID, respData.Token, nil
}

func (us *UserServiceAdapter) GetUserProfile(userID string) (*int64, string, int64, error) {
	var reqData = struct {
		UserID string `json:"user_id"`
	}{
		UserID: userID,
	}
	var respData struct {
		UserID    string `json:"user_id"`
		TgID      *int64 `json:"tg_id"`
		Role      string `json:"role"`
		LinkCount int64  `json:"link_count"`
	}

	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return nil, "", 0, err
	}

	path := us.BaseURL
	req, err := http.NewRequest(http.MethodPost, path, bytes.NewReader(reqBody))
	if err != nil {
		return nil, "", 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: us.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", 0, err
	}
	//nolint
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, "", 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, "", 0, err
	}

	return respData.TgID, respData.Role, respData.LinkCount, nil
}

func (us *UserServiceAdapter) LoginUser(token string) (string, *int64, error) {
	var reqData struct {
		ApiToken string `json:"api_token"`
	}
	var respData struct {
		UserID     string `json:"user_id"`
		TelegramId *int64 `json:"telegram_id"`
	}

	reqData.ApiToken = token
	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return "", nil, err
	}

	path := us.BaseURL + "login"
	req, err := http.NewRequest(http.MethodPost, path, bytes.NewReader(reqBody))
	if err != nil {
		return "", nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: us.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, err
	}
	//nolint
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", nil, err
	}

	return respData.UserID, respData.TelegramId, nil
}

func (us *UserServiceAdapter) SetTgID(userID string, tgID int64) error {
	var reqData struct {
		TgID int64 `json:"tg_id"`
	}
	var respData struct {
		UserID string `json:"user_id"`
	}

	reqData.TgID = tgID
	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return err
	}

	path := us.BaseURL + "update/" + userID
	req, err := http.NewRequest(http.MethodPost, path, bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: us.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	//nolint
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return err
	}

	return nil
}

func (us *UserServiceAdapter) RefreshToken(userID, token string) (string, error) {
	var reqData struct {
		UserID string `json:"user_id"`
		Token  string `json:"token"`
	}
	var respData struct {
		Token string `json:"token"`
	}
	reqData.UserID = userID
	reqData.Token = token
	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return "", err
	}

	path := us.BaseURL + "refresh"
	req, err := http.NewRequest(http.MethodPost, path, bytes.NewReader(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: us.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	//nolint
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", err
	}

	return respData.Token, nil
}

func (us *UserServiceAdapter) GetGRPCClient() (*grpc.ClientConn, *user_service.UserServiceClient, error) {
	grpcConn, err := grpc.NewClient(us.Address, us.DialOptions...)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot connect to user_service by gRPC: %v", err)
	}

	client := user_service.NewUserServiceClient(grpcConn)

	return grpcConn, &client, nil
}

func (us *UserServiceAdapter) GetTokenByTgID(tgID int64) (string, error) {
	conn, clientPointer, err := us.GetGRPCClient()
	if err != nil {
		return "", fmt.Errorf("cannot connect to user_service by gRPC: %v", err)
	}

	defer conn.Close() //nolint

	resultChan := make(chan string, 1)

	err = callers.Timeout(func() error {
		request := &user_service.GetTokenByTgIdRequest{TgId: tgID}
		response, grpcErr := (*clientPointer).GetTokenByTgId(context.Background(), request)
		if grpcErr != nil {
			return fmt.Errorf("error in timeout gRPC caller: %v", grpcErr)
		}
		resultChan <- response.GetToken()
		return nil
	}, us.Timeout)

	if err != nil {
		log.Printf("error in timeout gRPC caller: %v", err)
		return "", fmt.Errorf("couldn't get user_service.GetTokenByTgIdRequest gRPC response: %v", err)
	}

	token := <-resultChan
	close(resultChan)
	return token, nil
}

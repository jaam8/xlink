package shortener_adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ShortenerAdapter struct {
	BaseURL string
	Timeout time.Duration
}

func NewShortenerAdapter(baseURL string, timeout time.Duration) *ShortenerAdapter {
	return &ShortenerAdapter{
		BaseURL: baseURL,
		Timeout: timeout,
	}
}

func (s *ShortenerAdapter) CreateLink(userToken, targetUrl string, shortLink *string) (string, string, string, string, error) {
	var reqData struct {
		ShortLink string `json:"short_link,omitempty"`
		TargetUrl string `json:"target_url"`
	}
	var respData struct {
		LinkID    string `json:"link_id"`
		UserID    string `json:"user_id"`
		ShortLink string `json:"short_link"`
		TargetUrl string `json:"target_url"`
		CreatedAt string `json:"created_at"`
		ExpireAt  string `json:"expire_at"`
	}

	if shortLink != nil {
		reqData.ShortLink = *shortLink
	}
	reqData.TargetUrl = targetUrl
	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return "", "", "", "", err
	}

	path := s.BaseURL + "create"
	req, err := http.NewRequest(http.MethodPost, path, bytes.NewReader(reqBody))
	if err != nil {
		return "", "", "", "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", userToken)
	client := &http.Client{Timeout: s.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", "", err
	}
	// nolint
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		log.Println(req)
		log.Println(resp)
		return "", "", "", "", fmt.Errorf("shortener create failed, status code: %d", resp.StatusCode)
	}
	if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", "", "", "", err
	}
	log.Println(respData)
	createdAtDate, err := time.Parse(time.RFC3339, respData.CreatedAt)
	if err != nil {
		return "", "", "", "", err
	}
	createdAt := createdAtDate.Format("2006-01-02 15:04")
	expireAtDate, err := time.Parse(time.RFC3339, respData.ExpireAt)
	if err != nil {
		return "", "", "", "", err
	}
	expireAt := expireAtDate.Format("2006-01-02 15:04")

	return respData.ShortLink, respData.TargetUrl, createdAt, expireAt, nil
}

func (s *ShortenerAdapter) UpdateLink(userToken, shortLink, targetURL string,
	expireAt time.Time, regenerate bool) (string, string, string, string, error) {
	var reqData struct {
		ShortLink  string `json:"short_link"`
		TargetURL  string `json:"target_url"`
		ExpireAt   string `json:"expire_at"`
		Regenerate bool   `json:"regenerate"`
	}
	var respData struct {
		ShortLink string `json:"short_link"`
		TargetUrl string `json:"target_url"`
		CreatedAt string `json:"created_at"`
		ExpireAt  string `json:"expire_at"`
	}

	reqData.ShortLink = shortLink
	reqData.TargetURL = targetURL
	reqData.ExpireAt = expireAt.Format(time.RFC3339)
	reqData.Regenerate = regenerate
	reqBody, err := json.Marshal(reqData)
	if err != nil {
		return "", "", "", "", err
	}

	path := s.BaseURL + "update/" + shortLink
	req, err := http.NewRequest(http.MethodPut, path, bytes.NewReader(reqBody))
	if err != nil {
		return "", "", "", "", err
	}
	req.Header.Set("Authorization", userToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: s.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", "", err
	}
	// nolint
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", "", "", "", err
	}
	if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", "", "", "", err
	}

	createdAtDate, err := time.Parse(time.RFC3339, respData.CreatedAt)
	if err != nil {
		return "", "", "", "", err
	}
	createdAt := createdAtDate.Format("2006-01-02 15:04")
	expireAtDate, err := time.Parse(time.RFC3339, respData.ExpireAt)
	if err != nil {
		return "", "", "", "", err
	}
	expireAtResp := expireAtDate.Format("2006-01-02 15:04")

	return respData.ShortLink, respData.TargetUrl, createdAt, expireAtResp, nil
}

func (s *ShortenerAdapter) GetUserLinks(userToken string) ([]string, error) {
	var respData struct {
		Links []string `json:"shortLinks"`
	}
	path := s.BaseURL + "my-links"
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", userToken)

	client := &http.Client{Timeout: s.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, err
	}
	return respData.Links, nil
}

func (s *ShortenerAdapter) DeleteLink(userToken, shortLink string) error {
	path := s.BaseURL + "delete/" + shortLink
	req, err := http.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", userToken)

	client := &http.Client{Timeout: s.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return err
	}
	return nil
}

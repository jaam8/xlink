package analytics_adapter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"xlink/tg_bot/internal/models"
)

type AnalyticsAdapter struct {
	BaseURL string
	Timeout time.Duration
}

func NewAnalyticsAdapter(baseURL string, timeout time.Duration) *AnalyticsAdapter {
	return &AnalyticsAdapter{
		BaseURL: baseURL,
		Timeout: timeout,
	}
}

type Request struct {
	ShortLink string `json:"short_link"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func (a *AnalyticsAdapter) ClicksByCountry(userToken, shortLink, startDate, endDate string) ([]models.Stat[models.CountryDayStat], error) {
	var respData struct {
		Data []models.Stat[models.CountryDayStat] `json:"data"`
	}
	var reqData Request
	reqData.ShortLink = shortLink
	reqData.StartDate = startDate
	reqData.EndDate = endDate

	path := a.BaseURL + "by_country"
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("short_link", shortLink)
	q.Add("start_date", startDate)
	q.Add("end_date", endDate)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Authorization", userToken)

	client := &http.Client{Timeout: a.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get clicks by country: %s", resp.Status)
	}
	if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, err
	}
	if len(respData.Data) == 0 {
		return nil, fmt.Errorf("no data found")
	}

	return respData.Data, nil
}

func (a *AnalyticsAdapter) ClicksByRegion(userToken, shortLink, startDate, endDate string) ([]models.Stat[models.RegionDayStat], error) {
	var respData struct {
		Data []models.Stat[models.RegionDayStat] `json:"data"`
	}
	var reqData Request
	reqData.ShortLink = shortLink
	reqData.StartDate = startDate
	reqData.EndDate = endDate
	path := a.BaseURL + "by_region"
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("short_link", shortLink)
	q.Add("start_date", startDate)
	q.Add("end_date", endDate)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Authorization", userToken)

	client := &http.Client{Timeout: a.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get clicks by region: %s", resp.Status)
	}
	if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, err
	}
	if len(respData.Data) == 0 {
		return nil, fmt.Errorf("no data found")
	}
	return respData.Data, nil
}

func (a *AnalyticsAdapter) ClicksByBrowser(userToken, shortLink, startDate, endDate string) ([]models.Stat[models.BrowserDayStat], error) {
	var respData struct {
		Data []models.Stat[models.BrowserDayStat] `json:"data"`
	}
	var reqData Request
	reqData.ShortLink = shortLink
	reqData.StartDate = startDate
	reqData.EndDate = endDate
	path := a.BaseURL + "by_browser"
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("short_link", shortLink)
	q.Add("start_date", startDate)
	q.Add("end_date", endDate)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Authorization", userToken)

	client := &http.Client{Timeout: a.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get clicks by browser: %s", resp.Status)
	}
	if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, err
	}
	if len(respData.Data) == 0 {
		return nil, fmt.Errorf("no data found")
	}
	return respData.Data, nil
}

func (a *AnalyticsAdapter) ClicksByDevice(userToken, shortLink, startDate, endDate string) ([]models.Stat[models.DeviceDayStat], error) {
	var respData struct {
		Data []models.Stat[models.DeviceDayStat] `json:"data"`
	}
	var reqData Request
	reqData.ShortLink = shortLink
	reqData.StartDate = startDate
	reqData.EndDate = endDate
	path := a.BaseURL + "by-device-type"
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("short_link", shortLink)
	q.Add("start_date", startDate)
	q.Add("end_date", endDate)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Authorization", userToken)

	client := &http.Client{Timeout: a.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get clicks by device_type: %s", resp.Status)
	}
	if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, err
	}
	if len(respData.Data) == 0 {
		return nil, fmt.Errorf("no data found")
	}
	return respData.Data, nil
}

func (a *AnalyticsAdapter) ClicksByOS(userToken, shortLink, startDate, endDate string) ([]models.Stat[models.OSDayStat], error) {
	var respData struct {
		Data []models.Stat[models.OSDayStat] `json:"data"`
	}
	var reqData Request
	reqData.ShortLink = shortLink
	reqData.StartDate = startDate
	reqData.EndDate = endDate
	path := a.BaseURL + "by-os"
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("short_link", shortLink)
	q.Add("start_date", startDate)
	q.Add("end_date", endDate)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Authorization", userToken)

	client := &http.Client{Timeout: a.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get clicks by os: %s", resp.Status)
	}
	if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, err
	}
	if len(respData.Data) == 0 {
		return nil, fmt.Errorf("no data found")
	}
	return respData.Data, nil
}

func (a *AnalyticsAdapter) ClicksByReferrer(userToken, shortLink, startDate, endDate string) ([]models.Stat[models.ReferrerDayStat], error) {
	var respData struct {
		Data []models.Stat[models.ReferrerDayStat] `json:"data"`
	}
	var reqData Request
	reqData.ShortLink = shortLink
	reqData.StartDate = startDate
	reqData.EndDate = endDate
	path := a.BaseURL + "by-referrer"
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("short_link", shortLink)
	q.Add("start_date", startDate)
	q.Add("end_date", endDate)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Authorization", userToken)

	client := &http.Client{Timeout: a.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get clicks by referrer: %s", resp.Status)
	}
	if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, err
	}
	if len(respData.Data) == 0 {
		return nil, fmt.Errorf("no data found")
	}
	return respData.Data, nil
}

func (a *AnalyticsAdapter) ClicksByHour(userToken, shortLink, startDate, endDate string) ([]models.Stat[models.HourDayStat], error) {
	var respData struct {
		Data []models.Stat[models.HourDayStat] `json:"data"`
	}
	var reqData Request
	reqData.ShortLink = shortLink
	reqData.StartDate = startDate
	reqData.EndDate = endDate
	path := a.BaseURL + "by-hour"
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("short_link", shortLink)
	q.Add("start_date", startDate)
	q.Add("end_date", endDate)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Authorization", userToken)

	client := &http.Client{Timeout: a.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get clicks by hour: %s", resp.Status)
	}
	if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, err
	}
	if len(respData.Data) == 0 {
		return nil, fmt.Errorf("no data found")
	}
	return respData.Data, nil
}

func (a *AnalyticsAdapter) ClicksByDate(userToken, shortLink, startDate, endDate string) ([]models.DateStat, error) {
	var respData struct {
		Data []models.DateStat `json:"data"`
	}
	var reqData Request
	reqData.ShortLink = shortLink
	reqData.StartDate = startDate
	reqData.EndDate = endDate
	path := a.BaseURL + "by-date"
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("short_link", shortLink)
	q.Add("start_date", startDate)
	q.Add("end_date", endDate)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Authorization", userToken)

	client := &http.Client{Timeout: a.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get clicks by date: %s", resp.Status)
	}
	if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, err
	}
	if len(respData.Data) == 0 {
		return nil, fmt.Errorf("no data found")
	}
	return respData.Data, nil
}

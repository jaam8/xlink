package helper

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func ParseRegionAndCountry(ip string) (string, string, error) {
	url := fmt.Sprintf("https://ipinfo.io/%s/json", ip)

	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	//nolint
	defer resp.Body.Close()

	var info struct {
		Country string `json:"country"`
		Region  string `json:"region"`
		bogon   bool   `json:"bogon,omitempty"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return "", "", err
	}

	if info.bogon || info.Country == "" || info.Region == "" {
		return "Unknown", "Unknown", nil
	}
	return info.Region, info.Country, nil
}

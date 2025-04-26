package renderer_service_adapters

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type RendererServiceRepositoryHTTP struct {
	httpClient *http.Client
	address    string
}

func NewRendererServiceRepositoryHTTP(protocol string, hostName string, port string, timeout time.Duration) *RendererServiceRepositoryHTTP {
	return &RendererServiceRepositoryHTTP{
		address: fmt.Sprintf("%s://%s:%s", protocol, hostName, port),
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (r *RendererServiceRepositoryHTTP) Generate(shortLink string, param string, startDate time.Time, endDate time.Time, linkOwner string) ([]byte, error) {
	var err error

	var request *http.Request
	request, err = http.NewRequest("GET",
		fmt.Sprintf("%s/image?short_link=%s&param=%s&start_date=%s&end_date=%s&link_owner=%s",
			r.address, shortLink, param, startDate.Format(time.DateOnly), endDate.Format(time.DateOnly), linkOwner),
		nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP renderer request: %s", err)
	}

	var fileResponseBytes *http.Response
	fileResponseBytes, err = r.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error executing HTTP renderer request: %s", err)
	}

	defer fileResponseBytes.Body.Close() //nolint:all

	if fileResponseBytes.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error executing HTTP renderer request, status code: %d", fileResponseBytes.StatusCode)
	}

	var responseBytes []byte
	responseBytes, err = io.ReadAll(fileResponseBytes.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading HTTP renderer response body: %s", err)
	}

	return responseBytes, nil
}

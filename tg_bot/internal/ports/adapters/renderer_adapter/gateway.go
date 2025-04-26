package renderer_adapter

import (
	"fmt"
	"time"
)

type RendererAdapterGateway struct {
	BaseUrl string
}

func NewRendererAdapterGateway(baseUrl string) *RendererAdapterGateway {
	return &RendererAdapterGateway{BaseUrl: baseUrl}
}

func (s *RendererAdapterGateway) RenderChart(chartType string) (string, error) {
	panic("use GetImageUrl()")
}

func (s *RendererAdapterGateway) GetImageUrl(shortLink, token, param string, startDate time.Time, endDate time.Time) string {
	return fmt.Sprintf("%s/%s/?token=%s&param=%s&start_date=%s&end_date=%s",
		s.BaseUrl,
		shortLink,
		token, param,
		startDate.Format(time.DateOnly),
		endDate.Format(time.DateOnly),
	)
}

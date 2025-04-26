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
	url := fmt.Sprintf("%s%s/?param=%s&start_date=%s&end_date=%s&token=%s",
		s.BaseUrl,
		shortLink,
		param,
		startDate.Format(time.DateOnly),
		endDate.Format(time.DateOnly),
		token,
	)
	return url
}

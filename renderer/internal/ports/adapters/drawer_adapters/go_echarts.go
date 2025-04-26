package drawer_adapters

import (
	"bytes"
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"xlink/renderer/internal/statistics_data"
)

type DrawerRepositoryEcharts struct {
}

func NewDrawerRepositoryEcharts() *DrawerRepositoryEcharts {
	return &DrawerRepositoryEcharts{}
}

func generateBarChart(input statistics_data.StatisticsData, metric string, param string) *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: fmt.Sprintf("%s by %s", metric, param)}),
		charts.WithLegendOpts(opts.Legend{Show: opts.Bool(true)}),
		charts.WithTooltipOpts(opts.Tooltip{Show: opts.Bool(true)}),
	)

	valuesMap := make(map[string][]uint64)
	var dates []string
	for idx, item := range input.Stats {
		dates = append(dates, item.Date)
		tmp := map[string]uint64{}
		for _, s := range item.Items {
			var val uint64
			if metric == "Clicks" {
				val = s.Clicks
			} else {
				val = s.UniqueClicks
			}
			tmp[s.ParamValue] = val
		}
		for country := range tmp {
			if _, ok := valuesMap[country]; !ok {
				valuesMap[country] = make([]uint64, len(input.Stats))
			}
		}
		for country, v := range tmp {
			valuesMap[country][idx] = v
		}
	}

	bar.SetXAxis(dates)
	for country, values := range valuesMap {
		items := make([]opts.BarData, len(values))
		for i, v := range values {
			items[i] = opts.BarData{Value: v}
		}
		bar.AddSeries(country, items)
	}
	return bar
}

func takeScreenshot(htmlBytes []byte) ([]byte, error) {
	var screenshot []byte

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			// raw html
			chromedp.Navigate("data:text/html;charset=utf-8," + string(htmlBytes))
			return nil
		}),
		chromedp.Screenshot("body", &screenshot, chromedp.NodeVisible),
	)

	if err != nil {
		return nil, fmt.Errorf("couldn't take screenshot: %w", err)
	}

	return screenshot, nil
}

func (d DrawerRepositoryEcharts) Generate(input statistics_data.StatisticsData, paramName string) ([]byte, error) {
	page := components.NewPage()
	page.AddCharts(
		generateBarChart(input, "Clicks", paramName),
		generateBarChart(input, "UniqueClicks", paramName),
	)

	var fileBuf bytes.Buffer

	err := page.Render(&fileBuf)
	if err != nil {
		return nil, fmt.Errorf("couldn't render image: %w", err)
	}

	output := fileBuf.Bytes()

	var screenshot []byte
	screenshot, err = takeScreenshot(output)

	if err != nil {
		return nil, fmt.Errorf("couldn't take screenshot: %w", err)
	}

	return screenshot, nil
}

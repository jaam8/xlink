package storage

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
	"xlink/common/logger"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"xlink/analytics/internal/models"
)

type ClickHouseAdapter struct {
	Conn driver.Conn
}

func NewClickHouseAdapter(conn driver.Conn) *ClickHouseAdapter {
	return &ClickHouseAdapter{
		Conn: conn,
	}
}

func (c *ClickHouseAdapter) SaveClicks(ctx context.Context, clicks []*models.Click) error {
	batch, err := c.Conn.PrepareBatch(ctx, "INSERT INTO xlink.clicks")
	if err != nil {
		return err
	}

	for _, click := range clicks {
		logger.GetLoggerFromCtx(ctx).Info(ctx,
			"click", zap.Any("click", click))
		if click == nil {
			continue
		}

		if err = batch.AppendStruct(click); err != nil {
			return err
		}
	}

	return batch.Send()
}

func (c *ClickHouseAdapter) GetClicksByCountry(
	startDate, endDate *time.Time, linkOwner uuid.UUID, shortLink string,
) ([]models.CountryDayStats, error) {
	query := `
		SELECT
			toDate(clicked_at) as date,
			country,
			count() as clicks,
			sum(is_unique) as unique_clicks
		FROM xlink.clicks
		WHERE
			clicked_at BETWEEN ? AND ?
			AND link_owner = ?
			AND short_link = ?
		GROUP BY date, country
		ORDER BY date;
		`
	start := startDate.UTC().Format("2006-01-02 15:04:05")
	end := endDate.UTC().Format("2006-01-02 15:04:05")
	rows, err := c.Conn.Query(context.Background(), query, start, end, linkOwner, shortLink)
	if err != nil {
		return nil, fmt.Errorf("query clicks by country: %w", err)
	}
	statsMap := make(map[string][]models.CountryClicks)
	var dateOrder []string
	// nolint
	defer rows.Close()
	for rows.Next() {
		var (
			date         time.Time
			country      string
			clicks       uint64
			uniqueClicks uint64
		)
		if err = rows.Scan(&date, &country, &clicks, &uniqueClicks); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		day := date.Format("2006-01-02")
		if _, ok := statsMap[day]; !ok {
			dateOrder = append(dateOrder, day)
		}
		statsMap[day] = append(statsMap[day], models.CountryClicks{
			Country:      country,
			Clicks:       clicks,
			UniqueClicks: uniqueClicks,
		})
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %v", err)
	}

	result := make([]models.CountryDayStats, 0, len(dateOrder))
	for _, day := range dateOrder {
		result = append(result, models.CountryDayStats{
			Date:  day,
			Stats: statsMap[day],
		})
	}

	return result, nil
}

func (c *ClickHouseAdapter) GetClicksByRegion(
	startDate, endDate *time.Time, linkOwner uuid.UUID, shortLink string,
) ([]models.RegionDayStats, error) {
	query := `
		SELECT
			toDate(clicked_at) as date,
			region,
			count() as clicks,
			sum(is_unique) as unique_clicks
		FROM xlink.clicks
		WHERE
			clicked_at BETWEEN ? AND ?
			AND link_owner = ?
			AND short_link = ?
		GROUP BY date, region
		ORDER BY date;
		`
	start := startDate.UTC().Format("2006-01-02 15:04:05")
	end := endDate.UTC().Format("2006-01-02 15:04:05")
	rows, err := c.Conn.Query(context.Background(), query, start, end, linkOwner, shortLink)
	if err != nil {
		return nil, fmt.Errorf("query clicks by region: %w", err)
	}
	statsMap := make(map[string][]models.RegionClicks)
	var dateOrder []string
	// nolint
	defer rows.Close()
	for rows.Next() {
		var (
			date         time.Time
			region       string
			clicks       uint64
			uniqueClicks uint64
		)
		if err = rows.Scan(&date, &region, &clicks, &uniqueClicks); err != nil {
			return nil, fmt.Errorf("scan row: %v", err)
		}
		day := date.Format("2006-01-02")
		if _, ok := statsMap[day]; !ok {
			dateOrder = append(dateOrder, day)
		}
		statsMap[day] = append(statsMap[day], models.RegionClicks{
			Region:       region,
			Clicks:       clicks,
			UniqueClicks: uniqueClicks,
		})
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}

	result := make([]models.RegionDayStats, 0, len(dateOrder))
	for _, day := range dateOrder {
		result = append(result, models.RegionDayStats{
			Date:  day,
			Stats: statsMap[day],
		})
	}

	return result, nil
}

func (c *ClickHouseAdapter) GetClicksByReferrer(
	startDate, endDate *time.Time, linkOwner uuid.UUID, shortLink string,
) ([]models.ReferrerDayStats, error) {
	query := `
		SELECT
			toDate(clicked_at) as date,
			referrer,
			count() as clicks,
			sum(is_unique) as unique_clicks
		FROM xlink.clicks
		WHERE
			clicked_at BETWEEN ? AND ?
			AND link_owner = ?
			AND short_link = ?
		GROUP BY date, referrer
		ORDER BY date;
	`
	start := startDate.UTC().Format("2006-01-02 15:04:05")
	end := endDate.UTC().Format("2006-01-02 15:04:05")
	rows, err := c.Conn.Query(context.Background(), query, start, end, linkOwner, shortLink)
	if err != nil {
		return nil, fmt.Errorf("query clicks by referer: %w", err)
	}

	statsMap := make(map[string][]models.ReferrerClicks)
	var dateOrder []string
	// nolint
	defer rows.Close()
	for rows.Next() {
		var (
			date         time.Time
			referrer     string
			clicks       uint64
			uniqueClicks uint64
		)
		if err = rows.Scan(&date, &referrer, &clicks, &uniqueClicks); err != nil {
			return nil, fmt.Errorf("scan row: %v", err)
		}
		day := date.Format("2006-01-02")
		if _, ok := statsMap[day]; !ok {
			dateOrder = append(dateOrder, day)
		}
		statsMap[day] = append(statsMap[day], models.ReferrerClicks{
			Referrer:     referrer,
			Clicks:       clicks,
			UniqueClicks: uniqueClicks,
		})
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}

	result := make([]models.ReferrerDayStats, 0, len(dateOrder))
	for _, day := range dateOrder {
		result = append(result, models.ReferrerDayStats{
			Date:  day,
			Stats: statsMap[day],
		})
	}
	return result, nil
}

func (c *ClickHouseAdapter) GetClicksByBrowser(
	startDate, endDate *time.Time, linkOwner uuid.UUID, shortLink string,
) ([]models.BrowserDayStats, error) {
	query := `
		SELECT
			toDate(clicked_at) as date,
			browser,
			count() as clicks,
			sum(is_unique) as unique_clicks
		FROM xlink.clicks
		WHERE
			clicked_at BETWEEN ? AND ?
			AND link_owner = ?
			AND short_link = ?
		GROUP BY date, browser
		ORDER BY date;
		`
	start := startDate.UTC().Format("2006-01-02 15:04:05")
	end := endDate.UTC().Format("2006-01-02 15:04:05")
	rows, err := c.Conn.Query(context.Background(), query, start, end, linkOwner, shortLink)
	if err != nil {
		return nil, fmt.Errorf("query clicks by browser: %w", err)
	}
	statsMap := make(map[string][]models.BrowserClicks)
	var dateOrder []string
	// nolint
	defer rows.Close()
	for rows.Next() {
		var (
			date         time.Time
			browser      string
			clicks       uint64
			uniqueClicks uint64
		)
		if err = rows.Scan(&date, &browser, &clicks, &uniqueClicks); err != nil {
			return nil, fmt.Errorf("scan row: %v", err)
		}
		day := date.Format("2006-01-02")
		if _, ok := statsMap[day]; !ok {
			dateOrder = append(dateOrder, day)
		}
		statsMap[day] = append(statsMap[day], models.BrowserClicks{
			Browser:      browser,
			Clicks:       clicks,
			UniqueClicks: uniqueClicks,
		})
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}

	result := make([]models.BrowserDayStats, 0, len(dateOrder))
	for _, day := range dateOrder {
		result = append(result, models.BrowserDayStats{
			Date:  day,
			Stats: statsMap[day],
		})
	}

	return result, nil
}

func (c *ClickHouseAdapter) GetClicksByOS(
	startDate, endDate *time.Time, linkOwner uuid.UUID, shortLink string,
) ([]models.OSDayStats, error) {
	query := `
		SELECT
			toDate(clicked_at) as date,
			os,
			count() as clicks,
			sum(is_unique) as unique_clicks
		FROM xlink.clicks
		WHERE
			clicked_at BETWEEN ? AND ?
			AND link_owner = ?
			AND short_link = ?
		GROUP BY date, os
		ORDER BY date;
		`
	start := startDate.UTC().Format("2006-01-02 15:04:05")
	end := endDate.UTC().Format("2006-01-02 15:04:05")
	rows, err := c.Conn.Query(context.Background(), query, start, end, linkOwner, shortLink)
	if err != nil {
		return nil, fmt.Errorf("query clicks by OS: %w", err)
	}
	statsMap := make(map[string][]models.OSClicks)
	var dateOrder []string
	// nolint
	defer rows.Close()
	for rows.Next() {
		var (
			date         time.Time
			os           string
			clicks       uint64
			uniqueClicks uint64
		)
		if err = rows.Scan(&date, &os, &clicks, &uniqueClicks); err != nil {
			return nil, fmt.Errorf("scan row: %v", err)
		}
		day := date.Format("2006-01-02")
		if _, ok := statsMap[day]; !ok {
			dateOrder = append(dateOrder, day)
		}
		statsMap[day] = append(statsMap[day], models.OSClicks{
			OS:           os,
			Clicks:       clicks,
			UniqueClicks: uniqueClicks,
		})
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}

	result := make([]models.OSDayStats, 0, len(dateOrder))
	for _, day := range dateOrder {
		result = append(result, models.OSDayStats{
			Date:  day,
			Stats: statsMap[day],
		})
	}

	return result, nil
}

func (c *ClickHouseAdapter) GetClicksByDeviceType(
	startDate, endDate *time.Time, linkOwner uuid.UUID, shortLink string,
) ([]models.DeviceDayStats, error) {
	query := `
		SELECT
			toDate(clicked_at) as date,
			device_type,
			count() as clicks,
			sum(is_unique) as unique_clicks
		FROM xlink.clicks
		WHERE
			clicked_at BETWEEN ? AND ?
			AND link_owner = ?
			AND short_link = ?
		GROUP BY date, device_type
		ORDER BY date;
		`
	start := startDate.UTC().Format("2006-01-02 15:04:05")
	end := endDate.UTC().Format("2006-01-02 15:04:05")
	rows, err := c.Conn.Query(context.Background(), query, start, end, linkOwner, shortLink)
	if err != nil {
		return nil, fmt.Errorf("query clicks by device_type: %w", err)
	}
	statsMap := make(map[string][]models.DeviceClicks)
	var dateOrder []string
	// nolint
	defer rows.Close()
	for rows.Next() {
		var (
			date         time.Time
			deviceType   string
			clicks       uint64
			uniqueClicks uint64
		)
		if err = rows.Scan(&date, &deviceType, &clicks, &uniqueClicks); err != nil {
			return nil, fmt.Errorf("scan row: %v", err)
		}
		day := date.Format("2006-01-02")
		if _, ok := statsMap[day]; !ok {
			dateOrder = append(dateOrder, day)
		}
		statsMap[day] = append(statsMap[day], models.DeviceClicks{
			DeviceType:   deviceType,
			Clicks:       clicks,
			UniqueClicks: uniqueClicks,
		})
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}

	result := make([]models.DeviceDayStats, 0, len(dateOrder))
	for _, day := range dateOrder {
		result = append(result, models.DeviceDayStats{
			Date:  day,
			Stats: statsMap[day],
		})
	}

	return result, nil
}

func (c *ClickHouseAdapter) GetClicksByHour(
	startDate, endDate *time.Time, linkOwner uuid.UUID, shortLink string,
) ([]models.HourDayStats, error) {
	query := `
		SELECT
			toDate(clicked_at) as date,
			toHour(clicked_at) as hour,
			count() as clicks,
			sum(is_unique) as unique_clicks
		FROM xlink.clicks
		WHERE
			clicked_at BETWEEN ? AND ?
			AND link_owner = ?
			AND short_link = ?
		GROUP BY date, hour
		ORDER BY date;
		`
	start := startDate.UTC().Format("2006-01-02 15:04:05")
	end := endDate.UTC().Format("2006-01-02 15:04:05")
	rows, err := c.Conn.Query(context.Background(), query, start, end, linkOwner, shortLink)
	if err != nil {
		return nil, fmt.Errorf("query clicks by hour: %w", err)
	}
	statsMap := make(map[string][]models.HourStat)
	var dateOrder []string
	// nolint
	defer rows.Close()
	for rows.Next() {
		var (
			date         time.Time
			hour         uint8
			clicks       uint64
			uniqueClicks uint64
		)
		if err = rows.Scan(&date, &hour, &clicks, &uniqueClicks); err != nil {
			return nil, fmt.Errorf("scan row: %v", err)
		}
		day := date.Format("2006-01-02")
		if _, ok := statsMap[day]; !ok {
			dateOrder = append(dateOrder, day)
		}
		dayFormat, _ := time.Parse("2006-01-02 15:04:05", day)
		statsMap[day] = append(statsMap[day], models.HourStat{
			Hour:         uint32(dayFormat.Hour()),
			Clicks:       clicks,
			UniqueClicks: uniqueClicks,
		})
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}

	result := make([]models.HourDayStats, 0, len(dateOrder))
	for _, day := range dateOrder {
		result = append(result, models.HourDayStats{
			Date:  day,
			Stats: statsMap[day],
		})
	}

	return result, nil
}

func (c *ClickHouseAdapter) GetClicksByDate(
	startDate, endDate *time.Time, linkOwner uuid.UUID, shortLink string,
) ([]models.DateStat, error) {
	query := `
		SELECT
			toDate(clicked_at) as date,
			count() as clicks,
			sum(is_unique) as unique_clicks
		FROM xlink.clicks
		WHERE
			clicked_at BETWEEN ? AND ?
			AND link_owner = ?
			AND short_link = ?
		GROUP BY date
		ORDER BY date;
		`
	start := startDate.UTC().Format("2006-01-02 15:04:05")
	end := endDate.UTC().Format("2006-01-02 15:04:05")
	rows, err := c.Conn.Query(context.Background(), query, start, end, linkOwner, shortLink)
	if err != nil {
		return nil, fmt.Errorf("query clicks by date: %w", err)
	}
	var result []models.DateStat
	// nolint
	defer rows.Close()
	for rows.Next() {
		var (
			date         time.Time
			clicks       uint64
			uniqueClicks uint64
		)
		if err = rows.Scan(&date, &clicks, &uniqueClicks); err != nil {
			return nil, fmt.Errorf("scan row: %v", err)
		}
		day := date.Format("2006-01-02")
		result = append(result, models.DateStat{
			Date:         day,
			Clicks:       clicks,
			UniqueClicks: uniqueClicks,
		})
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}

	return result, nil
}

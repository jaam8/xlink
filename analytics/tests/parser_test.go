package tests

import (
	"testing"
	"xlink/analytics/internal/service/helper"
)

func TestParseRegionAndCountry(t *testing.T) {

	region, country, err := helper.ParseRegionAndCountry("1.1.1.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if region != "Queensland" || country != "AU" {
		t.Errorf("unexpected result: got region=%s, country=%s", region, country)
	}
}

package tests

import (
	"testing"
	"time"
	"xlink/common/gen/shortener"
	"xlink/shortener/internal/service/helper"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGetValidatedId(t *testing.T) {
	testRequest := shortener.DeleteLinkRequest{LinkId: "f9e71cb4-e1e1-4721-8eef-806338db2222"}

	testResponse, err := helper.GetValidatedId(&testRequest)
	if err != nil {
		t.Fatalf("Unable to get id, %v", err)
	}

	expected, err := uuid.Parse("f9e71cb4-e1e1-4721-8eef-806338db2222")
	if err != nil {
		t.Fatalf("Unable to parse id, %v", err)
	}

	assert.Equal(t, expected, testResponse)

}

func TestGetValidatedUserId(t *testing.T) {

	userIDStr := "f9e71cb4-e1e1-4721-8eef-806338db7282"
	shortLinkStr := "http://qwerty"

	testCreateLinkRequest := shortener.CreateLinkRequest{
		UserId:    userIDStr,
		ShortLink: &shortLinkStr,
		TargetUrl: "http://qwertysdijvnisdnc",
	}

	testResponse, err := helper.GetValidatedUserId(&testCreateLinkRequest)
	if err != nil {
		t.Fatalf("Unable to get id, %v", err)
	}

	expected, err := uuid.Parse(userIDStr)
	if err != nil {
		t.Fatalf("Unable to parse id, %v", err)
	}

	assert.Equal(t, expected, testResponse)

}

func TestGetValidatedExpireAt(t *testing.T) {

	idStr := "f9e71cb4-e1e1-4721-8eef-806338db2222"
	userIDStr := "f9e71cb4-e1e1-4721-8eef-806338db7282"
	shortLinkStr := "http://qwerty"
	expireAtStr := "2025-04-16T11:28:07+03:00"
	expireAtStr1 := "2025-04-16T11:28:08+03:00"

	expireAtTime, err := time.Parse(time.RFC3339, expireAtStr)
	require.NoError(t, err)
	expireAtTime1, err := time.Parse(time.RFC3339, expireAtStr1)
	require.NoError(t, err)

	expireAtTimestamppb := timestamppb.New(expireAtTime)
	require.NoError(t, err)

	regenerated := true

	testRequest := shortener.UpdateLinkRequest{
		LinkId:     idStr,
		UserId:     userIDStr,
		Regenerate: regenerated,
		ShortLink:  &shortLinkStr,
		TargetUrl:  "http://qwertysdijvnisdnc",
		ExpireAt:   expireAtTimestamppb,
	}

	testResponse, err := helper.GetValidatedExpireAt(&testRequest, expireAtTime1)
	if err != nil {
		t.Fatalf("Unable to get ExpireAt, %v", err)
	}

	expected := expireAtStr

	assert.Equal(t, expected, testResponse)

}

func TestValidatedStringNotEmpty_with_no_empty_string(t *testing.T) {

	str := "srgfuiybesuifh"

	resp := helper.ValidateStringNotEmpty(str)

	assert.NoError(t, resp)

}
func TestValidatedStringNotEmpty_with_empty_string(t *testing.T) {

	str := ""

	resp := helper.ValidateStringNotEmpty(str)

	assert.Error(t, resp)

}

func TestValidatedUrl_with_no_empty_url(t *testing.T) {

	url := "https://shbvgfuyesbfuybe"

	resp := helper.ValidateUrl(url)

	assert.NoError(t, resp)
}

func TestValidatedUrl_with_empty_url(t *testing.T) {

	url := ""

	resp := helper.ValidateUrl(url)

	assert.Error(t, resp)
}

func TestValidateShortLink_no_empty_string(t *testing.T) {
	str := "nsurgfhbsuyfb"

	resp, err := helper.ValidateShortLink(str)
	if resp == "" {
		t.Error("uxexpected empty string")
	}

	assert.NoError(t, err)
}

func TestValidateShortLink_with_empty_string(t *testing.T) {
	str := ""

	_, err := helper.ValidateShortLink(str)

	assert.Error(t, err)
}

func TestValidateIPAddress_correct(t *testing.T) {
	ip := "1.1.1.11.1."

	resp, err := helper.ValidateIPAddress(ip)
	if resp == "" {
		t.Error("uxexpected empty ip")
	}

	assert.NoError(t, err)

}

func TestValidateNotEmptyStr_with_no_empty_string(t *testing.T) {
	str := "srgfuiybesuifh"

	resp, err := helper.ValidateNotEmptyStr(str)
	if resp == "" {
		t.Error("uxexpected empty string")
	}

	assert.NoError(t, err)
}

func TestValidateNotEmptyStr_with_empty_string(t *testing.T) {
	str := ""

	_, err := helper.ValidateNotEmptyStr(str)

	assert.NoError(t, err)
}

package tests

import (
	"testing"
	"time"
	"xlink/shortener/internal/service/helper"
	"xlink/shortener/pkg/api/shortener"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetValidatedId(t *testing.T) {
	testRequest := shortener.DeleteLinkRequest{Id: "f9e71cb4-e1e1-4721-8eef-806338db2222"}

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

	idStr := "f9e71cb4-e1e1-4721-8eef-806338db2222"
	userIDStr := "f9e71cb4-e1e1-4721-8eef-806338db7282"
	groupIDStr := "f9e71cb4-e1e1-4721-8eef-806338db1111"
	expireAtStr := "2025-04-16T11:28:07+03:00"

	testRequest := shortener.UpdateLinkRequest{
		Id:        idStr,
		UserId:    userIDStr,
		GroupId:   &groupIDStr,
		Generated: true,
		ShortLink: "http://qwerty",
		Url:       "http://qwertysdijvnisdnc",
		ExpireAt:  &expireAtStr,
	}

	testResponse, err := helper.GetValidatedUserId(&testRequest)
	if err != nil {
		t.Fatalf("Unable to get id, %v", err)
	}

	expected, err := uuid.Parse(userIDStr)
	if err != nil {
		t.Fatalf("Unable to parse id, %v", err)
	}

	assert.Equal(t, expected, testResponse)

}

func TestGetValidatedGroupId(t *testing.T) {

	idStr := "f9e71cb4-e1e1-4721-8eef-806338db2222"
	userIDStr := "f9e71cb4-e1e1-4721-8eef-806338db7282"
	groupIDStr := "f9e71cb4-e1e1-4721-8eef-806338db1111"
	expireAtStr := "2025-04-16T11:28:07+03:00"

	testRequest := shortener.UpdateLinkRequest{
		Id:        idStr,
		UserId:    userIDStr,
		GroupId:   &groupIDStr,
		Generated: true,
		ShortLink: "http://qwerty",
		Url:       "http://qwertysdijvnisdnc",
		ExpireAt:  &expireAtStr,
	}

	testResponse, err := helper.GetValidatedGroupId(&testRequest, nil)
	if err != nil {
		t.Fatalf("Unable to get id, %v", err)
	}

	parsed, err := uuid.Parse(groupIDStr)
	expected := &parsed
	if err != nil {
		t.Fatalf("Unable to parse id, %v", err)
	}

	assert.Equal(t, expected, testResponse)

}

func TestGetValidatedExpireAt(t *testing.T) {

	idStr := "f9e71cb4-e1e1-4721-8eef-806338db2222"
	userIDStr := "f9e71cb4-e1e1-4721-8eef-806338db7282"
	groupIDStr := "f9e71cb4-e1e1-4721-8eef-806338db1111"
	expireAtStr := "2025-04-16T11:28:07+03:00"

	testRequest := shortener.UpdateLinkRequest{
		Id:        idStr,
		UserId:    userIDStr,
		GroupId:   &groupIDStr,
		Generated: true,
		ShortLink: "http://qwerty",
		Url:       "http://qwertysdijvnisdnc",
		ExpireAt:  &expireAtStr,
	}

	testResponse, err := helper.GetValidatedExpireAt(&testRequest, time.Now())
	if err != nil {
		t.Fatalf("Unable to get ExpireAt, %v", err)
	}

	expected, err := time.Parse(time.RFC3339, *testRequest.ExpireAt)
	if err != nil {
		t.Fatalf("Unable to parse ExpireAt, %v", err)
	}

	assert.Equal(t, expected, testResponse)

}

func TestValidatedStringNotEmpty(t *testing.T) {

	idStr := "f9e71cb4-e1e1-4721-8eef-806338db2222"
	userIDStr := "f9e71cb4-e1e1-4721-8eef-806338db7282"
	groupIDStr := "f9e71cb4-e1e1-4721-8eef-806338db1111"
	expireAtStr := "2025-04-16T11:28:07+03:00"

	testRequest := shortener.UpdateLinkRequest{
		Id:        idStr,
		UserId:    userIDStr,
		GroupId:   &groupIDStr,
		Generated: true,
		ShortLink: "http://qwerty",
		Url:       "http://qwertysdijvnisdnc",
		ExpireAt:  &expireAtStr,
	}

	testerror := helper.ValidateStringNotEmpty(testRequest.GetShortLink())
	if testerror != nil {
		t.Fatalf("ValidateStringNotEmpty error, %v", testerror)
	}

	assert.NoError(t, testerror)

}

func TestValidatedUrl(t *testing.T) {

	idStr := "f9e71cb4-e1e1-4721-8eef-806338db2222"
	userIDStr := "f9e71cb4-e1e1-4721-8eef-806338db7282"
	groupIDStr := "f9e71cb4-e1e1-4721-8eef-806338db1111"
	expireAtStr := "2025-04-16T11:28:07+03:00"

	testRequest := shortener.UpdateLinkRequest{
		Id:        idStr,
		UserId:    userIDStr,
		GroupId:   &groupIDStr,
		Generated: true,
		ShortLink: "http://qwerty",
		Url:       "http://qwertysdijvnisdnc",
		ExpireAt:  &expireAtStr,
	}

	testerror := helper.ValidateUrl(testRequest.GetUrl())
	if testerror != nil {
		t.Fatalf("ValidateUrl error, %v", testerror)
	}

	assert.NoError(t, testerror)

}

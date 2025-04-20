package tests

import (
	"testing"
	"time"
	"xlink/common/gen/shortener"
	"xlink/shortener/internal/models"
	"xlink/shortener/internal/service/helper"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinkModelFromLinkRequest(t *testing.T) {

	userIDStr := "f9e71cb4-e1e1-4721-8eef-806338db7282"
	groupIDStr := "f9e71cb4-e1e1-4721-8eef-806338db1111"
	expireAtStr := "2025-04-16T11:28:07+03:00"

	userUUID, err := uuid.Parse(userIDStr)
	require.NoError(t, err)

	groupUUID, err := uuid.Parse(groupIDStr)
	require.NoError(t, err)

	var createdAtTime time.Time

	expireAtTime, err := time.Parse(time.RFC3339, expireAtStr)
	require.NoError(t, err)

	testRequest := shortener.UpdateLinkRequest{
		UserId:    userIDStr,
		GroupId:   &groupIDStr,
		Generated: true,
		ShortLink: "http://qwerty",
		Url:       "http://qwertysdijvnisdnc",
		ExpireAt:  &expireAtStr,
	}

	expectedModel := &models.Link{
		UserId:    userUUID,
		GroupId:   &groupUUID,
		Generated: true,
		ShortLink: "http://qwerty",
		Url:       "http://qwertysdijvnisdnc",
		CreatedAt: createdAtTime,
		ExpireAt:  expireAtTime,
	}

	testResponse, err := helper.LinkModelFromLinkRequest(&testRequest, time.Now())
	if err != nil {
		t.Fatalf("Unable to get link model, %v", err)
	}

	assert.Equal(t, expectedModel, testResponse)

}

func TestLinkModelFromLinkRequestWithId(t *testing.T) {

	idStr := "f9e71cb4-e1e1-4721-8eef-806338db2222"
	userIDStr := "f9e71cb4-e1e1-4721-8eef-806338db7282"
	groupIDStr := "f9e71cb4-e1e1-4721-8eef-806338db1111"
	expireAtStr := "2025-04-16T11:28:07+03:00"

	idUUID, err := uuid.Parse(idStr)
	require.NoError(t, err)

	userUUID, err := uuid.Parse(userIDStr)
	require.NoError(t, err)

	groupUUID, err := uuid.Parse(groupIDStr)
	require.NoError(t, err)

	var createdAtTime time.Time

	expireAtTime, err := time.Parse(time.RFC3339, expireAtStr)
	require.NoError(t, err)

	testRequest := shortener.UpdateLinkRequest{
		Id:        idStr,
		UserId:    userIDStr,
		GroupId:   &groupIDStr,
		Generated: true,
		ShortLink: "http://qwerty",
		Url:       "http://qwertysdijvnisdnc",
		ExpireAt:  &expireAtStr,
	}

	expectedModel := &models.Link{
		Id:        idUUID,
		UserId:    userUUID,
		GroupId:   &groupUUID,
		Generated: true,
		ShortLink: "http://qwerty",
		Url:       "http://qwertysdijvnisdnc",
		CreatedAt: createdAtTime,
		ExpireAt:  expireAtTime,
	}

	testResponse, err := helper.LinkModelFromLinkRequestWithId(&testRequest, time.Now())
	if err != nil {
		t.Fatalf("Unable to get link model, %v", err)
	}

	assert.Equal(t, expectedModel, testResponse)

}

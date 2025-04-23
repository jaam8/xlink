package helper

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
	"time"
	"xlink/common/logger"
)

type ClicksRequest interface {
	GetLinkOwner() string
	GetShortLink() string
	GetStartDate() *timestamppb.Timestamp
	GetEndDate() *timestamppb.Timestamp
}

func ValidateStartDate(request ClicksRequest) (*time.Time, error) {
	var startDate *time.Time = nil
	if request.GetStartDate() != nil {
		startDateValue := request.GetStartDate().AsTime()
		if startDateValue.After(time.Now()) {
			return nil, fmt.Errorf("start date is out of date")
		}
		startDate = &startDateValue
	}

	return startDate, nil
}

func ValidateEndDate(request ClicksRequest) (*time.Time, error) {
	var endDate *time.Time = nil
	if request.GetEndDate() != nil {
		endDateValue := request.GetEndDate().AsTime()
		oneMonthAgo := time.Now().AddDate(0, -1, 0)
		if endDateValue.Before(oneMonthAgo) {
			return nil, fmt.Errorf("start date is too old (more than 1 month ago)")
		}
		endDate = &endDateValue
	}

	return endDate, nil
}

func ValidateNotEmptyStr(str string) (string, error) {
	if len(strings.TrimSpace(str)) == 0 {
		return "", fmt.Errorf("string cannot be empty")
	}
	return str, nil
}

func ValidateRequestDates(ctx context.Context, request ClicksRequest) (*time.Time, *time.Time, error) {
	startDate, err := ValidateStartDate(request)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to validate start date: %v", zap.Error(err))
		return nil, nil, err
	}
	endDate, err := ValidateEndDate(request)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to validate end date: %v", zap.Error(err))
		return nil, nil, err
	}
	return startDate, endDate, nil
}

func ValidateRequestLinkOwner(ctx context.Context, request ClicksRequest) (uuid.UUID, error) {
	requestLinkOwner, err := uuid.Parse(request.GetLinkOwner())
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx,
			"failed to parse link owner: %v", zap.Error(err))
		return uuid.UUID{}, err
	}

	return requestLinkOwner, nil
}

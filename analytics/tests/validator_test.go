package tests

import (
	"context"
	"testing"
	"time"

	"xlink/analytics/internal/service/helper"
	"xlink/common/logger"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type mockClicksRequest struct {
	startDate *timestamppb.Timestamp
}

func (m *mockClicksRequest) GetLinkOwner() string                 { return "" }
func (m *mockClicksRequest) GetShortLink() string                 { return "" }
func (m *mockClicksRequest) GetStartDate() *timestamppb.Timestamp { return m.startDate }
func (m *mockClicksRequest) GetEndDate() *timestamppb.Timestamp   { return nil }

func TestValidateStartDate(t *testing.T) {

	t.Run("valid past date", func(t *testing.T) {
		past := time.Now().Add(-24 * time.Hour)
		req := &mockClicksRequest{startDate: timestamppb.New(past)}

		result, err := helper.ValidateStartDate(req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result == nil || !result.Equal(past) {
			t.Errorf("expected %v, got %v", past, result)
		}
	})

	t.Run("invalid future date", func(t *testing.T) {
		future := time.Now().Add(24 * time.Hour)
		req := &mockClicksRequest{startDate: timestamppb.New(future)}

		_, err := helper.ValidateStartDate(req)
		if err == nil {
			t.Fatal("expected error, got none")
		}
	})
}

type mockRequest struct {
	startDate *timestamppb.Timestamp
	endDate   *timestamppb.Timestamp
}

func (m *mockRequest) GetStartDate() *timestamppb.Timestamp {
	return m.startDate
}

func (m *mockRequest) GetEndDate() *timestamppb.Timestamp {
	return m.endDate
}

func (m *mockRequest) GetLinkOwner() string {
	return "dummy-owner"
}

func (m *mockRequest) GetShortLink() string {
	return "dummy-short-link"
}

func TestValidateEndDate(t *testing.T) {
	tests := []struct {
		name        string
		endDate     time.Time
		expectError bool
	}{
		{
			name:        "valid end date",
			endDate:     time.Now().AddDate(0, 0, -10),
			expectError: false,
		},
		{
			name:        "too old end date",
			endDate:     time.Now().AddDate(0, -2, 0),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &mockRequest{
				endDate: timestamppb.New(tt.endDate),
			}
			result, err := helper.ValidateEndDate(req)
			if tt.expectError && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("did not expect error but got: %v", err)
			}
			if !tt.expectError && result == nil {
				t.Errorf("expected valid time, got nil")
			}
		})
	}
}

func TestValidateNotEmptyStr(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "non-empty string",
			input:       "valid input",
			expectError: false,
		},
		{
			name:        "empty string",
			input:       "",
			expectError: true,
		},
		{
			name:        "string with only spaces",
			input:       "   ",
			expectError: true,
		},
		{
			name:        "string with leading and trailing spaces",
			input:       "  valid  ",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := helper.ValidateNotEmptyStr(tt.input)
			if tt.expectError && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("did not expect error but got: %v", err)
			}
			if !tt.expectError && result != tt.input {
				t.Errorf("expected '%s', got '%s'", tt.input, result)
			}
		})
	}
}

type mockRequest1 struct {
	start *timestamppb.Timestamp
	end   *timestamppb.Timestamp
}

func (m *mockRequest1) GetStartDate() *timestamppb.Timestamp {
	return m.start
}
func (m *mockRequest1) GetEndDate() *timestamppb.Timestamp {
	return m.end
}
func (m *mockRequest1) GetLinkOwner() string {
	return "dummy"
}
func (m *mockRequest1) GetShortLink() string {
	return "short"
}

func TestValidateRequestDates(t *testing.T) {

	now := time.Now()
	validStart := timestamppb.New(now.Add(-24 * time.Hour))
	validEnd := timestamppb.New(now.Add(-1 * time.Hour))
	oldEnd := timestamppb.New(now.AddDate(0, -2, 0))
	futureStart := timestamppb.New(now.Add(24 * time.Hour))

	tests := []struct {
		name        string
		request     *mockRequest1
		expectError bool
	}{
		{
			name:        "valid dates",
			request:     &mockRequest1{start: validStart, end: validEnd},
			expectError: false,
		},
		{
			name:        "start date in future",
			request:     &mockRequest1{start: futureStart, end: validEnd},
			expectError: true,
		},
		{
			name:        "end date too old",
			request:     &mockRequest1{start: validStart, end: oldEnd},
			expectError: true,
		},
		{
			name:        "nil dates",
			request:     &mockRequest1{start: nil, end: nil},
			expectError: false,
		},
	}

	ctx, err := logger.New(context.Background())
	if err != nil {
		t.Errorf("cannot implement logger:%v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end, err := helper.ValidateRequestDates(ctx, tt.request)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, start)
				assert.Nil(t, end)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type mockRequest2 struct {
	linkOwner string
}

func (m *mockRequest2) GetLinkOwner() string                 { return m.linkOwner }
func (m *mockRequest2) GetShortLink() string                 { return "short" }
func (m *mockRequest2) GetStartDate() *timestamppb.Timestamp { return nil }
func (m *mockRequest2) GetEndDate() *timestamppb.Timestamp   { return nil }

func TestValidateRequestLinkOwner(t *testing.T) {

	validUUID := uuid.New().String()
	tests := []struct {
		name      string
		input     string
		expectErr bool
	}{
		{"Valid UUID", validUUID, false},
		{"Invalid UUID", "not-a-uuid", true},
	}

	ctx, err := logger.New(context.Background())
	if err != nil {
		t.Errorf("cannot implement logger:%v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &mockRequest2{linkOwner: tt.input}
			_, err := helper.ValidateRequestLinkOwner(ctx, req)
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
		})
	}
}

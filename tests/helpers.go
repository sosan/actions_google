package tests

import (
	"actions_google/pkg/domain/models"
	"time"
)

func StringPtr(s string) *string { return &s }
func Int64Ptr(i int64) *int64    { return &i }
func BoolPtr(b bool) *bool       { return &b }
func Float64Ptr(f float64) *float64 {
	return &f
}

func CustomTime(t time.Time) *models.CustomTime {
	return &models.CustomTime{Time: t}
}

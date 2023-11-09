package common

import "time"

type OrderInfo struct {
	TaskId    string
	UserEmail []string
}

const (
	ExpiredTimeRdb = 24 * time.Hour
)

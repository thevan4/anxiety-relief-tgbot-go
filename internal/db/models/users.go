package models

import "time"

type User struct {
	UserID        int64
	Username      string
	IsPremium     bool
	IsBot         bool
	TotalRequests int64
	FirstSeen     time.Time
	LastSeen      time.Time
}

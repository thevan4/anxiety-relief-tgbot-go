package models

import "time"

type TechniqueUsage struct {
	ID          int64
	UserID      int64
	TechniqueID string
	Category    string
	StartedAt   time.Time
	Completed   bool
	CompletedAt *time.Time
}

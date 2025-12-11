package statistic

import "sync"

// Stats ...
type Stats interface {
	IncreaseRequestsStatisticForUser(
		userID int64,
		username string,
		isPremium bool,
		isBot bool,
	)
	GetStatistics() map[int64]Statistic
}

// Statistics ...
type Statistics struct {
	mux *sync.Mutex
	// Key - user id.
	stats map[int64]*Statistic
}

// NewStatistics loads past statistics into memory.
func NewStatistics(oldStatistics map[int64]*Statistic) *Statistics {
	var newStats map[int64]*Statistic
	if len(oldStatistics) == 0 {
		newStats = make(map[int64]*Statistic)
	} else {
		newStats = oldStatistics
	}
	return &Statistics{
		mux:   new(sync.Mutex),
		stats: newStats,
	}
}

// IncreaseRequestsStatisticForUser increases the request counter for the user.
func (s *Statistics) IncreaseRequestsStatisticForUser(
	userID int64,
	username string,
	isPremium bool,
	isBot bool,
) {
	s.mux.Lock()
	defer s.mux.Unlock()
	oldStats, inMap := s.stats[userID]
	if !inMap {
		s.stats[userID] = &Statistic{
			username:      username,
			isPremium:     isPremium,
			isBot:         isBot,
			totalRequests: 1,
		}
		return
	}

	oldStats.totalRequests += 1
}

// GetStatistics gives the current statistics.
func (s *Statistics) GetStatistics() map[int64]Statistic {
	s.mux.Lock()
	defer s.mux.Unlock()
	statisticCopy := make(map[int64]Statistic, len(s.stats))
	for k, v := range s.stats {
		statisticCopy[k] = *v
	}

	return statisticCopy
}

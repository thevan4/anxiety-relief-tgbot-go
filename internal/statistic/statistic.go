package statistic

// Statistic ...
type Statistic struct {
	username      string
	isPremium     bool
	isBot         bool
	totalRequests int64
}

// NewStatistic ...
func NewStatistic(
	username string,
	isPremium bool,
	isBot bool,
	totalRequests int64,
) *Statistic {
	return &Statistic{
		username:      username,
		isPremium:     isPremium,
		isBot:         isBot,
		totalRequests: totalRequests,
	}
}

// GetUsername ...
func (s *Statistic) GetUsername() string {
	if s == nil {
		return ""
	}
	return s.username
}

// GetIsPremium ...
func (s *Statistic) GetIsPremium() bool {
	if s == nil {
		return false
	}
	return s.isPremium
}

// GetIsBot ...
func (s *Statistic) GetIsBot() bool {
	if s == nil {
		return false
	}
	return s.isBot
}

// GetTotalRequests ...
func (s *Statistic) GetTotalRequests() int64 {
	if s == nil {
		return 0
	}
	return s.totalRequests
}

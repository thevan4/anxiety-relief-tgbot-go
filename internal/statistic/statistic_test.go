package statistic

import (
	"testing"
)

func TestNewStatistic(t *testing.T) {
	s := NewStatistic("testuser", true, false, 10)
	if s == nil {
		t.Fatal("NewStatistic returned nil")
	}
	if s.GetUsername() != "testuser" {
		t.Errorf("GetUsername() = %q, want %q", s.GetUsername(), "testuser")
	}
	if !s.GetIsPremium() {
		t.Error("GetIsPremium() = false, want true")
	}
	if s.GetIsBot() {
		t.Error("GetIsBot() = true, want false")
	}
	if s.GetTotalRequests() != 10 {
		t.Errorf("GetTotalRequests() = %d, want %d", s.GetTotalRequests(), 10)
	}
}

func TestStatisticNilSafety(t *testing.T) {
	var s *Statistic = nil

	if s.GetUsername() != "" {
		t.Errorf("nil.GetUsername() = %q, want empty", s.GetUsername())
	}
	if s.GetIsPremium() {
		t.Error("nil.GetIsPremium() = true, want false")
	}
	if s.GetIsBot() {
		t.Error("nil.GetIsBot() = true, want false")
	}
	if s.GetTotalRequests() != 0 {
		t.Errorf("nil.GetTotalRequests() = %d, want 0", s.GetTotalRequests())
	}
}

func TestNewStatistics(t *testing.T) {
	// Test with nil
	stats := NewStatistics(nil)
	if stats == nil {
		t.Fatal("NewStatistics(nil) returned nil")
	}
	if len(stats.stats) != 0 {
		t.Errorf("NewStatistics(nil).stats has %d elements, want 0", len(stats.stats))
	}

	// Test with existing data
	oldStats := map[int64]*Statistic{
		123: NewStatistic("user1", true, false, 5),
		456: NewStatistic("user2", false, true, 10),
	}
	stats = NewStatistics(oldStats)
	if len(stats.stats) != 2 {
		t.Errorf("NewStatistics with data has %d elements, want 2", len(stats.stats))
	}

	// Verify data was copied (not shared)
	oldStats[123].username = "modified"
	if stats.stats[123].username == "modified" {
		t.Error("NewStatistics shares map with input")
	}
}

func TestIncreaseRequestsStatisticForUser(t *testing.T) {
	stats := NewStatistics(nil)

	// First request creates new entry
	stats.IncreaseRequestsStatisticForUser(123, "user1", true, false)
	result := stats.GetStatistics()
	if len(result) != 1 {
		t.Fatalf("GetStatistics() has %d elements, want 1", len(result))
	}
	stat := result[123]
	if stat.GetTotalRequests() != 1 {
		t.Errorf("TotalRequests = %d, want 1", stat.GetTotalRequests())
	}

	// Second request increments counter
	stats.IncreaseRequestsStatisticForUser(123, "user1", true, false)
	result = stats.GetStatistics()
	stat = result[123]
	if stat.GetTotalRequests() != 2 {
		t.Errorf("TotalRequests = %d, want 2", stat.GetTotalRequests())
	}

	// Updates user info
	stats.IncreaseRequestsStatisticForUser(123, "user1_updated", false, true)
	result = stats.GetStatistics()
	stat = result[123]
	if stat.GetUsername() != "user1_updated" {
		t.Errorf("Username = %q, want %q", stat.GetUsername(), "user1_updated")
	}
	if stat.GetIsPremium() {
		t.Error("IsPremium = true, want false")
	}
	if !stat.GetIsBot() {
		t.Error("IsBot = false, want true")
	}
}

func TestGetStatisticsCopy(t *testing.T) {
	stats := NewStatistics(nil)
	stats.IncreaseRequestsStatisticForUser(123, "user1", true, false)

	result1 := stats.GetStatistics()
	result2 := stats.GetStatistics()

	// Modifying result1 should not affect result2
	delete(result1, 123)
	if _, exists := result2[123]; !exists {
		t.Error("GetStatistics() returns shared map")
	}
}

func TestStatsInterface(t *testing.T) {
	var _ Stats = NewStatistics(nil)
}

func TestConcurrentStatistics(t *testing.T) {
	stats := NewStatistics(nil)
	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func(userID int64) {
			for j := 0; j < 100; j++ {
				stats.IncreaseRequestsStatisticForUser(userID, "user", false, false)
				stats.GetStatistics()
			}
			done <- true
		}(int64(i))
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	result := stats.GetStatistics()
	for i := int64(0); i < 10; i++ {
		stat := result[i]
		if stat.GetTotalRequests() != 100 {
			t.Errorf("user %d has %d requests, want 100", i, stat.GetTotalRequests())
		}
	}
}

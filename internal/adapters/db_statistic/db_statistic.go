package db_statistic

import (
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/db/models"
	"github.com/thevan4/anxiety-relief-tgbot-go/internal/statistic"
)

func DbModelsToStatistics(users []*models.User) map[int64]*statistic.Statistic {
	statistics := make(map[int64]*statistic.Statistic, len(users))
	for _, u := range users {
		statistics[u.UserID] = statistic.NewStatistic(u.Username, u.IsPremium, u.IsBot, u.TotalRequests)
	}
	return statistics
}

func StatisticsToDbModels(stats map[int64]statistic.Statistic) []*models.User {
	users := make([]*models.User, 0, len(stats))
	for userID, stat := range stats {
		users = append(users, &models.User{
			UserID:        userID,
			Username:      stat.GetUsername(),
			IsPremium:     stat.GetIsPremium(),
			IsBot:         stat.GetIsBot(),
			TotalRequests: stat.GetTotalRequests(),
		})
	}
	return users
}

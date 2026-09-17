package models

type DashboardStats struct {
	TotalUsers          int                `json:"total_users"`
	ActiveUsers         int                `json:"active_users"`
	TotalGames          int                `json:"total_games"`
	ActiveGames         int                `json:"active_games"`
	TotalDraws          int                `json:"total_draws"`
	PendingDraws        int                `json:"pending_draws"`
	TotalTickets        int                `json:"total_tickets"`
	TotalContributions  float64            `json:"total_contributions"`
	TotalSpent          float64            `json:"total_spent"`
	TotalPrizes         float64            `json:"total_prizes"`
	PendingPayments     []PendingPayment   `json:"pending_payments"`
	ContributionsByGame []GameContribution `json:"contributions_by_game"`
	MonthlyTrend        []MonthlyData      `json:"monthly_trend"`
}

type PendingPayment struct {
	UserID   int     `json:"user_id"`
	UserName string  `json:"user_name"`
	GameID   int     `json:"game_id"`
	GameName string  `json:"game_name"`
	Month    int     `json:"month"`
	Year     int     `json:"year"`
	Amount   float64 `json:"amount"`
}

type GameContribution struct {
	GameID   int     `json:"game_id"`
	GameName string  `json:"game_name"`
	Total    float64 `json:"total"`
	Count    int     `json:"count"`
}

type MonthlyData struct {
	Month  int     `json:"month"`
	Year   int     `json:"year"`
	Label  string  `json:"label"`
	Total  float64 `json:"total"`
	Count  int     `json:"count"`
}

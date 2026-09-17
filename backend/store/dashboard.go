package store

import (
	"context"
	"fmt"
	"time"

	"lottery-pool-manager/models"
)

func (s *PostgreSQLStore) GetDashboardStats() (models.DashboardStats, error) {
	var stats models.DashboardStats
	ctx := context.Background()

	// User counts
	err := s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&stats.TotalUsers)
	if err != nil {
		return stats, fmt.Errorf("count users: %w", err)
	}
	err = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE active = true`).Scan(&stats.ActiveUsers)
	if err != nil {
		return stats, fmt.Errorf("count active users: %w", err)
	}

	// Game counts
	err = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM lottery_games`).Scan(&stats.TotalGames)
	if err != nil {
		return stats, fmt.Errorf("count games: %w", err)
	}
	err = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM lottery_games WHERE active = true`).Scan(&stats.ActiveGames)
	if err != nil {
		return stats, fmt.Errorf("count active games: %w", err)
	}

	// Draw counts
	err = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM draws`).Scan(&stats.TotalDraws)
	if err != nil {
		return stats, fmt.Errorf("count draws: %w", err)
	}
	err = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM draws WHERE processed = false`).Scan(&stats.PendingDraws)
	if err != nil {
		return stats, fmt.Errorf("count pending draws: %w", err)
	}

	// Ticket counts
	err = s.pool.QueryRow(ctx, `SELECT COUNT(*) FROM tickets`).Scan(&stats.TotalTickets)
	if err != nil {
		return stats, fmt.Errorf("count tickets: %w", err)
	}

	// Financial totals
	err = s.pool.QueryRow(ctx, `SELECT COALESCE(SUM(amount), 0) FROM contributions`).Scan(&stats.TotalContributions)
	if err != nil {
		return stats, fmt.Errorf("sum contributions: %w", err)
	}
	err = s.pool.QueryRow(ctx, `SELECT COALESCE(SUM(cost), 0) FROM tickets`).Scan(&stats.TotalSpent)
	if err != nil {
		return stats, fmt.Errorf("sum spent: %w", err)
	}
	err = s.pool.QueryRow(ctx, `SELECT COALESCE(SUM(COALESCE(prize_amount, 0)), 0) FROM tickets`).Scan(&stats.TotalPrizes)
	if err != nil {
		return stats, fmt.Errorf("sum prizes: %w", err)
	}

	// Pending payments (current month)
	now := time.Now()
	currentMonth := int(now.Month())
	currentYear := now.Year()
	stats.PendingPayments, err = s.getPendingPayments(currentMonth, currentYear)
	if err != nil {
		return stats, fmt.Errorf("pending payments: %w", err)
	}

	// Contributions by game
	stats.ContributionsByGame, err = s.getContributionsByGame()
	if err != nil {
		return stats, fmt.Errorf("contributions by game: %w", err)
	}

	// Monthly trend (last 6 months)
	stats.MonthlyTrend, err = s.getMonthlyTrend(6)
	if err != nil {
		return stats, fmt.Errorf("monthly trend: %w", err)
	}

	return stats, nil
}

func (s *PostgreSQLStore) getPendingPayments(month, year int) ([]models.PendingPayment, error) {
	query := `
		SELECT u.id, u.name, lg.id, lg.name, c.month, c.year, c.amount
		FROM contributions c
		JOIN users u ON c.user_id = u.id
		JOIN lottery_games lg ON c.game_id = lg.id
		WHERE c.paid = false AND c.month = $1 AND c.year = $2
		ORDER BY u.name, lg.name
	`
	rows, err := s.pool.Query(context.Background(), query, month, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.PendingPayment
	for rows.Next() {
		var p models.PendingPayment
		if err := rows.Scan(&p.UserID, &p.UserName, &p.GameID, &p.GameName, &p.Month, &p.Year, &p.Amount); err != nil {
			return nil, err
		}
		payments = append(payments, p)
	}
	if payments == nil {
		payments = []models.PendingPayment{}
	}
	return payments, nil
}

func (s *PostgreSQLStore) getContributionsByGame() ([]models.GameContribution, error) {
	query := `
		SELECT lg.id, lg.name, COALESCE(SUM(c.amount), 0), COUNT(c.id)
		FROM lottery_games lg
		LEFT JOIN contributions c ON lg.id = c.game_id
		WHERE lg.active = true
		GROUP BY lg.id, lg.name
		ORDER BY SUM(c.amount) DESC
	`
	rows, err := s.pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.GameContribution
	for rows.Next() {
		var g models.GameContribution
		if err := rows.Scan(&g.GameID, &g.GameName, &g.Total, &g.Count); err != nil {
			return nil, err
		}
		results = append(results, g)
	}
	if results == nil {
		results = []models.GameContribution{}
	}
	return results, nil
}

func (s *PostgreSQLStore) getMonthlyTrend(months int) ([]models.MonthlyData, error) {
	query := `
		SELECT 
			EXTRACT(MONTH FROM created_at)::int,
			EXTRACT(YEAR FROM created_at)::int,
			COALESCE(SUM(amount), 0),
			COUNT(id)
		FROM contributions
		WHERE created_at >= NOW() - INTERVAL '6 months'
		GROUP BY EXTRACT(MONTH FROM created_at), EXTRACT(YEAR FROM created_at)
		ORDER BY EXTRACT(YEAR FROM created_at), EXTRACT(MONTH FROM created_at)
	`
	rows, err := s.pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	monthNames := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

	var results []models.MonthlyData
	for rows.Next() {
		var m models.MonthlyData
		if err := rows.Scan(&m.Month, &m.Year, &m.Total, &m.Count); err != nil {
			return nil, err
		}
		if m.Month >= 1 && m.Month <= 12 {
			m.Label = monthNames[m.Month-1]
		}
		results = append(results, m)
	}
	if results == nil {
		results = []models.MonthlyData{}
	}
	return results, nil
}

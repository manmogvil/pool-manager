package models

import "time"

type Ticket struct {
	ID             int        `json:"id"`
	DrawID         int        `json:"draw_id"`
	Numbers        string     `json:"numbers"`
	Stars          *string    `json:"stars"`
	Cost           float64    `json:"cost"`
	PurchasedAt    time.Time  `json:"purchased_at"`
	PrizeTier      *string    `json:"prize_tier"`
	PrizeAmount    *float64   `json:"prize_amount"`
	MatchedNumbers *int       `json:"matched_numbers"`
	MatchedStars   *int       `json:"matched_stars"`
}

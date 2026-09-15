package models

import (
	"time"
)

// PaymentDate is a pointer to time.Time to allow for null values in the database (optional field)
type Contribution struct {
	ID            int        `json:"id"`
	UserID        int        `json:"user_id"`
	GameID        int        `json:"game_id"`
	Month         int        `json:"month"`
	Year          int        `json:"year"`
	Amount        float64    `json:"amount"`
	Paid          bool       `json:"paid"`
	PaymentDate   *time.Time `json:"payment_date"`
	PaymentMethod string     `json:"payment_method"`
	Comments      string     `json:"comments"`
	CreatedAt     time.Time  `json:"created_at"`
}

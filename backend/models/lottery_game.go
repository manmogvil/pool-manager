package models

type LotteryGame struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	DrawDays    string  `json:"draw_days"`
	TicketPrice float64 `json:"ticket_price"`
	Active      bool    `json:"active"`
}

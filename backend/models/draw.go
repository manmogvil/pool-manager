package models

import "time"

type Draw struct {
	ID            int       `json:"id"`
	GameID        int       `json:"game_id"`
	DrawDate      time.Time `json:"draw_date"`
	ResultNumbers *string   `json:"result_numbers"`
	ResultStars   *string   `json:"result_stars"`
	DrawIDAPI     *string   `json:"draw_id_api"`
	Processed     bool      `json:"processed"`
	CreatedAt     time.Time `json:"created_at"`
}

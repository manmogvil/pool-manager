package models

import "time"

// Con la etiqueta json se indica cómo se deben serializar los campos de la estructura al convertirla a JSON.
// Por ejemplo, el campo Name como "name"
type Participant struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
}

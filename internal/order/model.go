package order

import "time"

type Order struct {
	ID        int64     `json:"id"`
	Customer  string    `json:"customer"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

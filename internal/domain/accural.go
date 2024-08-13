package domain

import "time"

type Accural struct {
	ID          uint64
	OrderID     uint64 `db:"order_id"`
	Sum         uint64
	ProcessedAt time.Time `db:"processed_at"`
}

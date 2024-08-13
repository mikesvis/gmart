package domain

import "time"

type Status string

const (
	StatusNew        Status = "NEW"
	StatusProcessing Status = "PROCESSING"
	StatusInvalid    Status = "INVALID"
	StatusProcessed  Status = "PROCESSED"
)

type Order struct {
	ID        uint64
	UserID    uint64 `db:"user_id"`
	Status    Status
	Accural   float64
	CreatedAt time.Time `db:"created_at"`
}

func (s Status) String() string {
	return string(s)
}

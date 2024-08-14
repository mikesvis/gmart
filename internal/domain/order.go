package domain

import (
	"github.com/mikesvis/gmart/pkg/json"
)

type Status string

const (
	StatusNew        Status = "NEW"
	StatusProcessing Status = "PROCESSING"
	StatusInvalid    Status = "INVALID"
	StatusProcessed  Status = "PROCESSED"
)

type Order struct {
	ID        uint64        `db:"id" json:"number,string"`
	UserID    uint64        `db:"user_id" json:"-"`
	Status    Status        `db:"status" json:"status"`
	Accural   json.Kopeykis `db:"amount" json:"accural,omitempty"`
	CreatedAt json.JSONTime `db:"created_at" json:"uploaded_at"`
}

func (s Status) String() string {
	return string(s)
}

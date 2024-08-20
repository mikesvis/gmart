package domain

import "time"

type User struct {
	ID       uint64
	Login    string
	Password string
}

type UserBalance struct {
	Current   uint64 // ВАЖНО! значения храним только в целых числах (копейки а не рубли)
	Withdrawn uint64 // ВАЖНО! значения храним только в целых числах (копейки а не рубли)
}

type UserWithdrawals struct {
	OrderID     uint64    `db:"id"`
	Sum         uint64    `db:"sum"` // ВАЖНО! значения храним только в целых числах (копейки а не рубли)
	ProcessedAt time.Time `db:"processed_at"`
}

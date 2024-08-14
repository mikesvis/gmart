package domain

type User struct {
	ID       uint64
	Login    string
	Password string
}

type UserBalance struct {
	Current   uint64 // ВАЖНО! значения храним только в целых числах (копейки а не рубли)
	Withdrawn uint64 // ВАЖНО! значения храним только в целых числах (копейки а не рубли)
}

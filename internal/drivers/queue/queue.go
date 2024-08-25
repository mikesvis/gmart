package drivers

type Queue struct {
	OrderQueue chan uint64
	Quit       chan bool
}

func NewQueue() *Queue {
	queue := Queue{
		OrderQueue: make(chan uint64),
		Quit:       make(chan bool, 1),
	}
	return &queue
}

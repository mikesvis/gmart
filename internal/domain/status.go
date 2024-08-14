package domain

type Status string

const (
	StatusNew        Status = "NEW"
	StatusRegistered Status = "REGISTERED"
	StatusProcessing Status = "PROCESSING"
	StatusInvalid    Status = "INVALID"
	StatusProcessed  Status = "PROCESSED"
)

func (s Status) String() string {
	return string(s)
}

package json

import (
	"strconv"
	"time"
)

type JSONTime time.Time

func (v JSONTime) MarshalJSON() ([]byte, error) {
	stamp := "\"" + time.Time(v).Local().Format(time.RFC3339) + "\""
	return []byte(stamp), nil
}

type Kopeykis uint64

func (v Kopeykis) MarshalJSON() ([]byte, error) {
	rubles := float64(v / 100)
	return []byte(strconv.FormatFloat(rubles, 'f', 6, 64)), nil
}

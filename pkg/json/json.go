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

type Rubles float64

func (v Rubles) MarshalJSON() ([]byte, error) {
	rubles := float64(v / 100)
	return []byte(strconv.FormatFloat(rubles, 'f', 2, 64)), nil
}

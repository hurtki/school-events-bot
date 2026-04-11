package domain

import (
	"errors"
	"strings"
	"time"
)

type Date struct {
	t time.Time
}

var (
	ErrWrongFormat = errors.New("wrong date format")
)

const dateLayout = "2.1.2006"

func NewDate(d string) (Date, error) {
	d = strings.TrimSpace(d)

	parsedTime, err := time.Parse(dateLayout, d)
	if err != nil {
		return Date{}, ErrWrongFormat
	}

	return Date{t: parsedTime}, nil
}

func (d Date) String() string {
	return d.t.Format(dateLayout)
}

func (d Date) Compare(other Date) int {
	if d.t.Before(other.t) {
		return -1
	}
	if d.t.After(other.t) {
		return 1
	}
	return 0
}

package domain

import (
	"errors"
	"strings"
	"time"
)

type Date struct {
	T time.Time
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

	return Date{T: parsedTime}, nil
}

// func NewTodaysDate() Date {
// 	now := time.Now()
// 	d, err := NewDate(
// 		fmt.Sprintf(
// 			"%d.%d.%d",
// 			now.Day(),
// 			now.Month(),
// 			now.Year(),
// 		),
// 	)
// 	if err != nil {
// 		panic("unexpected error from new date builder")
// 	}
// 	return d
// }

func (d Date) String() string {
	return d.T.Format(dateLayout)
}

func (d Date) Compare(other Date) int {
	if d.T.Before(other.T) {
		return -1
	}
	if d.T.After(other.T) {
		return 1
	}
	return 0
}

func (d Date) DaysUntil() int {
	now := time.Now().Truncate(24 * time.Hour)
	target := d.T.Truncate(24 * time.Hour)
	return int(target.Sub(now).Hours() / 24)
}

package domain

import (
	"errors"
	"strconv"
	"strings"
)

type Date struct {
	Year  int
	Month int
	Day   int
}

var (
	ErrWrongFormat = errors.New("wrong date format")
)

// day.month.year format
func NewDate(d string) (Date, error) {
	d = strings.TrimSpace(d)
	parts := strings.Split(d, ".")
	if len(parts) != 3 {
		return Date{}, ErrWrongFormat
	}
	day, err := strconv.Atoi(parts[0])
	if err != nil || day < 1 || day > 31 {
		return Date{}, ErrWrongFormat
	}
	month, err := strconv.Atoi(parts[1])
	if err != nil || month < 1 || month > 12 {
		return Date{}, ErrWrongFormat
	}
	year, err := strconv.Atoi(parts[2])
	if err != nil || year < 2000 || year > 2100 {
		return Date{}, ErrWrongFormat
	}
	return Date{Year: year, Month: month, Day: day}, nil
}

type Event struct {
	Date Date

	Group string
	Text  string
}

func NewEvent(date Date, gr string, text string) (Event, error) {
	return Event{
		Date:  date,
		Group: gr,
		Text:  text,
	}, nil
}

type Schedule struct {
	Events []Event
}

func NewSchedule(evs []Event) (Schedule, error) {
	return Schedule{
		Events: evs,
	}, nil
}

// t :=

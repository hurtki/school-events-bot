package domain

import "time"

type Date struct {
	Year int
	Month int
	Day int 
}

type Schedule struct {
	events map[Date]string
}

func test() {
	z, _ := time.LoadLocation("Asia/Jerusalem")
	t := time.Date(2026, time.April, 2, 0, 0, 0, 0, z)
	k := make(map[time.Time]string)
}

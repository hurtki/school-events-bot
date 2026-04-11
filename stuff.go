package main

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/hurtki/school-events-bot/internal/domain"
	"github.com/hurtki/school-events-bot/internal/infrastructure/spreadsheets"
	"github.com/hurtki/school-events-bot/internal/parser"
)

func main1(fetcher *spreadsheets.DocsFetcher) {
	ctx := context.Background()

	xlsx, err := fetcher.FetchXLSX(ctx)
	defer func() {
		if err := xlsx.Close(); err != nil {
			fmt.Println("coulnd't close xlsx doc")
		}
	}()
	sc, err := parser.ParseXLSX(xlsx)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Scanned %d events\n", len(sc.Events))

	start := time.Now()
	slices.SortFunc(sc.Events, func(a domain.Event, b domain.Event) int {
		if a.Date.Year != b.Date.Year {
			return a.Date.Year - b.Date.Year
		}
		if a.Date.Month != b.Date.Month {
			return a.Date.Month - b.Date.Month
		}
		if a.Date.Day != b.Date.Day {
			return a.Date.Day - b.Date.Day
		}
		return 0
	})
	fmt.Println("sorted in", time.Since(start).String())
	for _, ev := range sc.Events {
		// if !strings.Contains(ev.Group, "א") {
		if !strings.Contains(ev.Group, "ב") {
			continue
		}
		text := ""
		if hasHebrew(ev.Text) {
			text = reverseStringKeepLines(ev.Text)
		} else {
			text = ev.Text
		}

		fmt.Printf("[%d.%d.%d] [%s] [%s] \n%s",
			ev.Date.Day,
			ev.Date.Month,
			ev.Date.Year,
			reverseStringKeepLines(ev.Group),
			ev.Type.String(),
			text,
		)
	}

}

func hasHebrew(s string) bool {
	for _, r := range s {
		if r >= 0x0590 && r <= 0x05FF {
			return true
		}
	}
	return false
}

func reverseStringKeepLines(s string) string {
	lines := []rune(s)

	start := 0
	for i := 0; i <= len(lines); i++ {
		if i == len(lines) || lines[i] == '\n' {
			reverse(lines[start:i])
			start = i + 1
		}
	}
	return string(lines)
}

func reverse(s []rune) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

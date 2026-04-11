package main

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/hurtki/school-events-bot/internal/config"
	"github.com/hurtki/school-events-bot/internal/domain"
	"github.com/hurtki/school-events-bot/internal/infrastructure/spreadsheets"
	"github.com/hurtki/school-events-bot/internal/parser"
)

func main1(fetcher *spreadsheets.DocsFetcher, cfg config.AppConfig) {
	ctx := context.Background()

	xlsx, err := fetcher.FetchXLSX(ctx, cfg.SpreadsheetsDocumentID)
	defer func() {
		if err := xlsx.Close(); err != nil {
			fmt.Println("coulnd't close xlsx doc")
		}
	}()
	p, _ := parser.NewParser(xlsx, cfg.SpreadsheetsDocumentID)
	sc, err := p.ParseXLSX()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Scanned %d events\n", len(sc.Events))

	start := time.Now()
	slices.SortFunc(sc.Events, func(a, b domain.Event) int {
		return a.Date.Compare(b.Date)
	})
	fmt.Println("sorted in", time.Since(start).String())
	for _, ev := range sc.Events {
		// if !strings.Contains(ev.Group, "א") {
		if !(ev.Group == domain.EleventhGradeGroup) {
			continue
		}
		text := ""
		if hasHebrew(ev.Text) {
			text = reverseStringKeepLines(ev.Text)
		} else {
			text = ev.Text
		}

		fmt.Printf("[%s] [%s] [%s] \n%ssource link: %s\n",
			ev.Date.String(),
			ev.Group.String(),
			ev.Type.String(),
			text,
			ev.SourceURL,
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

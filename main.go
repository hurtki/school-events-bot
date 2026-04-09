package main

import (
	"fmt"
	"os"

	"github.com/hurtki/school-events-bot/internal/parser"
)

// const xlsxURL = "https://docs.google.com/spreadsheets/d/1WAqZExNwrM9w2p3IbOkS6ZMosKioh66h/export?format=xlsx"

func main() {
	f, err := os.Open("tbl.xlsx")
	if err != nil {
		fmt.Println("err when opening", err)
		return
	}
	sc, err := parser.ParseXLSX(f)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Scanned %d events\n", len(sc.Events))
	for _, ev := range sc.Events {
		text := ""
		if hasHebrew(ev.Text) {
			text = reverse(ev.Text)
		} else {
			text = ev.Text
		}

		fmt.Printf("%d.%d.%d [%s]: %s\n",
			ev.Date.Day,
			ev.Date.Month,
			ev.Date.Year,
			reverse(ev.Group),
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

func reverse(t string) string {
	s := []rune(t)
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return string(s)
}

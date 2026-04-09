package main

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

const xlsxURL = "https://docs.google.com/spreadsheets/d/1WAqZExNwrM9w2p3IbOkS6ZMosKioh66h/export?format=xlsx"

func main() {
	// res, err := http.Get(xlsxURL)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// f, _ := excelize.OpenReader(res.Body)
	// fmt.Println(f.GetSheetList())

	f, err := excelize.OpenFile("tbl.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	levels := f.GetSheetList()
	for _, l := range levels {
		rows, _ := f.GetRows(l)
		for _, r := range rows {
			for _, e := range r {
				if hasHebrew(e) {
					fmt.Printf("%s ", reverse(e))
				} else {
					fmt.Printf("%s ", e)
				}
			}
			fmt.Println()
		}
		fmt.Println("=======")
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

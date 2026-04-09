package parser

import (
	"fmt"
	"io"
	"strings"

	"github.com/hurtki/school-events-bot/internal/domain"
	"github.com/xuri/excelize/v2"
)

func ParseXLSX(data io.Reader) (domain.Schedule, error) {
	f, err := excelize.OpenReader(data)
	if err != nil {
		return domain.Schedule{}, fmt.Errorf("can't start reading xlsx: %w", err)
	}

	levels := f.GetSheetList()
	events := []domain.Event{}
	for _, l := range levels {
		rows, err := f.GetRows(l)
		if err != nil {
			return domain.Schedule{}, fmt.Errorf("can't get rows for scpecific sheet %s: %w", l, err)
		}
		evs, err := parseSheetXLSX(rows, l)
		if err != nil {
			return domain.Schedule{}, fmt.Errorf("can't parse sheet %s: %w", l, err)
		}
		events = append(events, evs...)
	}
	return domain.NewSchedule(events)
}

func parseSheetXLSX(rows [][]string, groupName string) ([]domain.Event, error) {
	if len(rows) < 1 {
		return nil, nil
	}
	columnsCount := len(rows[0])

	events := make([]domain.Event, 0)
	for columnIndex := range columnsCount {
		var lastFoundDate *domain.Date = nil
		for rowIndex := range rows {
			if columnIndex >= len(rows[rowIndex]) {
				continue
			}
			content := rows[rowIndex][columnIndex]
			content = strings.TrimSpace(content)
			if content == "" {
				continue
			}
			date, err := domain.NewDate(content)
			if err == nil {
				lastFoundDate = &date
				continue
			}
			if lastFoundDate == nil {
				continue
			}
			event, err := domain.NewEvent(*lastFoundDate, groupName, content)
			if err != nil {
				continue
			}
			events = append(events, event)
		}
	}

	return events, nil
}

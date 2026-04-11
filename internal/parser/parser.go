package parser

import (
	"fmt"
	"io"
	"strings"

	"github.com/hurtki/school-events-bot/internal/domain"
	"github.com/xuri/excelize/v2"
)

var (
	sheetToGroup = map[string]domain.Group{
		"שכבת י ":  domain.TenthGradeGroup,
		`שכבת י"א`: domain.EleventhGradeGroup,
		`שכבת י"ב`: domain.TwelfthGradeGroup,
		`מכללה`:    domain.CollegeGroup,
	}
)

type Parser struct {
	f *excelize.File

	DocID string
}

func NewParser(data io.Reader, docID string) (Parser, error) {
	f, err := excelize.OpenReader(data)
	if err != nil {
		return Parser{}, fmt.Errorf("can't start reading xlsx: %w", err)
	}

	return Parser{
		f:     f,
		DocID: docID,
	}, nil
}

func (p *Parser) ParseXLSX() (domain.Schedule, error) {
	levels := p.f.GetSheetList()
	events := []domain.Event{}
	for _, l := range levels {
		evs, err := p.parseSheetXLSX(l)
		if err != nil {
			return domain.Schedule{}, fmt.Errorf("can't parse sheet %s: %w", l, err)
		}
		events = append(events, evs...)
	}
	return domain.NewSchedule(events)
}

func (p *Parser) parseSheetXLSX(sheetName string) ([]domain.Event, error) {
	rows, _ := p.f.GetRows(sheetName)
	group := sheetToGroup[sheetName]
	if len(rows) < 1 {
		return nil, nil
	}
	columnsCount := len(rows[0])

	events := make([]domain.Event, 0)
	for columnIndex := range columnsCount {
		var lastFoundDate *domain.Date = nil
		var dayStartcellAddr string

		dayBuffer := ""
		for rowIndex := range rows {
			if columnIndex >= len(rows[rowIndex]) {
				continue
			}
			content := rows[rowIndex][columnIndex]
			content = strings.TrimSpace(content)

			date, err := domain.NewDate(content)
			if err == nil {
				if lastFoundDate == nil {
					dayStartcellAddr, _ = excelize.CoordinatesToCellName(columnIndex+1, rowIndex+1)
					dayBuffer = ""
					lastFoundDate = &date
					continue
				}

				// parse everything from last day into events
				daySrcURL := getSourceURL(dayStartcellAddr, group, p.DocID)
				evs := parseDayIntoEvents(dayBuffer, group, *lastFoundDate, daySrcURL)
				events = append(events, evs...)

				dayBuffer = ""
				dayStartcellAddr, _ = excelize.CoordinatesToCellName(columnIndex+1, rowIndex+1)
				lastFoundDate = &date
				continue
			}

			if lastFoundDate == nil {
				continue
			}

			dayBuffer += "\n" + content
		}
	}

	return events, nil
}

var (
	GroupGIDs = map[domain.Group]string{
		domain.TenthGradeGroup:    "1035022939",
		domain.EleventhGradeGroup: "336153840",
		domain.TwelfthGradeGroup:  "1710319946",
		domain.CollegeGroup:       "898691425",
	}
)

func getSourceURL(cellAddr string, group domain.Group, docID string) string {
	return fmt.Sprintf("%s%s/edit#gid=%s&range=%s",
		baseSpreadsheetURL,
		docID,
		GroupGIDs[group],
		cellAddr,
	)
}

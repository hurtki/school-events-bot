package parser

import (
	"testing"

	"github.com/hurtki/school-events-bot/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestParseDayIntoEvents(t *testing.T) {
	cases := []struct {
		dayInput       string
		groupInput     domain.Group
		dateInput      domain.Date
		daySrcURLInput string

		expectedEvents []domain.Event
	}{
		{
			dayInput: `
מגן פיזיקה
			`,
			groupInput:     domain.TwelfthGradeGroup,
			daySrcURLInput: "test_url",

			expectedEvents: []domain.Event{
				{
					Date: domain.Date{},

					Type:  domain.ProtectionBagrutTestEvent,
					Group: domain.TwelfthGradeGroup,
					Text:  "מגן פיזיקה\n",

					SourceURL: "test_url",
				},
			},
		},
		{
			dayInput: `
מגן פיזיקה

מגן פיזיקה
			`,
			groupInput:     domain.TwelfthGradeGroup,
			daySrcURLInput: "test_url",

			expectedEvents: []domain.Event{
				{
					Date: domain.Date{},

					Type:  domain.ProtectionBagrutTestEvent,
					Group: domain.TwelfthGradeGroup,
					Text:  "מגן פיזיקה\n",

					SourceURL: "test_url",
				},
				{
					Date: domain.Date{},

					Type:  domain.ProtectionBagrutTestEvent,
					Group: domain.TwelfthGradeGroup,
					Text:  "מגן פיזיקה\n",

					SourceURL: "test_url",
				},
			},
		},
		{
			dayInput: `
מגן פיזיקה
מגן 9:00
			`,
			groupInput:     domain.TwelfthGradeGroup,
			daySrcURLInput: "test_url",

			expectedEvents: []domain.Event{
				{
					Date: domain.Date{},

					Type:  domain.ProtectionBagrutTestEvent,
					Group: domain.TwelfthGradeGroup,
					Text:  "מגן פיזיקה\nמגן 9:00\n",

					SourceURL: "test_url",
				},
			},
		},
		{
			dayInput: `
מגן פיזיקה
9:00
			`,
			groupInput:     domain.TwelfthGradeGroup,
			daySrcURLInput: "test_url",

			expectedEvents: []domain.Event{
				{
					Date: domain.Date{},

					Type:  domain.ProtectionBagrutTestEvent,
					Group: domain.TwelfthGradeGroup,
					Text:  "מגן פיזיקה\n9:00\n",

					SourceURL: "test_url",
				},
			},
		},
		{
			dayInput: `
מבחן
history
			`,
			groupInput:     domain.TwelfthGradeGroup,
			daySrcURLInput: "test_url",

			expectedEvents: []domain.Event{
				{
					Date: domain.Date{},

					Type:  domain.ExamEvent,
					Group: domain.TwelfthGradeGroup,
					Text:  "מבחן\nhistory\n",

					SourceURL: "test_url",
				},
			},
		},
	}

	// מבחן
	for _, c := range cases {
		evs := parseDayIntoEvents(
			c.dayInput,
			c.groupInput,
			c.dateInput,
			c.daySrcURLInput,
		)
		require.Equal(t, c.expectedEvents, evs)
	}
}

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
	}

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

package bbk

import (
	"fmt"
	"time"

	"github.com/FaintLocket424/rc-timing-api/internal/models"
)

func (suite *BBKTestSuite) TestParseRaceResult() {
	tests := []struct {
		name        string
		timestamp   time.Time
		url         string
		expectValue *models.ResultScrape
	}{
		{
			name:      "Initial Testing",
			timestamp: time.Date(2026, time.April, 11, 0, 0, 0, 0, time.UTC),
			url:       "forcc.co.uk/live",
			expectValue: &models.ResultScrape{
				HeatNumber:  ptr(4),
				ClassName:   ptr("2 Wheel Drive"),
				Round:       ptr(2),
				ElapsedTime: ptr(time.Duration(4)*time.Minute + time.Duration(37)*time.Second),
				Drivers: []models.DriverRaceResult{
					{
						CarNumber:       ptr(2),
						Name:            ptr("Stuart Clifford"),
						Laps:            ptr(23),
						Time:            ptr(time.Duration(4)*time.Minute + time.Duration(37_088)*time.Millisecond),
						BestLapDuration: ptr(time.Duration(10_936) * time.Millisecond),
					},
					{
						CarNumber:       ptr(1),
						Name:            ptr("Chris Goldsmith"),
						Laps:            ptr(22),
						Time:            ptr(time.Duration(4)*time.Minute + time.Duration(27_041)*time.Millisecond),
						BestLapDuration: ptr(time.Duration(10_945) * time.Millisecond),
					},
					{
						CarNumber:       ptr(5),
						Name:            ptr("David Armstrong"),
						Laps:            ptr(20),
						Time:            ptr(time.Duration(4)*time.Minute + time.Duration(26_546)*time.Millisecond),
						BestLapDuration: ptr(time.Duration(11_653) * time.Millisecond),
					},
					{
						CarNumber:       ptr(3),
						Name:            ptr("Jodie Moon"),
						Laps:            ptr(20),
						Time:            ptr(time.Duration(4)*time.Minute + time.Duration(29_788)*time.Millisecond),
						BestLapDuration: ptr(time.Duration(11_723) * time.Millisecond),
					},
					{
						CarNumber:       ptr(6),
						Name:            ptr("Mark Mummery"),
						Laps:            ptr(19),
						Time:            ptr(time.Duration(4)*time.Minute + time.Duration(36_804)*time.Millisecond),
						BestLapDuration: ptr(time.Duration(11_471) * time.Millisecond),
					},
					{
						CarNumber:       ptr(7),
						Name:            ptr("Graeme Wood"),
						Laps:            ptr(17),
						Time:            ptr(time.Duration(4)*time.Minute + time.Duration(25_767)*time.Millisecond),
						BestLapDuration: ptr(time.Duration(12_457) * time.Millisecond),
					},
					{
						CarNumber: ptr(4),
						Name:      ptr("Dave Mackins"),
					},
				},
				BestLap: &models.BestLap{
					DriverName: ptr("Stuart Clifford"),
					Time:       ptr(time.Duration(10_936) * time.Millisecond),
					LapNumber:  ptr(9),
				},
				ClassFT: &models.ClassFT{
					DriverName:     ptr("Kyle Moon"),
					Laps:           ptr(28),
					Time:           ptr(time.Duration(5)*time.Minute + time.Duration(5_566)*time.Millisecond),
					AvgLapDuration: ptr(time.Duration(10_913) * time.Millisecond),
					Round:          ptr(1),
				},
				ClassBestLap: &models.ClassBestLap{
					DriverName: ptr("James Procter"),
					Time:       ptr(time.Duration(10_192) * time.Millisecond),
				},
			},
		},
		{
			name:      "Potteries National 25th April 2026 15:22",
			timestamp: time.Date(2026, time.April, 25, 15, 22, 0, 0, time.UTC),
			url:       "www.brca-results.org/sections/10or/LiveResults",
			expectValue: &models.ResultScrape{
				HeatNumber:  ptr(4),
				ClassName:   ptr("2 Wheel Drive"),
				Round:       ptr(4),
				ElapsedTime: nil,
				Drivers: []models.DriverRaceResult{
					{
						CarNumber:       ptr(1),
						Name:            ptr("Scott Whyman"),
						Laps:            ptr(11),
						Time:            ptr(time.Duration(5)*time.Minute + time.Duration(4_230)*time.Millisecond),
						BestLapDuration: ptr(time.Duration(26_800) * time.Millisecond),
						BestLapNumber:   ptr(9),
					},
					{
						CarNumber:       ptr(4),
						Name:            ptr("Chris Boden"),
						Laps:            ptr(11),
						Time:            ptr(time.Duration(5)*time.Minute + time.Duration(6_160)*time.Millisecond),
						BestLapDuration: ptr(time.Duration(27_020) * time.Millisecond),
						BestLapNumber:   ptr(4),
					},
					{
						CarNumber:       ptr(6),
						Name:            ptr("Mark Christopher"),
						Laps:            ptr(11),
						Time:            ptr(time.Duration(5)*time.Minute + time.Duration(7_820)*time.Millisecond),
						BestLapDuration: ptr(time.Duration(27_370) * time.Millisecond),
						BestLapNumber:   ptr(4),
					},
					{
						CarNumber:       ptr(3),
						Name:            ptr("Andrew Shillito"),
						Laps:            ptr(11),
						Time:            ptr(time.Duration(5)*time.Minute + time.Duration(7_990)*time.Millisecond),
						BestLapDuration: ptr(time.Duration(27_430) * time.Millisecond),
						BestLapNumber:   ptr(10),
					},
					{
						CarNumber:       ptr(2),
						Name:            ptr("Chris Elworthy"),
						Laps:            ptr(11),
						Time:            ptr(time.Duration(5)*time.Minute + time.Duration(8_140)*time.Millisecond),
						BestLapDuration: ptr(time.Duration(26_410) * time.Millisecond),
						BestLapNumber:   ptr(11),
					},
					{
						CarNumber:       ptr(7),
						Name:            ptr("Jamie Banks"),
						Laps:            ptr(11),
						Time:            ptr(time.Duration(5)*time.Minute + time.Duration(10_940)*time.Millisecond),
						BestLapDuration: ptr(time.Duration(27_870) * time.Millisecond),
						BestLapNumber:   ptr(3),
					},
					{
						CarNumber:       ptr(10),
						Name:            ptr("Lee Stokes"),
						Laps:            ptr(11),
						Time:            ptr(time.Duration(5)*time.Minute + time.Duration(12_780)*time.Millisecond),
						BestLapDuration: ptr(time.Duration(27_090) * time.Millisecond),
						BestLapNumber:   ptr(4),
					},
					{
						CarNumber:       ptr(5),
						Name:            ptr("Roger Mills"),
						Laps:            ptr(11),
						Time:            ptr(time.Duration(5)*time.Minute + time.Duration(18_310)*time.Millisecond),
						BestLapDuration: ptr(time.Duration(26_780) * time.Millisecond),
						BestLapNumber:   ptr(10),
					},
					{
						CarNumber:       ptr(8),
						Name:            ptr("Colin Mulligan"),
						Laps:            ptr(11),
						Time:            ptr(time.Duration(5)*time.Minute + time.Duration(23_630)*time.Millisecond),
						BestLapDuration: ptr(time.Duration(28_020) * time.Millisecond),
						BestLapNumber:   ptr(2),
					},
					{
						CarNumber:       ptr(9),
						Name:            ptr("Miklos Szabados"),
						Laps:            ptr(11),
						Time:            ptr(time.Duration(5)*time.Minute + time.Duration(26_030)*time.Millisecond),
						BestLapDuration: ptr(time.Duration(28_430) * time.Millisecond),
						BestLapNumber:   ptr(3),
					},
				},
				BestLap: &models.BestLap{
					DriverName: ptr("Chris Elworthy"),
					Time:       ptr(time.Duration(26_410) * time.Millisecond),
					LapNumber:  ptr(11),
				},
				ClassFT: &models.ClassFT{
					DriverName:     ptr("Josh Holdsworth"),
					Laps:           ptr(13),
					Time:           ptr(time.Duration(5)*time.Minute + time.Duration(17_300)*time.Millisecond),
					AvgLapDuration: ptr(time.Duration(24_410) * time.Millisecond),
				},
				ClassBestLap: &models.ClassBestLap{
					DriverName: ptr("Phil Campbell"),
					Time:       ptr(time.Duration(22_510) * time.Millisecond),
				},
			},
		},
		{
			name:      "Potteries National 25th April 2026 16:49",
			timestamp: time.Date(2026, time.April, 25, 16, 49, 0, 0, time.UTC),
			url:       "www.brca-results.org/sections/10or/LiveResults",
			expectValue: &models.ResultScrape{
				FinalNumber: ptr(12),
				ClassName:   ptr("2 Wheel Drive"),
				Round:       ptr(1),
				ElapsedTime: ptr(2*time.Minute + 9*time.Second),
				Drivers: []models.DriverRaceResult{
					{
						CarNumber:       ptr(1),
						Name:            ptr("Colin Mulligan"),
						Laps:            ptr(4),
						Time:            ptr(time.Duration(1)*time.Minute + time.Duration(57_170)*time.Millisecond),
						LastLapDuration: ptr(time.Duration(28_230) * time.Millisecond),
					},
					{
						CarNumber:       ptr(4),
						Name:            ptr("Stuart Robinson"),
						Laps:            ptr(4),
						Time:            ptr(time.Duration(1)*time.Minute + time.Duration(57_170+700)*time.Millisecond),
						LastLapDuration: ptr(time.Duration(28_580) * time.Millisecond),
					},
					{
						CarNumber:       ptr(5),
						Name:            ptr("Trevor Follington"),
						Laps:            ptr(4),
						Time:            ptr(time.Duration(1)*time.Minute + time.Duration(57_170+6_780)*time.Millisecond),
						LastLapDuration: ptr(time.Duration(29_110) * time.Millisecond),
					},
					{
						CarNumber:       ptr(7),
						Name:            ptr("Nathan Stanmore"),
						Laps:            ptr(4),
						Time:            ptr(time.Duration(1)*time.Minute + time.Duration(57_170+7_940)*time.Millisecond),
						LastLapDuration: ptr(time.Duration(30_470) * time.Millisecond),
					},
					{
						CarNumber:       ptr(3),
						Name:            ptr("Dominic Mawby"),
						Laps:            ptr(4),
						Time:            ptr(time.Duration(1)*time.Minute + time.Duration(57_170+8_270)*time.Millisecond),
						LastLapDuration: ptr(time.Duration(30_360) * time.Millisecond),
					},
					{
						CarNumber:       ptr(6),
						Name:            ptr("Jamie Careless"),
						Laps:            ptr(4),
						Time:            ptr(time.Duration(1)*time.Minute + time.Duration(57_170+9_320)*time.Millisecond),
						LastLapDuration: ptr(time.Duration(30_680) * time.Millisecond),
					},
					{
						CarNumber:       ptr(8),
						Name:            ptr("Paul Steadman"),
						Laps:            ptr(4),
						Time:            ptr(time.Duration(1)*time.Minute + time.Duration(57_170+10_550)*time.Millisecond),
						LastLapDuration: ptr(time.Duration(29_990) * time.Millisecond),
					},
					{
						CarNumber:       ptr(2),
						Name:            ptr("Miklos Szabados"),
						Laps:            ptr(3),
						Time:            ptr(time.Duration(1)*time.Minute + time.Duration(36_290)*time.Millisecond),
						LastLapDuration: ptr(time.Duration(29_250) * time.Millisecond),
					},
				},
				BestLap: &models.BestLap{
					DriverName: ptr("Colin Mulligan"),
					Time:       ptr(time.Duration(28_130) * time.Millisecond),
					LapNumber:  ptr(2),
				},
			},
		},
	}

	for _, tc := range tests {
		suite.Run(tc.name, func() {
			scraper := &BBKScraper{
				Target: fmt.Sprintf("%s/%s/%s", suite.server.URL, tc.timestamp.Format("2006-01-02_15-04-05"), tc.url),
				Client: suite.server.Client(),
			}

			data, err := scraper.GetLiveTiming()

			suite.Require().NoError(err)
			suite.Require().NotNil(data)

			suite.Equal(tc.expectValue, data)
		})
	}
}

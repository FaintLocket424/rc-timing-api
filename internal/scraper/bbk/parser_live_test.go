package bbk

import (
	"fmt"
	"time"

	"github.com/FaintLocket424/rc-timing-api/internal/models"
)

func ptr[T any](v T) *T {
	return &v
}

func (suite *BBKTestSuite) TestParseRaceResult() {
	tests := []struct {
		name        string
		timestamp   time.Time
		url         string
		expectValue *models.LiveTimingScrape
	}{
		{
			name:      "Initial Testing",
			timestamp: time.Date(2026, time.April, 11, 0, 0, 0, 0, time.UTC),
			url:       "forcc.co.uk/live",
			expectValue: &models.LiveTimingScrape{
				HeatNumber:  4,
				ClassName:   "2 Wheel Drive",
				Round:       2,
				ElapsedTime: time.Duration(4)*time.Minute + time.Duration(37)*time.Second,
				Drivers: []models.DriverRaceResult{
					{
						CarNumber:       2,
						Name:            "Stuart Clifford",
						Laps:            23,
						Time:            time.Duration(4)*time.Minute + time.Duration(37_088)*time.Millisecond,
						BestLapDuration: time.Duration(10_936) * time.Millisecond,
					},
					{
						CarNumber:       1,
						Name:            "Chris Goldsmith",
						Laps:            22,
						Time:            time.Duration(4)*time.Minute + time.Duration(27_041)*time.Millisecond,
						BestLapDuration: time.Duration(10_945) * time.Millisecond,
					},
					{
						CarNumber:       5,
						Name:            "David Armstrong",
						Laps:            20,
						Time:            time.Duration(4)*time.Minute + time.Duration(26_546)*time.Millisecond,
						BestLapDuration: time.Duration(11_653) * time.Millisecond,
					},
					{
						CarNumber:       3,
						Name:            "Jodie Moon",
						Laps:            20,
						Time:            time.Duration(4)*time.Minute + time.Duration(29_788)*time.Millisecond,
						BestLapDuration: time.Duration(11_723) * time.Millisecond,
					},
					{
						CarNumber:       6,
						Name:            "Mark Mummery",
						Laps:            19,
						Time:            time.Duration(4)*time.Minute + time.Duration(36_804)*time.Millisecond,
						BestLapDuration: time.Duration(11_471) * time.Millisecond,
					},
					{
						CarNumber:       7,
						Name:            "Graeme Wood",
						Laps:            17,
						Time:            time.Duration(4)*time.Minute + time.Duration(25_767)*time.Millisecond,
						BestLapDuration: time.Duration(12_457) * time.Millisecond,
					},
					{
						CarNumber: 4,
						Name:      "Dave Mackins",
					},
				},
				BestLap: models.BestLap{
					DriverName: "Stuart Clifford", Time: time.Duration(10_936) * time.Millisecond, LapNumber: 9,
				},
				ClassFT: models.ClassFT{
					DriverName:     "Kyle Moon",
					Laps:           28,
					Time:           time.Duration(5)*time.Minute + time.Duration(5_566)*time.Millisecond,
					AvgLapDuration: time.Duration(10_913) * time.Millisecond,
					Round:          1,
				},
				ClassBestLap: models.ClassBestLap{
					DriverName: "James Procter",
					Time:       time.Duration(10_192) * time.Millisecond,
				},
			},
		},
		{
			name:      "Potteries National 25th April 2026 15:22",
			timestamp: time.Date(2026, time.April, 25, 15, 22, 0, 0, time.UTC),
			url:       "www.brca-results.org/sections/10or/LiveResults",
			expectValue: &models.LiveTimingScrape{
				HeatNumber:  4,
				ClassName:   "2 Wheel Drive",
				Round:       4,
				ElapsedTime: 0,
				Drivers: []models.DriverRaceResult{
					{
						CarNumber:       1,
						Name:            "Scott Whyman",
						Laps:            11,
						Time:            time.Duration(5)*time.Minute + time.Duration(4_230)*time.Millisecond,
						BestLapDuration: time.Duration(26_800) * time.Millisecond,
						BestLapNumber:   9,
					},
					{
						CarNumber:       4,
						Name:            "Chris Boden",
						Laps:            11,
						Time:            time.Duration(5)*time.Minute + time.Duration(6_160)*time.Millisecond,
						BestLapDuration: time.Duration(27_020) * time.Millisecond,
						BestLapNumber:   4,
					},
					{
						CarNumber:       6,
						Name:            "Mark Christopher",
						Laps:            11,
						Time:            time.Duration(5)*time.Minute + time.Duration(7_820)*time.Millisecond,
						BestLapDuration: time.Duration(27_370) * time.Millisecond,
						BestLapNumber:   4,
					},
					{
						CarNumber:       3,
						Name:            "Andrew Shillito",
						Laps:            11,
						Time:            time.Duration(5)*time.Minute + time.Duration(7_990)*time.Millisecond,
						BestLapDuration: time.Duration(27_430) * time.Millisecond,
						BestLapNumber:   10,
					},
					{
						CarNumber:       2,
						Name:            "Chris Elworthy",
						Laps:            11,
						Time:            time.Duration(5)*time.Minute + time.Duration(8_140)*time.Millisecond,
						BestLapDuration: time.Duration(26_410) * time.Millisecond,
						BestLapNumber:   11,
					},
					{
						CarNumber:       7,
						Name:            "Jamie Banks",
						Laps:            11,
						Time:            time.Duration(5)*time.Minute + time.Duration(10_940)*time.Millisecond,
						BestLapDuration: time.Duration(27_870) * time.Millisecond,
						BestLapNumber:   3,
					},
					{
						CarNumber:       10,
						Name:            "Lee Stokes",
						Laps:            11,
						Time:            time.Duration(5)*time.Minute + time.Duration(12_780)*time.Millisecond,
						BestLapDuration: time.Duration(27_090) * time.Millisecond,
						BestLapNumber:   4,
					},
					{
						CarNumber:       5,
						Name:            "Roger Mills",
						Laps:            11,
						Time:            time.Duration(5)*time.Minute + time.Duration(18_310)*time.Millisecond,
						BestLapDuration: time.Duration(26_780) * time.Millisecond,
						BestLapNumber:   10,
					},
					{
						CarNumber:       8,
						Name:            "Colin Mulligan",
						Laps:            11,
						Time:            time.Duration(5)*time.Minute + time.Duration(23_630)*time.Millisecond,
						BestLapDuration: time.Duration(28_020) * time.Millisecond,
						BestLapNumber:   2,
					},
					{
						CarNumber:       9,
						Name:            "Miklos Szabados",
						Laps:            11,
						Time:            time.Duration(5)*time.Minute + time.Duration(26_030)*time.Millisecond,
						BestLapDuration: time.Duration(28_430) * time.Millisecond,
						BestLapNumber:   3,
					},
				},
				BestLap: models.BestLap{
					DriverName: "Chris Elworthy", Time: time.Duration(26_410) * time.Millisecond, LapNumber: 11,
				},
				ClassFT: models.ClassFT{
					DriverName:     "Josh Holdsworth",
					Laps:           13,
					Time:           time.Duration(5)*time.Minute + time.Duration(17_300)*time.Millisecond,
					AvgLapDuration: time.Duration(24_410) * time.Millisecond,
					Round:          0,
				},
				ClassBestLap: models.ClassBestLap{
					DriverName: "Phil Campbell",
					Time:       time.Duration(22_510) * time.Millisecond,
				},
			},
		},
		{
			name:      "Potteries National 25th April 2026 16:49",
			timestamp: time.Date(2026, time.April, 25, 16, 49, 0, 0, time.UTC),
			url:       "www.brca-results.org/sections/10or/LiveResults",
			expectValue: &models.LiveTimingScrape{
				ClassName:   "2 Wheel Drive",
				Round:       1,
				ElapsedTime: 2*time.Minute + 9*time.Second,
				Drivers: []models.DriverRaceResult{
					{
						CarNumber:       1,
						Name:            "Colin Mulligan",
						Laps:            4,
						Time:            time.Duration(1)*time.Minute + time.Duration(57_170)*time.Millisecond,
						LastLapDuration: time.Duration(28_230) * time.Millisecond,
					},
					{
						CarNumber:       4,
						Name:            "Stuart Robinson",
						Laps:            4,
						Time:            time.Duration(1)*time.Minute + time.Duration(57_170+700)*time.Millisecond,
						LastLapDuration: time.Duration(28_580) * time.Millisecond,
					},
					{
						CarNumber:       5,
						Name:            "Trevor Follington",
						Laps:            4,
						Time:            time.Duration(1)*time.Minute + time.Duration(57_170+6_780)*time.Millisecond,
						LastLapDuration: time.Duration(29_110) * time.Millisecond,
					},
					{
						CarNumber:       7,
						Name:            "Nathan Stanmore",
						Laps:            4,
						Time:            time.Duration(1)*time.Minute + time.Duration(57_170+7_940)*time.Millisecond,
						LastLapDuration: time.Duration(30_470) * time.Millisecond,
					},
					{
						CarNumber:       3,
						Name:            "Dominic Mawby",
						Laps:            4,
						Time:            time.Duration(1)*time.Minute + time.Duration(57_170+8_270)*time.Millisecond,
						LastLapDuration: time.Duration(30_360) * time.Millisecond,
					},
					{
						CarNumber:       6,
						Name:            "Jamie Careless",
						Laps:            4,
						Time:            time.Duration(1)*time.Minute + time.Duration(57_170+9_320)*time.Millisecond,
						LastLapDuration: time.Duration(30_680) * time.Millisecond,
					},
					{
						CarNumber:       8,
						Name:            "Paul Steadman",
						Laps:            4,
						Time:            time.Duration(1)*time.Minute + time.Duration(57_170+10_550)*time.Millisecond,
						LastLapDuration: time.Duration(29_990) * time.Millisecond,
					},
					{
						CarNumber:       2,
						Name:            "Miklos Szabados",
						Laps:            3,
						Time:            time.Duration(1)*time.Minute + time.Duration(36_290)*time.Millisecond,
						LastLapDuration: time.Duration(29_250) * time.Millisecond,
					},
				},
				BestLap: models.BestLap{
					DriverName: "Colin Mulligan", Time: time.Duration(28_130) * time.Millisecond, LapNumber: 2,
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

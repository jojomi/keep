package keep

import (
	"math/rand"
	"testing"
	"time"
)

func TestListForDate(t *testing.T) {
	testTime, _ := time.Parse(time.RFC3339, "2018-12-12T11:45:26.371Z")

	shuffledHours := getTimesEveryHourFrom(testTime, 300)
	rand.Shuffle(len(shuffledHours), func(i, j int) { shuffledHours[i], shuffledHours[j] = shuffledHours[j], shuffledHours[i] })

	type args struct {
		input []TimeResource
		reqs  Requirements
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantLast  string
		wantErr   bool
	}{
		{
			name: "hours only",
			args: args{
				input: getTimesEveryHourFrom(testTime, 10),
				reqs: Requirements{
					Hours: 3,
				},
			},
			wantCount: 3,
			wantLast:  "2018-12-12T08:45:26.371Z",
			wantErr:   false,
		}, {
			name: "cutoff",
			args: args{
				input: getTimesEveryHourFrom(testTime, 2),
				reqs: Requirements{
					Hours: 30,
				},
			},
			wantCount: 2,
			wantLast:  "2018-12-12T09:45:26.371Z",
			wantErr:   false,
		}, {
			name: "days only",
			args: args{
				input: getTimesEveryHourFrom(testTime, 50),
				reqs: Requirements{
					Days: 3,
				},
			},
			wantCount: 3,
			wantLast:  "2018-12-10T10:45:26.371Z",
			wantErr:   false,
		}, {
			name: "months only",
			args: args{
				input: getTimesEveryDayFrom(testTime, 50),
				reqs: Requirements{
					Months: 2,
				},
			},
			wantCount: 2,
			wantLast:  "2018-11-11T11:45:26.371Z",
			wantErr:   false,
		}, {
			name: "years only",
			args: args{
				input: getTimesEveryDayFrom(testTime, 1000),
				reqs: Requirements{
					Years: 2,
				},
			},
			wantCount: 2,
			wantLast:  "2017-12-11T11:45:26.371Z",
			wantErr:   false,
		}, {
			name: "stretch hours",
			args: args{
				input: getTimesEveryDayFrom(testTime, 10),
				reqs: Requirements{
					Hours: 6,
				},
			},
			wantCount: 6,
			wantLast:  "2018-12-06T11:45:26.371Z",
			wantErr:   false,
		}, {
			name: "mixed",
			args: args{
				input: getTimesEveryHourFrom(testTime, 300),
				reqs: Requirements{
					Hours:  5,
					Days:   3,
					Weeks:  2,
					Months: 6,
					Years:  4,
				},
			},
			wantCount: 10,
			wantLast:  "2018-12-02T05:45:26.371Z",
			wantErr:   false,
		}, {
			name: "shuffled",
			args: args{
				input: shuffledHours,
				reqs: Requirements{
					Hours:  5,
					Days:   3,
					Weeks:  2,
					Months: 6,
					Years:  4,
				},
			},
			wantCount: 10,
			wantLast:  "2018-12-02T05:45:26.371Z",
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListForDate(testTime, tt.args.input, tt.args.reqs)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListForDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantCount {
				t.Errorf("ListForDate() count = %v, wantCount %v", len(got), tt.wantCount)
				return
			}
			last := got[len(got)-1]
			wantLast, err := time.Parse(time.RFC3339, tt.wantLast)
			if err != nil {
				panic(err)
			}
			if !wantLast.Equal(last.GetTime()) {
				t.Errorf("ListForDate() error = %v, wantLast %v", last.GetTime(), tt.wantLast)
				return
			}
		})
	}
}

func TestListRemovableForDate(t *testing.T) {
	testTime, _ := time.Parse(time.RFC3339, "2018-12-12T11:45:26.371Z")

	type args struct {
		input []TimeResource
		reqs  Requirements
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		{
			name: "mixed",
			args: args{
				input: getTimesEveryHourFrom(testTime, 50),
				reqs: Requirements{
					Hours: 3,
					Days:  2,
				},
			},
			wantCount: 50 - 3 - 2,
			wantErr:   false,
		},
		{
			name: "all used",
			args: args{
				input: getTimesEveryDayFrom(testTime, 5),
				reqs: Requirements{
					Days: 20,
				},
			},
			wantCount: 0,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListRemovableForDate(testTime, tt.args.input, tt.args.reqs)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListForDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantCount {
				t.Errorf("ListForDate() count = %v, wantCount %v", len(got), tt.wantCount)
				return
			}
		})
	}
}

type TestTimeResource struct {
	t time.Time
}

func (x TestTimeResource) GetTime() time.Time {
	return x.t
}

func getTimesEveryHourFrom(referenceTime time.Time, hourCount int) []TimeResource {
	result := make([]TimeResource, hourCount)
	hour, _ := time.ParseDuration("-1h")
	for i := 0; i < hourCount; i++ {
		referenceTime = referenceTime.Add(hour)
		result[i] = TestTimeResource{
			t: referenceTime,
		}
	}
	return result
}

func getTimesEveryDayFrom(referenceTime time.Time, dayCount int) []TimeResource {
	result := make([]TimeResource, dayCount)
	for i := 0; i < dayCount; i++ {
		referenceTime = referenceTime.AddDate(0, 0, -1)
		result[i] = TestTimeResource{
			t: referenceTime,
		}
	}
	return result
}

package keep

import (
	"sort"
	"testing"
	"time"
)

func TestJailhouse_KeptElements(t *testing.T) {
	testDate := time.Date(2024, time.January, 20, 0, 0, 0, 0, time.UTC)

	type fields struct {
		elements     []TimeResource
		requirements *Requirements
	}
	tests := []struct {
		name   string
		fields fields
		want   []*JailhouseTimeResource
	}{
		{
			name: "last 1 (unordered)", // requires sorting of elements by time and older ones
			fields: fields{
				elements: []TimeResource{
					date("2023-03-09"),
					date("2023-01-01"),
					date("2023-02-22"),
				},
				requirements: NewRequirements().Add(LAST, 1),
			},
			want: []*JailhouseTimeResource{
				NewJailhouseTimeResource(date("2023-03-09")).AddTag(TimeRangeTagFrom(LAST, 1)),
			},
		},
		{
			name: "last n", // selects more than one
			fields: fields{
				elements: []TimeResource{
					date("2023-01-01"),
					date("2023-02-22"),
					date("2023-03-09"),
				},
				requirements: NewRequirements().Add(LAST, 2),
			},
			want: []*JailhouseTimeResource{
				NewJailhouseTimeResource(date("2023-02-22")).AddTag(TimeRangeTagFrom(LAST, 2)),
				NewJailhouseTimeResource(date("2023-03-09")).AddTag(TimeRangeTagFrom(LAST, 1)),
			},
		},
		{
			name: "last 1 and day 1", // combines two levels
			fields: fields{
				elements: []TimeResource{
					date("2023-01-01"),
					date("2023-02-22"),
					date("2023-03-09"),
				},
				requirements: NewRequirements().Add(LAST, 1).Add(DAY, 1),
			},
			want: []*JailhouseTimeResource{
				NewJailhouseTimeResource(date("2023-03-09")).AddTag(TimeRangeTagFrom(LAST, 1)),
				NewJailhouseTimeResource(date("2023-02-22")).AddTag(TimeRangeTagFrom(DAY, 1)),
			},
		},
		{
			name: "year 2", // skips in-between elements
			fields: fields{
				elements: []TimeResource{
					date("2022-05-10"), // slightly too old for year 1
					date("2022-05-01"),
					date("2023-02-22"),
					date("2023-05-09"),
				},
				requirements: NewRequirements().Add(YEAR, 2),
			},
			want: []*JailhouseTimeResource{
				NewJailhouseTimeResource(date("2023-05-09")).AddTag(TimeRangeTagFrom(YEAR, 1)),
				NewJailhouseTimeResource(date("2022-05-01")).AddTag(TimeRangeTagFrom(YEAR, 2)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := NewDefaultJailhouse()
			x.AddElements(tt.fields.elements...)
			x.ApplyRequirementsForDate(*tt.fields.requirements, testDate)
			assertSameElements(t, tt.want, x.KeptElements())
		})
	}
}

func assertSameElements(t *testing.T, expected, seen []*JailhouseTimeResource) {
	// sort both lists
	sort.SliceStable(expected, func(i, j int) bool {
		return expected[i].TimeResource.GetTime().After(expected[j].TimeResource.GetTime())
	})
	sort.SliceStable(seen, func(i, j int) bool {
		return seen[i].TimeResource.GetTime().After(seen[j].TimeResource.GetTime())
	})

	if len(expected) != len(seen) {
		t.Errorf("different set sizes, expected %d entries, got %d, got %v", len(expected), len(seen), seen)
		t.FailNow()
	}

	// now loop and compare
	var (
		s            *JailhouseTimeResource
		sTags, eTags []TimeRangeTag
		sTag         TimeRangeTag
	)
	for i, e := range expected {
		s = seen[i]

		if !s.GetTime().Equal(e.GetTime()) {
			t.Errorf("times did not match (elem #%d), expected %v, got %v", i, e.GetTime(), s.GetTime())
			t.FailNow()
		}

		eTags = e.GetTags()
		sTags = s.GetTags()
		sort.SliceStable(eTags, func(i, j int) bool {
			return int(eTags[i].TimeRange) < int(eTags[j].TimeRange)
		})
		sort.SliceStable(sTags, func(i, j int) bool {
			return int(sTags[i].TimeRange) < int(sTags[j].TimeRange)
		})

		for i, eTag := range eTags {
			sTag = sTags[i]

			if eTag.TimeRange != sTag.TimeRange {
				t.Errorf("timerange tag did not match")
				t.FailNow()
			}
			if eTag.Index != sTag.Index {
				t.Errorf("timerange index did not match")
				t.FailNow()
			}
		}
	}
}

func date(date string) *TestTimeResource {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		panic(err)
	}
	return &TestTimeResource{
		t: t,
	}
}

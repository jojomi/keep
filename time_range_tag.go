package keep

import "fmt"

func TimeRangeTagFrom(timeRange TimeRange, index uint16) TimeRangeTag {
	return TimeRangeTag{
		TimeRange: timeRange,
		Index:     index,
	}
}

type TimeRangeTag struct {
	TimeRange TimeRange
	Index     uint16
}

func (x TimeRangeTag) String() string {
	return fmt.Sprintf("%s-%d", x.TimeRange.String(), x.Index)
}

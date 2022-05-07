package keep

import (
	"fmt"
	"strings"
)

// Requirements define which elements in an input slice should be kept
type Requirements struct {
	ranges map[TimeRange]uint16
}

// NewRequirements creates a new empty Requirement definition.
func NewRequirements() *Requirements {
	return &Requirements{
		ranges: make(map[TimeRange]uint16, 0),
	}
}

// NewRequirementsFromMap makes a Requirement from a map with TimeRange as keys and their number as values.
func NewRequirementsFromMap(data map[TimeRange]uint16) *Requirements {
	r := NewRequirements()
	for key, value := range data {
		r.ranges[key] = value
	}
	return r
}

// IsEmpty is true iff no files should be kept.
func (x Requirements) IsEmpty() bool {
	for _, v := range x.ranges {
		if v > 0 {
			return false
		}
	}
	return true
}

// Get returns the number of elements in this Requirement for a given TimeRange.
func (x Requirements) Get(timeRange TimeRange) uint16 {
	if v, ok := x.ranges[timeRange]; ok {
		return v
	}
	return 0
}

// Add adds a number of required elements for a given TimeRange.
func (x *Requirements) Add(timeRange TimeRange, value int8) *Requirements {
	if _, ok := x.ranges[timeRange]; !ok {
		x.ranges[timeRange] = 0
	}
	x.ranges[timeRange] = uint16(int(x.ranges[timeRange]) + int(value))

	return x
}

// DeepCopy returns a Requirement copy with the same properties.
func (x Requirements) DeepCopy() Requirements {
	r := NewRequirements()
	for key, value := range x.ranges {
		r.ranges[key] = value
	}
	return *r
}

// String prints a Requirement configuration
func (x Requirements) String() string {
	var elems []string
	for _, r := range TimeRangeNames() {
		if v, ok := x.ranges[MustParseTimeRange(r)]; ok {
			elems = append(elems, fmt.Sprintf("%s=%d", r, v))
		}
	}
	return strings.Join(elems, ", ")
}

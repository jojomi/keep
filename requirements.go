package keep

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Requirements define which elements in an input slice should be kept
type Requirements struct {
	ranges map[TimeRange]uint16
}

// NewRequirements creates a new empty Requirement definition.
func NewRequirements() *Requirements {
	return &Requirements{
		ranges: make(map[TimeRange]uint16),
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

func NewRequirementsFromString(source string) *Requirements {
	r := NewRequirements()
	re := regexp.MustCompile(`(?i)(\d+)\s+(last|seconds?|minutes?|hours?|days?|weeks?|months?|quarters?|years?)`)
	matches := re.FindAllStringSubmatch(source, -1)
	for _, match := range matches {
		num, err := strconv.Atoi(match[1])
		if err != nil {
			continue
		}
		switch strings.ToLower(match[2]) {
		case "last":
			r.ranges[LAST] += uint16(num)
		case "second", "seconds":
			r.ranges[SECOND] += uint16(num)
		case "minute", "minutes":
			r.ranges[MINUTE] += uint16(num)
		case "hour", "hours":
			r.ranges[HOUR] += uint16(num)
		case "day", "days":
			r.ranges[DAY] += uint16(num)
		case "week", "weeks":
			r.ranges[WEEK] += uint16(num)
		case "month", "months":
			r.ranges[MONTH] += uint16(num)
		case "quarter", "quarters":
			r.ranges[QUARTER] += uint16(num)
		case "year", "years":
			r.ranges[YEAR] += uint16(num)
		}
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

package keep

import (
	"time"

	"github.com/juju/errors"
	"golang.org/x/exp/slices"
)

type Jailhouse struct {
	elements []*JailhouseTimeResource
	levels   []TimeRange
}

func NewDefaultJailhouse() *Jailhouse {
	jailhouse := &Jailhouse{
		elements: make([]*JailhouseTimeResource, 0),
		levels: []TimeRange{
			LAST, SECOND, MINUTE, HOUR, DAY, WEEK, MONTH, QUARTER, YEAR, DECADE, CENTURY, MILLENIUM,
		},
	}
	return jailhouse
}

func (x *Jailhouse) GetLevels() []TimeRange {
	return x.levels
}

func (x *Jailhouse) AddElements(elems ...TimeResource) *Jailhouse {
	// add
	for _, e := range elems {
		x.elements = append(x.elements, NewJailhouseTimeResource(e))
	}

	// sort
	x.sortResources(x.elements)

	return x
}

func (x *Jailhouse) ApplyRequirements(reqs Requirements) *Jailhouse {
	return x.ApplyRequirementsForDate(reqs, time.Now())
}

func (x *Jailhouse) ApplyRequirementsForDate(reqs Requirements, referenceDate time.Time) *Jailhouse {
	// clear previous results
	for _, item := range x.elements {
		item.ClearTags()
	}

	// loop and keep or pass
	var (
		currentTime time.Time
		lastOfLevel bool
		nextTime    time.Time
	)

	for _, level := range x.GetLevels() {
		levelElementIndex := 1

		if reqs.Get(level) == 0 {
			continue
		}

		// reset time
		currentTime = referenceDate

		for _, item := range x.elements {
			// skip due to time?
			if item.GetTime().After(currentTime) {
				// this one is not in the output -> can be dropped
				continue
			}

			// continue?
			nextTime, lastOfLevel, reqs = x.nextTickForLevel(currentTime, reqs, level)

			// mark
			item.AddTag(TimeRangeTagFrom(level, uint16(levelElementIndex)))
			levelElementIndex++

			if lastOfLevel {
				break
			}

			currentTime = nextTime
		}
	}

	return x
}

func (x Jailhouse) FilteredElements(filter func(*JailhouseTimeResource) bool) []*JailhouseTimeResource {
	result := make([]*JailhouseTimeResource, 0)
	for _, element := range x.elements {
		if !filter(element) {
			continue
		}
		result = append(result, element)
	}
	return result
}

func (x Jailhouse) KeptElements() []*JailhouseTimeResource {
	return x.FilteredElements(func(element *JailhouseTimeResource) bool {
		return !element.IsFree()
	})
}

func (x Jailhouse) KeptElementsByLevel(level TimeRange) []*JailhouseTimeResource {
	return x.FilteredElements(func(element *JailhouseTimeResource) bool {
		return element.HasLevel(level)
	})
}

func (x Jailhouse) FreeElements() []*JailhouseTimeResource {
	return x.FilteredElements(func(element *JailhouseTimeResource) bool {
		return element.IsFree()
	})
}

func (x Jailhouse) Elements() []*JailhouseTimeResource {
	return x.elements
}

func (x Jailhouse) nextTickForLevel(current time.Time, requirements Requirements, level TimeRange) (newTime time.Time, lastOfType bool, newRequirements Requirements) {
	if requirements.Get(level) == 0 {
		return
	}

	switch level {
	case LAST:
		newTime = current
	case SECOND:
		d, _ := time.ParseDuration("-1s")
		newTime = current.Add(d)
	case MINUTE:
		d, _ := time.ParseDuration("-1m")
		newTime = current.Add(d)
	case HOUR:
		d, _ := time.ParseDuration("-1h")
		newTime = current.Add(d)
	case DAY:
		newTime = current.AddDate(0, 0, -1)
	case WEEK:
		newTime = current.AddDate(0, 0, -7)
	case MONTH:
		newTime = current.AddDate(0, -1, 0)
	case QUARTER:
		newTime = current.AddDate(0, -3, 0)
	case YEAR:
		newTime = current.AddDate(-1, 0, 0)
	case DECADE:
		newTime = current.AddDate(-10, 0, 0)
	case CENTURY:
		newTime = current.AddDate(-100, 0, 0)
	case MILLENIUM:
		newTime = current.AddDate(-1000, 0, 0)
	default:
		err := errors.Errorf("could not find level %s", level)
		panic(err)
	}

	newRequirements = requirements.DeepCopy()
	newRequirements.Add(level, -1)
	lastOfType = newRequirements.Get(level) <= 0
	return
}

func (x Jailhouse) sortResources(input []*JailhouseTimeResource) {
	slices.SortFunc(input, func(a, b *JailhouseTimeResource) int {
		if a.GetTime().Equal(b.GetTime()) {
			return 0
		}
		if a.GetTime().After(b.GetTime()) {
			return -1
		}
		return 1
	})
}

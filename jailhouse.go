package keep

import (
	"math"
	"time"

	"github.com/juju/errors"
	"golang.org/x/exp/slices"
)

type Jailhouse[T TimeResource] struct {
	elements []*JailhouseTimeResource[T]
	levels   []TimeRange
}

func NewDefaultJailhouse[T TimeResource]() *Jailhouse[T] {
	jailhouse := &Jailhouse[T]{
		elements: make([]*JailhouseTimeResource[T], 0),
		levels: []TimeRange{
			LAST, SECOND, MINUTE, HOUR, DAY, WEEK, MONTH, QUARTER, YEAR, DECADE, CENTURY, MILLENIUM,
		},
	}
	return jailhouse
}

func (x *Jailhouse[T]) GetLevels() []TimeRange {
	return x.levels
}

func (x *Jailhouse[T]) AddElements(elems ...T) *Jailhouse[T] {
	// add
	for _, e := range elems {
		x.elements = append(x.elements, NewJailhouseTimeResource[T](e))
	}

	// sort
	x.sortResources(x.elements)

	return x
}

func (x *Jailhouse[T]) ApplyRequirements(reqs Requirements) *Jailhouse[T] {
	return x.ApplyRequirementsForDate(reqs, time.Now())
}

func (x *Jailhouse[T]) ApplyRequirementsForDate(reqs Requirements, referenceDate time.Time) *Jailhouse[T] {
	// clear previous results
	for _, item := range x.elements {
		item.ClearTags()
	}

	// loop and keep or pass
	var (
		currentTime       = referenceDate
		lastOfLevel       bool
		nextTime          time.Time
		extendedTime      time.Time
		levelElementIndex int
		startElementIndex = 0
		item              *JailhouseTimeResource[T]
		nextItem          *JailhouseTimeResource[T]
		elementCount      = len(x.elements)
		levelStart        int
	)

	for _, level := range x.GetLevels() {
		levelElementIndex = 0

		if reqs.Get(level) == 0 {
			continue
		}

		// reset time
		// currentTime = referenceDate

		levelStart = startElementIndex
		for i := levelStart; i < elementCount; i++ {
			item = x.elements[i]

			// ignore the future
			if item.GetTime().After(referenceDate) {
				continue
			}

			// - first element for level is always kept
			// - for LAST we keep any element
			// - we do need at least one more and this one is the last? -> keep it
			// - skip because the next one is still "in (extended) range" and close to the target date?
			if level > LAST && i > levelStart && i < elementCount-1 {
				nextItem = x.elements[i+1]
				if !nextItem.GetTime().Before(extendedTime) {
					// select either this one or the next, depending on which is closer to the "current time" we aim for
					if math.Abs(float64(nextItem.GetTime().Sub(currentTime))) < math.Abs(float64(currentTime.Sub(item.GetTime()))) {
						// this one is not in the output -> can be dropped
						continue
					}
				}
			}

			// continue?
			nextTime, extendedTime, lastOfLevel, reqs = x.nextTickForLevel(item.GetTime(), reqs, level)

			// mark
			item.AddTag(TimeRangeTagFrom(level, uint16(levelElementIndex+1)))
			levelElementIndex++
			startElementIndex = i + 1

			if lastOfLevel {
				break
			}

			currentTime = nextTime
		}
	}

	// os.Exit(1)
	return x
}

func (x *Jailhouse[T]) FilteredElements(filter func(*JailhouseTimeResource[T]) bool) []*JailhouseTimeResource[T] {
	result := make([]*JailhouseTimeResource[T], 0)
	for _, element := range x.elements {
		if !filter(element) {
			continue
		}
		result = append(result, element)
	}
	return result
}

func (x *Jailhouse[T]) KeptElements() []*JailhouseTimeResource[T] {
	return x.FilteredElements(func(element *JailhouseTimeResource[T]) bool {
		return !element.IsFree()
	})
}

func (x *Jailhouse[T]) KeptElementsByLevel(level TimeRange) []*JailhouseTimeResource[T] {
	return x.FilteredElements(func(element *JailhouseTimeResource[T]) bool {
		return element.HasLevel(level)
	})
}

func (x *Jailhouse[T]) FreeElements() []*JailhouseTimeResource[T] {
	return x.FilteredElements(func(element *JailhouseTimeResource[T]) bool {
		return element.IsFree()
	})
}

func (x *Jailhouse[T]) Elements() []*JailhouseTimeResource[T] {
	return x.elements
}

func (x *Jailhouse[T]) nextTickForLevel(current time.Time, requirements Requirements, level TimeRange) (newTime, extendedTime time.Time, lastOfType bool, newRequirements Requirements) {
	if requirements.Get(level) == 0 {
		return
	}

	newTime = x.addLevelStep(level, current)
	extendedTime = newTime
	if level >= MINUTE {
		extendedTime = x.addLevelStep(TimeRange(int(level-1)), newTime)
	}

	newRequirements = requirements.DeepCopy()
	newRequirements.Add(level, -1)
	lastOfType = newRequirements.Get(level) <= 0
	return
}

func (x *Jailhouse[T]) addLevelStep(level TimeRange, current time.Time) time.Time {
	switch level {
	case LAST:
		return current
	case SECOND:
		d, _ := time.ParseDuration("-1s")
		return current.Add(d)
	case MINUTE:
		d, _ := time.ParseDuration("-1m")
		return current.Add(d)
	case HOUR:
		d, _ := time.ParseDuration("-1h")
		return current.Add(d)
	case DAY:
		return current.AddDate(0, 0, -1)
	case WEEK:
		return current.AddDate(0, 0, -7)
	case MONTH:
		return current.AddDate(0, -1, 0)
	case QUARTER:
		return current.AddDate(0, -3, 0)
	case YEAR:
		return current.AddDate(-1, 0, 0)
	case DECADE:
		return current.AddDate(-10, 0, 0)
	case CENTURY:
		return current.AddDate(-100, 0, 0)
	case MILLENIUM:
		return current.AddDate(-1000, 0, 0)
	default:
		err := errors.Errorf("could not find level %s", level)
		panic(err)
	}
}

func (x *Jailhouse[T]) sortResources(input []*JailhouseTimeResource[T]) {
	slices.SortFunc(input, func(a, b *JailhouseTimeResource[T]) int {
		if a.GetTime().Equal(b.GetTime()) {
			return 0
		}
		if a.GetTime().After(b.GetTime()) {
			return -1
		}
		return 1
	})
}

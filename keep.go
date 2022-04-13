package keep

import (
	"fmt"
	"golang.org/x/exp/slices"
	"time"
)

type TimeResource interface {
	GetTime() time.Time
}

type Requirements struct {
	Hours  int
	Days   int
	Weeks  int
	Months int
	Years  int
}

func (x Requirements) Empty() bool {
	return x.Hours == 0 && x.Days == 0 && x.Weeks == 0 && x.Months == 0 && x.Years == 0
}

func (x Requirements) String() string {
	return fmt.Sprintf("%d hours, %d days, %d weeks, %d months, %d years", x.Hours, x.Days, x.Weeks, x.Months, x.Years)
}

func List(input []TimeResource, reqs Requirements) ([]TimeResource, error) {
	return ListForDate(time.Now(), input, reqs)
}

func ListForDate(referenceDate time.Time, input []TimeResource, reqs Requirements) ([]TimeResource, error) {
	result := make([]TimeResource, 0)

	// 1. sort input slice in-place (youngest first)
	sortResources(input)

	// 2. loop and keep or pass
	var (
		currentTime = referenceDate
		nextTime    time.Time
	)
	for _, item := range input {
		if reqs.Empty() {
			break
		}

		// skip due to time?
		if item.GetTime().After(currentTime) {
			// this one is not in the output -> can be dropped
			continue
		}

		// keep this one and move the current time
		result = append(result, item)
		currentTime = item.GetTime()

		// continue?
		nextTime, reqs = nextTick(currentTime, reqs)

		currentTime = nextTime
	}
	return result, nil
}

func sortResources(input []TimeResource) {
	slices.SortFunc(input, func(a, b TimeResource) bool {
		return a.GetTime().After(b.GetTime())
	})
}

func ListRemovableForDate(referenceDate time.Time, input []TimeResource, reqs Requirements) ([]TimeResource, error) {
	keep, err := ListForDate(referenceDate, input, reqs)
	result := make([]TimeResource, 0, len(input)-len(keep))
	if err != nil {
		return result, err
	}

	// 1. sort input slice in-place (youngest first)
	sortResources(input)

out:
	for _, in := range input {
		for _, keepResource := range keep {
			if keepResource == in {
				continue out
			}
		}
		result = append(result, in)
	}

	return result, nil
}

func ListRemovable(input []TimeResource, reqs Requirements) ([]TimeResource, error) {
	return ListRemovableForDate(time.Now(), input, reqs)
}

func nextTick(current time.Time, requirements Requirements) (time.Time, Requirements) {
	newRequirements := requirements
	var newTime time.Time

	if requirements.Hours > 0 {
		d, _ := time.ParseDuration("-1h")
		newTime = current.Add(d)
		newRequirements.Hours--
	} else if requirements.Days > 0 {
		newTime = current.AddDate(0, 0, -1)
		newRequirements.Days--
	} else if requirements.Weeks > 0 {
		newTime = current.AddDate(0, 0, -7)
		newRequirements.Weeks--
	} else if requirements.Months > 0 {
		newTime = current.AddDate(0, -1, 0)
		newRequirements.Months--
	} else if requirements.Years > 0 {
		newTime = current.AddDate(-1, 0, 0)
		newRequirements.Years--
	}

	return newTime, newRequirements
}

# keep

[![Godoc Reference](https://godoc.org/github.com/jojomi/keep?status.svg)](http://godoc.org/github.com/jojomi/keep)
![Go Version](https://img.shields.io/github/go-mod/go-version/jojomi/keep)
![Last Commit](https://img.shields.io/github/last-commit/jojomi/keep)
![Coverage](https://img.shields.io/badge/Coverage-84.3%25-brightgreen)
[![Go Report Card](https://goreportcard.com/badge/jojomi/keep)](https://goreportcard.com/report/jojomi/keep)
[![License](https://img.shields.io/badge/License-MIT-orange.svg)](https://github.com/jojomi/keep/blob/master/LICENSE)

This library helps you decide which elements should be kept for a given number of hours, days, weeks, months, and years.
Most obvious use case are backup files.

## Installation

### CLI tool `keep`

``` go
go install github.com/jojomi/keep/command/keep@latest
```

### Golang library `keep`

``` go
go get github.com/jojomi/keep
```

## Algorithm

Given the configuration below (3 last, 12 hours, 7 days, 12 weeks, 12 months, and 4 years) this library will operate on the list of elements sorted by date, youngest first.
* It will select the first 3 elements from the list (last)
* It then will select 12 elements that are not more than an hour apart (plus one minute extended range), the first of them being the youngest file not yet processed in the set.
* It will then select 7 elements that are not less than a day (plus one hour) apart, the first of them being the youngest file not yet processed in the set.
* and so on as far as levels are defined.
* It is allowed to skip definitions, so you don't have to select daily elements even if you specify hourly and weekly selections.

## Usage

``` go
import "github.com/jojomi/keep"

func main() {
    j := NewDefaultJailhouse()
    j.AddElements(...)

    reqs := NewRequirementsFromMap(map[TimeRange]uint16{
        LAST: 3,
        DAY:  2,
        MONTH: 4,
    })

    j.ApplyRequirements(*reqs)
    kept := j.KeptElements()
    deletable := j.FreeElements()

    // print kept elements
    for _, k := range kept {
        // handle here
    }
	
    // delete now
    for _, k := range deletable {
        // handle here
    }
}
```

Your input data must implement the `TimeResource` interface:

``` go 
type TimeResource interface {
	GetTime() time.Time
}
```

## Development

Add git hooks:

``` shell
git config --local core.hooksPath .githooks/
```
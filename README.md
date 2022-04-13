# keep

[![Godoc Reference](https://godoc.org/github.com/jojomi/keep?status.svg)](http://godoc.org/github.com/jojomi/keep)
![Go Version](https://img.shields.io/github/go-mod/go-version/jojomi/keep)
![Last Commit](https://img.shields.io/github/last-commit/jojomi/keep)
[![Go Report Card](https://goreportcard.com/badge/jojomi/keep)](https://goreportcard.com/report/jojomi/keep)
[![License](https://img.shields.io/badge/License-MIT-orange.svg)](https://github.com/jojomi/keep/blob/master/LICENSE)

This library helps you decide which elements should be kept for a given number of hours, days, weeks, months, and years.
Most obvious use case are backup files.

## Installation

``` go
go get github.com/jojomi/keep
```

## Usage

``` go
import "github.com/jojomi/keep"

func main() {
    reqs := Requirements{
        Hours:  12,
        Days:   7,
        Weeks:  12,
        Months: 12,
        Years:  4,
    },
    k, err := keep.List(input, reqs)
    if err != nil {
        panic(err)
    }
    for _, kk := range k {
        fmt.Println(k.GetTime())
    }
}
```

Your input data must implement the `TimeResource` interface:

``` go 
type TimeResource interface {
	GetTime() time.Time
}
```

If you are interested in the elements that can be _removed_ under the rules of your, swap [`List`](https://pkg.go.dev/github.com/jojomi/keep#List) for [`ListRemovable`](https://pkg.go.dev/github.com/jojomi/keep#ListRemovable).

For tests or special needs, there are also functions named [`ListForDate`](https://pkg.go.dev/github.com/jojomi/keep#ListForDate) and [`ListRemovableForDate`](https://pkg.go.dev/github.com/jojomi/keep#ListRemovableForDate).
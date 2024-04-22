package keep

import (
	"fmt"
	"time"
)

type File struct {
	Filename string
	Time     time.Time
}

func (x File) GetTime() time.Time {
	return x.Time
}

func (x File) String() string {
	return fmt.Sprintf("File %s with date %s", x.Filename, x.GetTime().Format("02.01.2006 15:04:05 Uhr"))
}

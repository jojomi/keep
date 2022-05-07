package keep

import "time"

// TimeResource is an element with an associated time
type TimeResource interface {
	GetTime() time.Time
}

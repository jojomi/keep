package keep

import (
	"fmt"
	"strings"
	"time"
)

// JailhouseTimeResource is an TimeResource with annotated Jailhouse data
type JailhouseTimeResource[T TimeResource] struct {
	Tags         []TimeRangeTag
	TimeResource T
}

func NewJailhouseTimeResource[T TimeResource](resource T) *JailhouseTimeResource[T] {
	x := JailhouseTimeResource[T]{
		TimeResource: resource,
	}
	x.ClearTags()
	return &x
}

func (x JailhouseTimeResource[T]) GetTime() time.Time {
	return x.TimeResource.GetTime()
}

func (x JailhouseTimeResource[T]) GetTags() []TimeRangeTag {
	return x.Tags
}

func (x *JailhouseTimeResource[T]) ClearTags() *JailhouseTimeResource[T] {
	x.Tags = []TimeRangeTag{}
	return x
}

func (x *JailhouseTimeResource[T]) AddTag(tag TimeRangeTag) *JailhouseTimeResource[T] {
	x.Tags = append(x.Tags, tag)
	return x
}

func (x *JailhouseTimeResource[T]) HasLevel(level TimeRange) bool {
	for _, l := range x.Tags {
		if l.TimeRange == level {
			return true
		}
	}
	return false
}

func (x JailhouseTimeResource[T]) IsFree() bool {
	return x.Tags == nil || len(x.Tags) <= 0
}

func (x JailhouseTimeResource[T]) String() string {
	elems := make([]string, len(x.Tags))
	for i, v := range x.Tags {
		elems[i] = v.String()
	}
	return fmt.Sprintf("%s: %s", x.GetTime().Format(time.RFC3339), strings.Join(elems, ", "))
}

package keep

import (
	"fmt"
	"strings"
	"time"
)

// JailhouseTimeResource is an TimeResource with annotated Jailhouse data
type JailhouseTimeResource struct {
	Tags         []TimeRangeTag
	TimeResource TimeResource
}

func NewJailhouseTimeResource(resource TimeResource) *JailhouseTimeResource {
	x := JailhouseTimeResource{
		TimeResource: resource,
	}
	x.ClearTags()
	return &x
}

func (x JailhouseTimeResource) GetTime() time.Time {
	return x.TimeResource.GetTime()
}

func (x JailhouseTimeResource) GetTags() []TimeRangeTag {
	return x.Tags
}

func (x *JailhouseTimeResource) ClearTags() *JailhouseTimeResource {
	x.Tags = []TimeRangeTag{}
	return x
}

func (x *JailhouseTimeResource) AddTag(tag TimeRangeTag) *JailhouseTimeResource {
	x.Tags = append(x.Tags, tag)
	return x
}

func (x *JailhouseTimeResource) HasLevel(level TimeRange) bool {
	for _, l := range x.Tags {
		if l.TimeRange == level {
			return true
		}
	}
	return false
}

func (x JailhouseTimeResource) IsFree() bool {
	return x.Tags == nil || len(x.Tags) <= 0
}

func (x JailhouseTimeResource) String() string {
	elems := make([]string, len(x.Tags))
	for i, v := range x.Tags {
		elems[i] = v.String()
	}
	return fmt.Sprintf("%s: %s", x.GetTime().Format(time.RFC3339), strings.Join(elems, ", "))
}

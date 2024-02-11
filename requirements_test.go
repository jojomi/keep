package keep

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRequirements(t *testing.T) {
	asrt := assert.New(t)

	r := NewRequirements()
	asrt.True(r.IsEmpty())
}

func TestRequirements_String(t *testing.T) {
	type fields struct {
		ranges map[TimeRange]uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "basic",
			fields: fields{
				ranges: map[TimeRange]uint16{
					CENTURY: 2,
					HOUR:    5,
				},
			},
			want: "HOUR=5, CENTURY=2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := Requirements{
				ranges: tt.fields.ranges,
			}
			assert.Equalf(t, tt.want, x.String(), "String()")
		})
	}
}

func TestRequirements_DeepCopy(t *testing.T) {
	r := NewRequirements().Add(WEEK, 3).Add(LAST, 2)

	c := r.DeepCopy()
	assert.Equalf(t, uint16(2), c.Get(LAST), "DeepCopy()")
	assert.Equalf(t, uint16(3), c.Get(WEEK), "DeepCopy()")
}

func TestRequirements_IsEmpty(t *testing.T) {
	type fields struct {
		ranges map[TimeRange]uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "empty",
			fields: fields{
				map[TimeRange]uint16{},
			},
			want: true,
		},
		{
			name: "not empty",
			fields: fields{
				map[TimeRange]uint16{
					YEAR: 3,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := Requirements{
				ranges: tt.fields.ranges,
			}
			assert.Equalf(t, tt.want, x.IsEmpty(), "IsEmpty()")
		})
	}
}

func TestNewRequirementsFromMap(t *testing.T) {
	type args struct {
		data map[TimeRange]uint16
	}
	tests := []struct {
		name string
		args args
		want *Requirements
	}{
		{
			name: "basic",
			args: args{
				map[TimeRange]uint16{
					CENTURY: 4,
					LAST:    1,
				},
			},
			want: NewRequirements().Add(LAST, 1).Add(CENTURY, 4),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewRequirementsFromMap(tt.args.data), "NewRequirementsFromMap(%v)", tt.args.data)
		})
	}
}

func TestRequirements_Get(t *testing.T) {
	type fields struct {
		ranges map[TimeRange]uint16
	}
	type args struct {
		timeRange TimeRange
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint16
	}{
		{
			name: "matching get",
			fields: fields{
				map[TimeRange]uint16{
					HOUR: 10,
				},
			},
			args: args{
				HOUR,
			},
			want: 10,
		},
		{
			name: "get unset element",
			fields: fields{
				map[TimeRange]uint16{
					MILLENIUM: 10,
				},
			},
			args: args{
				MINUTE,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := Requirements{
				ranges: tt.fields.ranges,
			}
			assert.Equalf(t, tt.want, x.Get(tt.args.timeRange), "Get(%v)", tt.args.timeRange)
		})
	}
}

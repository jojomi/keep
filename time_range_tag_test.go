package keep

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTimeRangeTagFrom(t *testing.T) {
	type args struct {
		timeRange TimeRange
		index     uint16
	}
	tests := []struct {
		name string
		args args
		want TimeRangeTag
	}{
		{
			name: "making",
			args: args{
				timeRange: SECOND,
				index:     2,
			},
			want: TimeRangeTag{
				TimeRange: SECOND,
				Index:     2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, TimeRangeTagFrom(tt.args.timeRange, tt.args.index), "TimeRangeTagFrom(%v, %v)", tt.args.timeRange, tt.args.index)
		})
	}
}

func TestTimeRangeTag_String(t *testing.T) {
	type fields struct {
		TimeRange TimeRange
		Index     uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "basic",
			fields: fields{
				TimeRange: QUARTER,
				Index:     3,
			},
			want: "QUARTER-3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := TimeRangeTag{
				TimeRange: tt.fields.TimeRange,
				Index:     tt.fields.Index,
			}
			assert.Equalf(t, tt.want, x.String(), "String()")
		})
	}
}

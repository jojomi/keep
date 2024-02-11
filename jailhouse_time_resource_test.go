package keep

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJailhouseTimeResource_AddLevel(t *testing.T) {
	type fields struct {
		Levels       []TimeRangeTag
		TimeResource TimeResource
	}
	type args struct {
		level TimeRangeTag
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   TimeRange
	}{
		{
			name: "add level",
			fields: fields{
				Levels:       make([]TimeRangeTag, 0),
				TimeResource: nil,
			},
			args: args{
				level: TimeRangeTagFrom(WEEK, 1),
			},
			want: WEEK,
		},
		{
			name: "add multiple levels",
			fields: fields{
				Levels:       make([]TimeRangeTag, 0),
				TimeResource: nil,
			},
			args: args{
				level: TimeRangeTagFrom(DAY, 1),
			},
			want: DAY,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &JailhouseTimeResource{
				Tags:         tt.fields.Levels,
				TimeResource: tt.fields.TimeResource,
			}
			x.AddTag(tt.args.level)
			assert.Truef(t, x.HasLevel(tt.want), "AddTag(%v)", tt.args.level)
		})
	}
}

func TestJailhouseTimeResource_GetTime(t *testing.T) {
	testTime, err := time.Parse("20060102", "20020518")
	assert.NoError(t, err)

	type fields struct {
		Levels       []TimeRangeTag
		TimeResource TimeResource
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		{
			name: "get time",
			fields: fields{
				Levels: nil,
				TimeResource: TestTimeResource{
					t: testTime,
				},
			},
			want: testTime,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := JailhouseTimeResource{
				Tags:         tt.fields.Levels,
				TimeResource: tt.fields.TimeResource,
			}
			assert.Equalf(t, tt.want, x.GetTime(), "GetTime()")
		})
	}
}

func TestJailhouseTimeResource_HasLevel(t *testing.T) {
	type fields struct {
		Levels       []TimeRangeTag
		TimeResource TimeResource
	}
	type args struct {
		level TimeRange
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "has level",
			fields: fields{
				Levels: []TimeRangeTag{
					TimeRangeTagFrom(DAY, 1),
				},
				TimeResource: nil,
			},
			args: args{
				DAY,
			},
			want: true,
		},
		{
			name: "has not level",
			fields: fields{
				Levels: []TimeRangeTag{
					TimeRangeTagFrom(WEEK, 1),
				},
				TimeResource: nil,
			},
			args: args{
				DAY,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &JailhouseTimeResource{
				Tags:         tt.fields.Levels,
				TimeResource: tt.fields.TimeResource,
			}
			assert.Equalf(t, tt.want, x.HasLevel(tt.args.level), "HasLevel(%v)", tt.args.level)
		})
	}
}

func TestJailhouseTimeResource_IsFree(t *testing.T) {
	type fields struct {
		Levels       []TimeRangeTag
		TimeResource TimeResource
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "free",
			fields: fields{
				Levels:       []TimeRangeTag{},
				TimeResource: nil,
			},
			want: true,
		},
		{
			name: "free (nil)",
			fields: fields{
				Levels:       nil,
				TimeResource: nil,
			},
			want: true,
		},
		{
			name: "not free",
			fields: fields{
				Levels:       []TimeRangeTag{TimeRangeTagFrom(DAY, 1)},
				TimeResource: nil,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := JailhouseTimeResource{
				Tags:         tt.fields.Levels,
				TimeResource: tt.fields.TimeResource,
			}
			assert.Equalf(t, tt.want, x.IsFree(), "IsFree()")
		})
	}
}

func TestJailhouseTimeResource_String(t *testing.T) {
	testTime, err := time.Parse("20060102", "20020518")
	assert.NoError(t, err)

	type fields struct {
		Levels       []TimeRangeTag
		TimeResource TimeResource
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "to string",
			fields: fields{
				Levels: []TimeRangeTag{TimeRangeTagFrom(DAY, 3), TimeRangeTagFrom(YEAR, 1)},
				TimeResource: TestTimeResource{
					t: testTime,
				},
			},
			want: "2002-05-18T00:00:00Z: DAY-3, YEAR-1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := JailhouseTimeResource{
				Tags:         tt.fields.Levels,
				TimeResource: tt.fields.TimeResource,
			}
			assert.Equalf(t, tt.want, x.String(), "String()")
		})
	}
}

func TestNewJailhouseTimeResource(t *testing.T) {
	ttr := TestTimeResource{}

	type args struct {
		resource TimeResource
	}
	tests := []struct {
		name string
		args args
		want TimeResource
	}{
		{
			name: "basic",
			args: args{
				resource: ttr,
			},
			want: ttr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewJailhouseTimeResource(tt.args.resource).TimeResource, "NewJailhouseTimeResource(%v)", tt.args.resource)
		})
	}
}

func TestJailhouseTimeResource_GetLevels(t *testing.T) {
	levels := []TimeRangeTag{TimeRangeTagFrom(HOUR, 2)}

	type fields struct {
		Levels       []TimeRangeTag
		TimeResource TimeResource
	}
	tests := []struct {
		name   string
		fields fields
		want   []TimeRangeTag
	}{
		{
			name: "",
			fields: fields{
				Levels:       levels,
				TimeResource: nil,
			},
			want: levels,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := JailhouseTimeResource{
				Tags:         tt.fields.Levels,
				TimeResource: tt.fields.TimeResource,
			}
			assert.Equalf(t, tt.want, x.GetTags(), "GetTags()")
		})
	}
}

type TestTimeResource struct {
	t time.Time
}

func (x TestTimeResource) GetTime() time.Time {
	return x.t
}

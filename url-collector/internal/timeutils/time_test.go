package timeutils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	now := time.Now()
	type args struct {
		from time.Time
		to   time.Time
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "zero times",
			args: args{},
			want: 1,
		},
		{
			name: "from greater than to",
			args: args{
				from: now.AddDate(0, 0, 1),
				to:   now,
			},
			want: 0,
		},
		{
			name: "from equal to to",
			args: args{
				from: now,
				to:   now,
			},
			want: 1,
		},
		{
			name: "to greater than from by 1 day",
			args: args{
				from: now,
				to:   now.AddDate(0, 0, 1),
			},
			want: 2,
		},
		{
			name: "to greater than from by 7 days",
			args: args{
				from: now,
				to:   now.AddDate(0, 0, 7),
			},
			want: 8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, len(List(tt.args.from, tt.args.to)))
		})
	}
}

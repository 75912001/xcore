package time

import (
	"testing"
	"time"
)

// "github.com/stretchr/testify/assert"
func TestAbleUTC(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "normal",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AbleUTC()
			if getInstance().utcAble != tt.want {
				t.Errorf("AbleUTC() = %v, want %v", getInstance().utcAble, tt.want)
			}
		})
	}
}

func TestDayBeginSec(t *testing.T) {
	type args struct {
		timestamp int64
	}
	tests := []struct {
		name    string
		args    args
		preFunc func()
		want    int64
	}{
		{
			name:    "normal-DisableUTC",
			args:    args{timestamp: 1635025000}, //2021-10-24 05:36:40
			preFunc: DisableUTC,
			want:    1635004800, //2021-10-24 00:00:00
		},
		{
			name:    "normal-AbleUTC",
			args:    args{timestamp: 1635025000}, //2021-10-24 05:36:40
			preFunc: AbleUTC,
			want:    1634947200, //2021-10-23 08:00:00
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.preFunc()
			if got := DayBeginSec(tt.args.timestamp); got != tt.want {
				t.Errorf("DayBeginSec() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDayBeginSecByTime(t *testing.T) {
	type args struct {
		t *time.Time
	}
	timeCST := time.Unix(int64(1635025000), 0)       //2021-10-24 05:36:40
	timeUTC := time.Unix(int64(1635025000), 0).UTC() //2021-10-24 05:36:40
	tests := []struct {
		name    string
		args    args
		preFunc func()
		want    int64
	}{
		{
			name:    "normal-DisableUTC",
			args:    args{t: &timeCST},
			preFunc: DisableUTC,
			want:    1635004800, //2021-10-24 00:00:00
		},
		{
			name:    "normal-AbleUTC",
			args:    args{t: &timeUTC},
			preFunc: AbleUTC,
			want:    1634947200, //2021-10-23 08:00:00
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.preFunc()
			if got := DayBeginSecByTime(tt.args.t); got != tt.want {
				t.Errorf("DayBeginSecByTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDisableUTC(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "normal",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DisableUTC()
			if getInstance().utcAble != tt.want {
				t.Errorf("AbleUTC() = %v, want %v", getInstance().utcAble, tt.want)
			}
		})
	}
}

func TestGenYMD(t *testing.T) {
	type args struct {
		timestamp int64
	}
	tests := []struct {
		name    string
		args    args
		preFunc func()
		want    int
	}{
		{
			name:    "normal-DisableUTC",
			args:    args{timestamp: 1635025000}, //2021-10-24 05:36:40
			preFunc: DisableUTC,
			want:    20211024,
		},
		{
			name:    "normal-AbleUTC",
			args:    args{timestamp: 1635025000}, //2021-10-24 05:36:40
			preFunc: AbleUTC,
			want:    20211023,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.preFunc()
			if got := GenYMD(tt.args.timestamp); got != tt.want {
				t.Errorf("GenYMD() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNowTime(t *testing.T) {
	tests := []struct {
		name    string
		preFunc func()
	}{
		{
			name:    "normal-AbleUTC",
			preFunc: AbleUTC,
		},
		{
			name:    "normal-DisableUTC",
			preFunc: DisableUTC,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.preFunc()
			now := NowTime()
			t.Log(now)
		})
	}
}

func TestShadowTimestampSecond(t *testing.T) {
	Update()
	SetTimestampSecondOffset(10)
	tests := []struct {
		name    string
		preFunc func(offset int64)
		want    int64
	}{
		{
			name: "normal",
			want: getInstance().TimestampSecondOffset + getInstance().TimestampSecond,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShadowTimestampSecond(); got != tt.want {
				t.Errorf("ShadowTimestampSecond() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "normal",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Update()
		})
	}
}

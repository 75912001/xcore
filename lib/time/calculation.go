package time

import (
	"strconv"
	"time"
)

// DayBeginSecByTime 当天开始时间戳
func (p *Mgr) DayBeginSecByTime(t *time.Time) int64 {
	if p.utcAble {
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC).Unix()
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
}

// DayBeginSec 返回给定时间戳所在天的开始时间戳
func (p *Mgr) DayBeginSec(timestamp int64) int64 {
	if p.utcAble {
		t := time.Unix(timestamp, 0).UTC()
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC).Unix()
	}
	t := time.Unix(timestamp, 0)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Unix()
}

// GenYMD 获取 e.g.:20210819
//
//	返回YMD
func (p *Mgr) GenYMD(timestamp int64) int {
	var strYMD string
	if p.utcAble {
		strYMD = time.Unix(timestamp, 0).UTC().Format("20060102")
	} else {
		strYMD = time.Unix(timestamp, 0).Format("20060102")
	}
	ymd, _ := strconv.Atoi(strYMD)
	return ymd
}

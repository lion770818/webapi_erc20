package utils

import (
	"time"
)

const (
	TimezoneUTC = "UTC"

	// 改使用官方包的 format
	//TimeLayoutSecond = time.DateTime
)

func TimeNowUTC() time.Time {
	return time.Now().UTC()
}

func TimeNowUnix() int64 {
	return time.Now().Unix()
}

func TimestampConvTime(timestamp int64, timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(timestamp, 0).In(loc), nil
}

func ParseInLocation(timeStr, timezone, layout string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}

	t, err := time.ParseInLocation(layout, timeStr, loc)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

type TimestampConv struct {
	Timestamp int64

	Timezone string
	Layout   string
}

func (tv *TimestampConv) ToTime() (time.Time, error) {
	timeStr, err := tv.ToStr()
	if err != nil {
		return time.Time{}, err
	}

	t, err := ParseInLocation(timeStr, tv.Timezone, tv.Layout)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

func (tv *TimestampConv) ToStr() (string, error) {
	tempTime, err := TimestampConvTime(tv.Timestamp, tv.Timezone)
	if err != nil {
		return "", err
	}

	timeStr := tempTime.Format(tv.Layout)
	return timeStr, nil
}

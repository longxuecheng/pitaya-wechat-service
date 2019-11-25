package utils

import "time"

const (
	TimePrecision_Date TimePrecision = iota
	TimePrecision_Seconds
)

const (
	TimeFormatUnit_Date    string = "2006-1-2"
	TimeFormatUnit_Seconds string = "2006-1-2 15:04:05"
	TimeFormatUnit_Minutes string = "2006-1-2 15:04"
)

type TimePrecision int8

func FormatTime(t time.Time, precision TimePrecision) string {
	formatedTime := ""
	switch precision {
	case TimePrecision_Date:
		formatedTime = t.Format(TimeFormatUnit_Date)
	case TimePrecision_Seconds:
		formatedTime = t.Format(TimeFormatUnit_Seconds)
	}
	return formatedTime
}

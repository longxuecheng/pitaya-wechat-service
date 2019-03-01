package utils

import "time"

const (
	TimePrecision_Year TimePrecision = iota
	TimePrecision_Date
	TimePrecision_Seconds
)

type TimeFormatUnit string

const (
	TimeFormatUnit_Year TimeFormatUnit = ""
)

type TimePrecision int8

func FormatTime(t time.Time, layout string, precision TimePrecision) string {

	switch precision {
	case TimePrecision_Year:
	case TimePrecision_Date:
	case TimePrecision_Seconds:

	}
	return ""
}

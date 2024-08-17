package helper

import (
	"fmt"
	"strconv"
	"time"
)

func ParseWindowsTime(windowsTime string) (time.Time, error) {
	layout := "20060102150405.000000-070"

	t, err := time.Parse(layout, windowsTime)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse Windows time: %v", err)
	}

	return t, nil
}

func FormatWindowsTimeCustom(t time.Time) string {
	return t.Format("2006-01-02-15-04-05")
}
func parseUint32(s string) uint32 {
	v, _ := strconv.ParseUint(s, 10, 32)
	return uint32(v)
}
func parseUint64(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

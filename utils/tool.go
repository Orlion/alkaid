package utils

import (
	"bytes"
	"strconv"
	"time"
)

func FriendlyTime(t time.Time) string {
	var (
		now time.Time
		nowUnix int64
		nowYear int
		diffUnix int64
		tYear int
		formatTime string
		buffer bytes.Buffer
	)

	now = time.Now()
	nowUnix = now.Unix()
	diffUnix = nowUnix - t.Unix()
	nowYear = now.Year()
	tYear = t.Year()
	if tYear == nowYear {
		if diffUnix < 60 {
			formatTime = "刚刚"
		} else if diffUnix >= 60 && diffUnix < 3600 {
			buffer.WriteString(strconv.Itoa(int(diffUnix / 60)))
			buffer.WriteString("分钟前")
			formatTime = buffer.String()
		} else if diffUnix >= 3600 && diffUnix < 86400 {
			buffer.WriteString(strconv.Itoa(int(diffUnix / 3600)))
			buffer.WriteString("小时前")
			formatTime = buffer.String()
		} else {
			formatTime = t.Format("01-02 15:04")
		}
	} else {
		formatTime = t.Format("2006-01-02")
	}

	return formatTime
}
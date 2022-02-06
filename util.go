package utils

import (
	"fmt"
	"strconv"
	"time"
)

func FormatUnixTimestampToPgComparisonValue(timestamp int) string {
	df := time.Unix(int64(timestamp), 0)
	year, month, date := df.Date()
	hour, minute, second := df.Clock()
	var months, days, hours, minutes, seconds string

	if Months[month] < 9 {
		months = "0" + strconv.Itoa(Months[month])
	} else {
		months = strconv.Itoa(Months[month])
	}
	if date < 9 {
		days = "0" + strconv.Itoa(date)
	} else {
		days = strconv.Itoa(date)
	}

	if hour <= 9 {
		hours = "0" + strconv.Itoa(hour)
	} else {
		hours = strconv.Itoa(hour)
	}
	if minute <= 9 {
		minutes = "0" + strconv.Itoa(minute)
	} else {
		minutes = strconv.Itoa(minute)
	}
	if second <= 9 {
		seconds = "0" + strconv.Itoa(second)
	} else {
		seconds = strconv.Itoa(second)
	}

	return fmt.Sprintf("'%d-%s-%s %s:%s:%s.000000'::timestamp",
		year, months, days, hours, minutes, seconds)
}

var Months = map[time.Month]int{
	time.January:   1,
	time.February:  2,
	time.March:     3,
	time.April:     4,
	time.May:       5,
	time.June:      6,
	time.July:      7,
	time.August:    8,
	time.September: 9,
	time.October:   10,
	time.November:  11,
	time.December:  12,
}

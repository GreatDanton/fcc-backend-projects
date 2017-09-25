package utilities

import (
	"fmt"
	"time"
)

// helper constants for calculating time difference
const diffMonths = 30.0 * 24      // how many hours there are in an average month
const diffYears = diffMonths * 12 // how many hours there are in an average year

// TimeDiff returns time difference between inserted time and time now in suitable
//format according to number of:
// seconds if time < 1min
// minutes if time < 1hour
// hours if time < 1 day
// days if time < 1 month
// months if time < 1 year
// years if time > 1 year
func TimeDiff(t time.Time) string {
	diff := time.Since(t)
	h := diff.Hours()
	m := diff.Minutes()
	s := diff.Seconds()

	// if number of years > 1 display number of years since the inserted time t
	numOfYears := h / diffYears
	str, ok := timeDiffFormat(numOfYears, "year", "years")
	// if everything is okay, display result in number of years
	if ok {
		return str
	}

	// display number of months ago
	numOfMonths := h / diffMonths
	str, ok = timeDiffFormat(numOfMonths, "month", "months")
	if ok {
		return str
	}

	// check number of days ago
	numOfDays := h / 24
	str, ok = timeDiffFormat(numOfDays, "day", "days")
	if ok {
		return str
	}

	// check number of hours ago
	str, ok = timeDiffFormat(h, "hour", "hours")
	if ok {
		return str
	}

	// check number of minutes ago
	str, ok = timeDiffFormat(m, "minute", "minutes")
	if ok {
		return str
	}

	// if nothing of above applies, return number of seconds ago
	str, _ = timeDiffFormat(s, "second", "seconds")
	return str
}

// timeDiffFormat checks if timeInput is suitable for displaying timeInput in correct
// format ex: [1 second, 2 seconds]. It's returning false if timeInput is 0 and therefore
// not suitable for displaying
func timeDiffFormat(timeInput float64, singleUnitName string, multiUnitName string) (string, bool) {
	time := Round(timeInput)
	if time < 1 {
		return "", false
	}
	ending := multiUnitName
	if time == 1 {
		ending = singleUnitName
	}
	msg := fmt.Sprintf("%v %v", time, ending)
	ok := true
	return msg, ok
}

// Round is utility function for rounding floats to nearest integer
func Round(val float64) int {
	if val < 0 {
		return int(val - 0.5)
	}
	return int(val + 0.5)
}

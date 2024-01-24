package util

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var now = func() time.Time { return time.Now() }
var loc, _ = time.LoadLocation("Europe/Amsterdam")

type TimeParseResult struct {
	Time        time.Time
	isDateEmpty bool
	isTimeEmpty bool
	isSecsEmpty bool
}

func ParseTimeFrom(input string) (time.Time, error) {
	res, err := ParseTime(input)

	if res.isDateEmpty {
		timeString := res.Time.Format("15:04:05")
		today := now().Format("2006-01-02")

		newRes, _ := time.ParseInLocation("2006-01-02 15:04:05", today+" "+timeString, loc)

		if newRes.After(now()) {
			yesterday := now().Add(-24 * time.Hour).Format("2006-01-02")
			newRes, _ = time.ParseInLocation("2006-01-02 15:04:05", yesterday+" "+timeString, loc)
		}

		res.Time = newRes
	}
	if res.isTimeEmpty {
		// No need to do anything, Time defaults to 00:00:00
	}

	return res.Time, err
}

func ParseTimeUntil(input string) (time.Time, error) {
	res, err := ParseTime(input)

	if res.isDateEmpty {
		timeString := res.Time.Format("15:04:05")
		today := now().Format("2006-01-02")

		newRes, _ := time.ParseInLocation("2006-01-02 15:04:05", today+" "+timeString, loc)

		if newRes.Before(now()) {
			tomorrow := now().Add(24 * time.Hour).Format("2006-01-02")
			newRes, _ = time.ParseInLocation("2006-01-02 15:04:05", tomorrow+" "+timeString, loc)
		}

		res.Time = newRes
	}
	if res.isTimeEmpty {
		res.Time = res.Time.Add(24 * time.Hour).Add(-1 * time.Second)
	}
	if res.isSecsEmpty {
		res.Time = res.Time.Add(59 * time.Second)
	}

	return res.Time, err
}

func ParseTime(input string) (TimeParseResult, error) {
	if strings.TrimSpace(input) == "" {
		return TimeParseResult{
			Time:        time.Now(),
			isDateEmpty: false,
			isTimeEmpty: false,
			isSecsEmpty: false,
		}, nil
	}
	dateFormats := []string{
		"",
		"2006-01-02",
		"2006-1-02",
		"02-01-2006",
		"2-1-2006",
		"02/01/2006",
		"2/1/2006",
	}
	timeFormats := []string{
		"",
		"15:04",
		"15:04:05",
		"3:04",
		"3:04:05",
	}
	for _, dateFormat := range dateFormats {
		for _, timeFormat := range timeFormats {
			separator := func() string {
				if dateFormat == "" || timeFormat == "" {
					return ""
				}
				return " "
			}()
			res, err := time.ParseInLocation(dateFormat+separator+timeFormat, input, loc)
			if err == nil {
				return TimeParseResult{
					Time:        res,
					isDateEmpty: dateFormat == "",
					isTimeEmpty: timeFormat == "",
					isSecsEmpty: strings.Count(timeFormat, ":") == 1,
				}, nil
			}
		}
	}

	return TimeParseResult{}, errors.New(fmt.Sprintf("failed to parse Time '%v'", input))
}

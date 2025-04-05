package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	ErrEmptyRepeat      = errors.New("empty repeat rule")
	ErrInvalidFormat    = errors.New("invalid repeat format")
	ErrInvalidDate      = errors.New("invalid date format")
	ErrInvalidDay       = errors.New("invalid day value")
	ErrInvalidMonth     = errors.New("invalid month value")
	ErrInvalidWeekday   = errors.New("invalid weekday value")
	ErrIntervalTooLarge = errors.New("interval too large")
)

func NextDate(now time.Time, dateStr string, repeat string) (string, error) {
	if repeat == "" {
		return "", ErrEmptyRepeat
	}

	date, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		return "", ErrInvalidDate
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", ErrInvalidFormat
	}

	switch parts[0] {
	case "d":
		return dailyRepeat(now, date, parts)
	case "y":
		return yearlyRepeat(now, date)
	case "w":
		return weeklyRepeat(now, date, parts)
	case "m":
		return monthlyRepeat(now, date, parts)
	default:
		return "", ErrInvalidFormat
	}
}

func dailyRepeat(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", ErrInvalidFormat
	}

	interval, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", ErrInvalidFormat
	}

	if interval <= 0 || interval > 400 {
		return "", ErrIntervalTooLarge
	}

	for {
		date = date.AddDate(0, 0, interval)
		if afterNow(date, now) {
			break
		}
	}

	return date.Format(DateFormat), nil
}

func yearlyRepeat(now, date time.Time) (string, error) {
	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			break
		}
	}
	return date.Format(DateFormat), nil
}

func weeklyRepeat(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", ErrInvalidFormat
	}

	daysStr := strings.Split(parts[1], ",")
	days := make([]int, 0, len(daysStr))

	for _, dayStr := range daysStr {
		day, err := strconv.Atoi(dayStr)
		if err != nil || day < 1 || day > 7 {
			return "", ErrInvalidWeekday
		}
		days = append(days, day)
	}

	for {
		date = date.AddDate(0, 0, 1)
		if afterNow(date, now) {
			weekday := int(date.Weekday())
			if weekday == 0 {
				weekday = 7 // Sunday is 7
			}

			for _, d := range days {
				if weekday == d {
					return date.Format(DateFormat), nil
				}
			}
		}
	}
}

func monthlyRepeat(now, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", ErrInvalidFormat
	}

	daysStr := strings.Split(parts[1], ",")
	days := make([]int, 0, len(daysStr))

	for _, dayStr := range daysStr {
		day, err := strconv.Atoi(dayStr)
		if err != nil || day < -2 || day == 0 || day > 31 {
			return "", ErrInvalidDay
		}
		days = append(days, day)
	}

	var months []int
	if len(parts) > 2 {
		monthsStr := strings.Split(parts[2], ",")
		months = make([]int, 0, len(monthsStr))
		for _, monthStr := range monthsStr {
			month, err := strconv.Atoi(monthStr)
			if err != nil || month < 1 || month > 12 {
				return "", ErrInvalidMonth
			}
			months = append(months, month)
		}
	}

	for {
		date = date.AddDate(0, 0, 1)
		if afterNow(date, now) {
			_, month, day := date.Date()
			currentMonth := int(month)

			if len(months) > 0 {
				monthMatch := false
				for _, m := range months {
					if m == currentMonth {
						monthMatch = true
						break
					}
				}
				if !monthMatch {
					continue
				}
			}

			for _, d := range days {
				switch {
				case d > 0 && day == d:
					return date.Format(DateFormat), nil
				case d == -1 && isLastDayOfMonth(date):
					return date.Format(DateFormat), nil
				case d == -2 && isPenultimateDayOfMonth(date):
					return date.Format(DateFormat), nil
				}
			}
		}
	}
}

func afterNow(date, now time.Time) bool {
	return date.After(now)
}

func isLastDayOfMonth(date time.Time) bool {
	return date.AddDate(0, 0, 1).Month() != date.Month()
}

func isPenultimateDayOfMonth(date time.Time) bool {
	return date.AddDate(0, 0, 2).Month() != date.Month()
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, "Invalid now parameter", http.StatusBadRequest)
			return
		}
	}

	result, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, result)
}

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
	ErrIntervalTooLarge = errors.New("interval must be between 1 and 400")
)

// NextDate вычисляет следующую дату выполнения задачи
func NextDate(now time.Time, dateStr string, repeat string) (string, error) {
	if repeat == "" {
		return "", ErrEmptyRepeat
	}

	date, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		return "", ErrInvalidDate
	}

	parts := strings.Fields(repeat)
	if len(parts) < 1 {
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

// dailyRepeat обрабатывает ежедневное повторение
func dailyRepeat(now, date time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", ErrInvalidFormat
	}

	interval, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", ErrInvalidFormat
	}

	if interval < 1 || interval > 400 {
		return "", ErrIntervalTooLarge
	}

	nextDate := date
	for {
		nextDate = nextDate.AddDate(0, 0, interval)
		if isAfter(trimToDay(nextDate), trimToDay(now)) {
			break
		}
	}

	return nextDate.Format(DateFormat), nil
}

// yearlyRepeat обрабатывает ежегодное повторение
func yearlyRepeat(now, date time.Time) (string, error) {
	nextDate := date
	for {
		nextDate = nextDate.AddDate(1, 0, 0)
		if isAfter(trimToDay(nextDate), trimToDay(now)) {
			break
		}
	}
	return nextDate.Format(DateFormat), nil
}

// weeklyRepeat обрабатывает еженедельное повторение
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

	nextDate := date
	for {
		nextDate = nextDate.AddDate(0, 0, 1)
		if isAfter(trimToDay(nextDate), trimToDay(now)) {
			weekday := int(nextDate.Weekday())
			if weekday == 0 {
				weekday = 7 // Воскресенье = 7
			}
			for _, d := range days {
				if weekday == d {
					return nextDate.Format(DateFormat), nil
				}
			}
		}
	}
}

// monthlyRepeat обрабатывает ежемесячное повторение
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

	nextDate := date
	for {
		nextDate = nextDate.AddDate(0, 0, 1)
		if isAfter(trimToDay(nextDate), trimToDay(now)) {
			currentMonth := int(nextDate.Month())
			currentDay := nextDate.Day()

			if len(months) > 0 {
				validMonth := false
				for _, m := range months {
					if m == currentMonth {
						validMonth = true
						break
					}
				}
				if !validMonth {
					continue
				}
			}

			for _, d := range days {
				switch {
				case d > 0 && currentDay == d:
					return nextDate.Format(DateFormat), nil
				case d == -1 && isLastDayOfMonth(nextDate):
					return nextDate.Format(DateFormat), nil
				case d == -2 && isPenultimateDayOfMonth(nextDate):
					return nextDate.Format(DateFormat), nil
				}
			}
		}
	}
}

// Вспомогательные функции
func trimToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func isAfter(a, b time.Time) bool {
	return a.After(b)
}

func isLastDayOfMonth(date time.Time) bool {
	return date.AddDate(0, 0, 1).Month() != date.Month()
}

func isPenultimateDayOfMonth(date time.Time) bool {
	return date.AddDate(0, 0, 2).Month() != date.Month()
}

// NextDateHandler обрабатывает HTTP-запросы для вычисления следующей даты
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

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

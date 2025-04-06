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

// NextDate - функция для вычисления следующей даты на основе текущей даты, даты задачи и правила повторения.
// Она выполняет следующие шаги:
// 1. Проверка наличия правила повторения. Если правило не указано, то возвращается ошибка.
// 2. Парсинг даты задачи. Если дата имеет неверный формат, то возвращается ошибка.
// 3. Разделение правила повторения на части.
// 4. В зависимости от первой части правила повторения вызывается соответствующая функция для вычисления следующей даты:
//   - "d" - dailyRepeat
//   - "y" - yearlyRepeat
//   - "w" - weeklyRepeat
//   - "m" - monthlyRepeat
//
// 5. Если первая часть правила повторения не поддерживается, то возвращается ошибка.
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

// dailyRepeat - функция для вычисления следующей даты на основе текущей даты, даты задачи и правила повторения "d".
// Она выполняет следующие шаги:
// 1. Проверка формата правила повторения. Если правило имеет неверный формат, то возвращается ошибка.
// 2. Преобразование второй части правила повторения в целое число. Если преобразование не удалось, то возвращается ошибка.
// 3. Проверка интервала повторения. Если интервал меньше или равен 0 или больше 400, то возвращается ошибка.
// 4. Вычисление следующей даты с учетом интервала повторения. Если следующая дата уже наступила, то вычисление продолжается.
// 5. Возвращение следующей даты в формате YYYYMMDD.
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

// yearlyRepeat - функция для вычисления следующей даты на основе текущей даты и даты задачи с ежегодным повторением.
// Она выполняет следующие шаги:
// 1. Вычисление следующей даты с учетом ежегодного повторения. Если следующая дата уже наступила, то вычисление продолжается.
// 2. Возвращение следующей даты в формате YYYYMMDD.
func yearlyRepeat(now, date time.Time) (string, error) {
	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			break
		}
	}
	return date.Format(DateFormat), nil
}

// weeklyRepeat - функция для вычисления следующей даты на основе текущей даты, даты задачи и правила повторения "w".
// Она выполняет следующие шаги:
// 1. Проверка формата правила повторения. Если правило имеет неверный формат, то возвращается ошибка.
// 2. Преобразование второй части правила повторения в список дней недели. Если преобразование не удалось или день недели имеет неверное значение, то возвращается ошибка.
// 3. Вычисление следующей даты с учетом еженедельного повторения. Если следующая дата уже наступила, то вычисление продолжается.
// 4. Проверка, является ли день недели следующей даты одним из дней, указанных в правиле повторения. Если да, то возвращается следующая дата в формате YYYYMMDD.
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

// monthlyRepeat - функция для вычисления следующей даты на основе текущей даты, даты задачи и правила повторения "m".
// Она выполняет следующие шаги:
// 1. Проверка формата правила повторения. Если правило имеет неверный формат, то возвращается ошибка.
// 2. Преобразование второй части правила повторения в список дней месяца. Если преобразование не удалось или день месяца имеет неверное значение, то возвращается ошибка.
// 3. Если в правиле повторения указаны месяцы, то преобразование третьей части правила повторения в список месяцев. Если преобразование не удалось или месяц имеет неверное значение, то возвращается ошибка.
// 4. Вычисление следующей даты с учетом ежемесячного повторения. Если следующая дата уже наступила, то вычисление продолжается.
// 5. Проверка, является ли день месяца следующей даты одним из дней, указанных в правиле повторения. Если да, то возвращается следующая дата в формате YYYYMMDD.
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

// afterNow - функция для проверки, является ли дата после текущей даты.
// Она выполняет следующие шаги:
// 1. Сравнение даты с текущей датой.
// 2. Возвращение true, если дата после текущей даты, иначе false.
func afterNow(date, now time.Time) bool {
	return date.After(now)
}

// isLastDayOfMonth - функция для проверки, является ли дата последним днем месяца.
// Она выполняет следующие шаги:
// 1. Добавление одного дня к дате.
// 2. Сравнение месяца следующей даты с месяцем исходной даты.
// 3. Возвращение true, если месяцы различаются, иначе false.
func isLastDayOfMonth(date time.Time) bool {
	return date.AddDate(0, 0, 1).Month() != date.Month()
}

// isPenultimateDayOfMonth - функция для проверки, является ли дата предпоследним днем месяца.
// Она выполняет следующие шаги:
// 1. Добавление двух дней к дате.
// 2. Сравнение месяца следующей даты с месяцем исходной даты.
// 3. Возвращение true, если месяцы различаются, иначе false.
func isPenultimateDayOfMonth(date time.Time) bool {
	return date.AddDate(0, 0, 2).Month() != date.Month()
}

// NextDateHandler - обработчик HTTP-запросов для вычисления следующей даты на основе текущей даты, даты задачи и правила повторения.
// Он выполняет следующие шаги:
// 1. Проверка метода запроса. Если метод не GET, то возвращается ошибка.
// 2. Получение параметров запроса: now, date и repeat.
// 3. Если параметр now не указан, то используется текущая дата. В противном случае, параметр now парсится в формате YYYYMMDD. Если парсинг не удался, то возвращается ошибка.
// 4. Вычисление следующей даты с использованием функции NextDate. Если возникает ошибка, то возвращается ошибка.
// 5. Установка заголовка Content-Type в "text/plain".
// 6. Возвращение следующей даты в формате YYYYMMDD.
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

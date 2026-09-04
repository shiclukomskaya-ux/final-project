package planner

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("Пустая строка")
	}

	t, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", err
	}

	parts := strings.Split(repeat, " ")
	switch parts[0] {
	case "y":
		for {
			t = t.AddDate(1, 0, 0)
			if t.After(now) {
				break
			}
		}
		return t.Format(dateFormat), nil

	case "d":
		if len(parts) < 2 {
			return "", fmt.Errorf("Не указан интервал в днях")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("Некорректный интервал")
		}
		if days <= 0 || days > 400 {
			return "", fmt.Errorf("Превышен допустимы интервал")
		}
		for {
			t = t.AddDate(0, 0, days)
			if t.After(now) {
				break
			}
		}
		return t.Format(dateFormat), nil
	default:
		return "", fmt.Errorf("Неподдерживаемый формат")
	}
}

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
		return "", fmt.Errorf("Правило повторения не указано")
	}

	t, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("Некорректная начальная дата: %w", err)
	}

	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return "", fmt.Errorf("Неверный формат правила")
	}
	switch parts[0] {
	case "y":
		for {
			t = t.AddDate(1, 0, 0)
			if afterNow(t, now) {
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
			return "", fmt.Errorf("Некорректный интервал: %w", err)
		}
		if days <= 0 || days > 400 {
			return "", fmt.Errorf("Интервал должен быть от 1 до 400 дней")
		}
		for {
			t = t.AddDate(0, 0, days)
			if afterNow(t, now) {
				break
			}
		}
		return t.Format(dateFormat), nil
	}
	return "", fmt.Errorf("Неподдерживаемый формат")
}
func afterNow(date, now time.Time) bool {
	return date.After(now)
}

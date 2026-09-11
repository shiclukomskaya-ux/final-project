package api

import (
	"encoding/json"
	"final-project/internal/database"
	"final-project/internal/planner"
	"net/http"
	"strconv"
	"time"
)

func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task database.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "Ошибка десериализации JSON"})
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "Некорректный формат даты или правила повторения"})
		return
	}

	id, err := database.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "Не удалось добавить задачу в базу данных"})
		return
	}
	writeJson(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}
func checkDate(task *database.Task) error {
	now := time.Now()
	todayStr := now.Format("20060102")

	now, _ = time.Parse("20060102", todayStr)

	if task.Date == "" {
		task.Date = todayStr
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		return err
	}

	if planner.AfterNow(now, t) {
		if task.Repeat == "" {
			task.Date = todayStr
		} else {
			next, err := planner.NextDate(now, todayStr, task.Repeat)
			if err != nil {
				return err
			}
			task.Date = next
		}
		return nil
	}

	if task.Repeat != "" {
		_, err := planner.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	return nil
}

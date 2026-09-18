package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"final-project/internal/database"
	"final-project/internal/planner"
)

func writeJson(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task database.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Ошибка десериализации JSON"})
		return
	}

	if task.Title == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Некорректный формат даты или правила повторения"})
		return
	}

	id, err := database.AddTask(&task)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Не удалось добавить задачу в базу данных"})
		return
	}
	writeJson(w, http.StatusOK, map[string]string{"id": strconv.FormatInt(id, 10)})
}
func checkDate(task *database.Task) error {
	now := time.Now()
	todayStr := now.Format(apiDateFormat)

	now, _ = time.Parse(apiDateFormat, todayStr)

	if task.Date == "" {
		task.Date = todayStr
	}

	t, err := time.Parse(apiDateFormat, task.Date)
	if err != nil {
		return err
	}

	if planner.AfterNow(now, t) {
		if task.Repeat == "" {
			task.Date = todayStr
		} else {
			next, err := planner.NextDate(now, task.Date, task.Repeat)
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

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

	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		writeJson(w, map[string]string{"error": "Некорректный формат даты"})
		return
	}

	var next string
	if task.Repeat != "" {
		next, err = planner.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJson(w, map[string]string{"error": "Некорректный формат повторения"})
			return
		}
	}

	if planner.AfterNow(now, t) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format("20060102")
		} else {
			task.Date = next
		}
	}

	id, err := database.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "Не удалось добавить задачу в базу данных"})
		return
	}
	writeJson(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}

package api

import (
	"encoding/json"
	"net/http"
	"time"

	"final-project/internal/database"
	"final-project/internal/planner"
)

const apiDateFormat = "20060102"

func Init() {
	http.HandleFunc("GET /api/nextdate", NextDayHandler)
	http.HandleFunc("GET /api/tasks", TasksHandler)
	http.HandleFunc("POST /api/task/done", taskDoneHandler)

	http.HandleFunc("GET /api/task", getTaskHandler)
	http.HandleFunc("POST /api/task", addTaskHandler)
	http.HandleFunc("PUT /api/task", updateTaskHandler)
	http.HandleFunc("DELETE /api/task", deleteTaskHandler)
}
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}
	task, err := database.GetTask(id)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Задача не найдена"})
		return
	}
	writeJson(w, http.StatusOK, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task database.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Ошибка десериализации JSON"})
		return
	}
	if task.ID == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
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

	err = database.UpdateTask(&task)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Не удалось обновить задачу в базе данных"})
		return
	}
	writeJson(w, http.StatusOK, map[string]string{})
}

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := database.GetTask(id)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Задача не найдена"})
		return
	}

	if task.Repeat == "" {
		err = database.DeleteTask(id)
		if err != nil {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "Не удалось удалить задачу"})
			return
		}
		writeJson(w, http.StatusOK, map[string]string{})
		return
	}
	now := time.Now()
	nextDate, err := planner.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Некорректные данные"})
		return
	}
	err = database.UpdateDate(nextDate, id)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Не удалось обновить дату задачи"})
		return
	}

	writeJson(w, http.StatusOK, map[string]string{})
}
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}
	err := database.DeleteTask(id)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Не удалось удалить задачу"})
		return
	}
	writeJson(w, http.StatusOK, map[string]string{})
}

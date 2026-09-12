package api

import (
	"encoding/json"
	"final-project/internal/database"
	"final-project/internal/planner"
	"net/http"
	"time"
)

func Init() {
	http.HandleFunc("/api/nextdate", NextDayHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", TasksHandler)
	http.HandleFunc("/api/task/done", taskDoneHandler)
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		getTaskHandler(w, r)

	case http.MethodPut:
		updateTaskHandler(w, r)

	case http.MethodDelete:
		deleteTaskHandler(w, r)

	case http.MethodPost:
		addTaskHandler(w, r)
	}
}
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}
	task, err := database.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "Задача не найдена"})
		return
	}
	writeJson(w, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task database.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "Ошибка десериализации JSON"})
		return
	}
	if task.ID == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор"})
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

	err = database.UpdateTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": "Не удалось обновить задачу в базе данных"})
		return
	}
	writeJson(w, map[string]string{})
}

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := database.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "Задача не найдена"})
		return
	}

	if task.Repeat == "" {
		err = database.DeleteTask(id)
		if err != nil {
			writeJson(w, map[string]string{"error": "Не удалось удалить задачу"})
			return
		}
		writeJson(w, map[string]string{})
		return
	}
	now := time.Now()
	nextDate, err := planner.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeJson(w, map[string]string{"error": "Некорректные данные"})
		return
	}
	err = database.UpdateDate(nextDate, id)
	if err != nil {
		writeJson(w, map[string]string{"error": "Не удалось обновить дату задачи"})
		return
	}

	writeJson(w, map[string]string{})
}
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}
	err := database.DeleteTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": "Не удалось удалить задачу"})
		return
	}
	writeJson(w, map[string]string{})
}

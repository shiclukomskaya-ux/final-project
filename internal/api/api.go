package api

import "net/http"

func Init() {
	http.HandleFunc("/api/nextdate", NextDayHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", TasksHandler)
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		getTaskHandler(w, r)

	// case http.MethodDelete:

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

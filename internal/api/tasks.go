package api

import (
	"final-project/internal/database"
	"net/http"
)

type TasksResp struct {
	Tasks []*database.Task `json:"tasks"`
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := database.Tasks(50)
	if err != nil {
		writeJson(w, map[string]string{"error": "Не удалось получить список задач"})
		return
	}
	if tasks == nil {
		tasks = make([]*database.Task, 0)
	}

	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}

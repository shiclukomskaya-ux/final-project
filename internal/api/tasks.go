package api

import (
	"net/http"

	"final-project/internal/database"
)

const taskLimit = 50

type TasksResp struct {
	Tasks []*database.Task `json:"tasks"`
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := database.Tasks(taskLimit)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Не удалось получить список задач"})
		return
	}
	if tasks == nil {
		tasks = make([]*database.Task, 0)
	}

	writeJson(w, http.StatusOK, TasksResp{
		Tasks: tasks,
	})
}

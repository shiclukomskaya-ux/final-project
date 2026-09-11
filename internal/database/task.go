package database

type Task struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler(date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}
func Tasks(limit int) ([]*Task, error) {
	var tasks []*Task

	rows, err := DB.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?", limit)

	if err != nil {
		return tasks, err
	}
	defer rows.Close()
	for rows.Next() {
		var task Task
		rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		tasks = append(tasks, &task)
	}
	return tasks, nil
}
func UpdateTask(task *Task) error {
	query := 
}

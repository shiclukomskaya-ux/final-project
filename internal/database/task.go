package database

import "fmt"

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
	tasks := make([]*Task, 0)

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

func GetTask(id string) (*Task, error) {

	var task Task
	err := DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}
	return &task, nil

}

func UpdateTask(task *Task) error {
	res, err := DB.Exec("UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?", task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("Задача с таким id не найдена")
	}
	return nil
}
func DeleteTask(id string) error {
	res, err := DB.Exec("DELETE FROM scheduler WHERE id=?", id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("Задача не найдена")
	}
	return nil
}
func UpdateDate(next string, id string) error {
	res, err := DB.Exec("UPDATE scheduler SET date=? WHERE id=?", next, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("Задача не найдена")
	}
	return nil
}

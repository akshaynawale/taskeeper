package taskdb

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

type Task struct {
	TaskID int
	UserID int // task corrosponding to the UserID
	Title  string
	Body   string
}

// CreateTask creates a task entry into the database
func CreateTask(db *TaskeeperDB, userID int, title string, body string) error {
	query := sq.Insert("tasks").Columns("userID", "title", "body").Values(userID, title, body)
	_, err := query.RunWith(db.db).Query()
	return err
}

// GetTasks gets all tasks for a user from db
func GetTasks(db *TaskeeperDB, userID int) ([]Task, error) {
	query := sq.Select("taskID", "userID", "title", "body").From("tasks").Where(sq.Eq{"userID": userID})
	rows, err := query.RunWith(db.db).Query()
	if err != nil {
		return []Task{}, fmt.Errorf("failed to run select query on db: %v", err)
	}
	// scan each row
	tasks := []Task{}
	for rows.Next() {
		t := Task{}
		if err := rows.Scan(&t.TaskID, &t.UserID, &t.Title, &t.Body); err != nil {
			return tasks, fmt.Errorf("failed to read task rows: %v", err)
		}
		tasks = append(tasks, t)
	}
	return tasks, err
}

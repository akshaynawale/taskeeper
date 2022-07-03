package taskdb

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

type User struct {
	UserID   int64
	Username string // username
	Password string // password encoded string
}

// CreateUser creates new user in the db
func CreateUser(db *TaskeeperDB, uname, pass string) error {
	query := sq.Insert("users").Columns("username", "password").
		Values(uname, pass)
	_, err := query.RunWith(db.db).Query()
	return err
}

//GetUser gets the User struct from db with userID
func GetUser(db *TaskeeperDB, userID int) (User, error) {
	query := sq.Select("userID, username, password").From("users").Where(sq.Eq{"userID": userID})
	rows, err := query.RunWith(db.db).Query()
	if err != nil {
		return User{}, fmt.Errorf("failed to run query select with error: %v", err)
	}
	// if we  dont see any error  return the first result
	usrs := []User{}
	for rows.Next() {
		u := User{}
		if err := rows.Scan(&u.UserID, &u.Username, &u.Password); err != nil {
			return User{}, fmt.Errorf("failed to read row from db: %v", err)
		}
		usrs = append(usrs, u)
	}
	if len(usrs) != 1 {
		return User{}, fmt.Errorf("something went wrong when trying to find user in db with userid: %d", userID)
	}
	return usrs[0], err
}

//GetUserByName gets the User struct from db with userID
func GetUserByName(db *TaskeeperDB, username string) (*User, error) {
	query := sq.Select("userID, username, password").From("users").Where(sq.Eq{"username": username})
	rows, err := query.RunWith(db.db).Query()
	if err != nil {
		return nil, fmt.Errorf("failed to run query select with error: %v", err)
	}
	// if we  dont see any error  return the first result
	usrs := []User{}
	for rows.Next() {
		u := User{}
		if err := rows.Scan(&u.UserID, &u.Username, &u.Password); err != nil {
			return nil, fmt.Errorf("failed to read row from db: %v", err)
		}
		usrs = append(usrs, u)
	}
	if len(usrs) > 1 {
		return nil, fmt.Errorf("something went wrong when trying to find user in db with username: %s", username)
	}
	if len(usrs) == 1 {
		return &usrs[0], err
	}
	return nil, nil
}

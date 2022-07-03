package taskdb

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
)

type TaskeeperDB struct {
	db *sql.DB // mysql db connection
}

// ConnectDb function creates database connection and returns CaterDB struct
func ConnectDb() (TaskeeperDB, error) {
	db, err := sql.Open("mysql", "root:MyLocalDbForFun@tcp(127.0.0.1:3306)/taskeeper")
	if err != nil {
		return TaskeeperDB{}, err
	}
	// store the db connection into the CaterDB struct
	glog.Info("now pinging the database")
	err = db.Ping()
	if err != nil {
		return TaskeeperDB{}, err
	}

	glog.Info("sucessfully ping the database :) ")
	return TaskeeperDB{db: db}, nil
}

// to close the connection
func (d TaskeeperDB) Close() {
	d.db.Close()
}

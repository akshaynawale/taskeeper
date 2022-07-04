package taskdb

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
)

type TaskeeperDB struct {
	db *sql.DB // mysql db connection
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Panicf("%s environment variable not set.", k)
	}
	return v
}

// ConnectDb function creates database connection and returns CaterDB struct
func ConnectDb() (TaskeeperDB, error) {
	// get env variables from the environment variables
	var (
		connectionName = mustGetenv("CLOUDSQL_CONNECTION_NAME")
		userName       = mustGetenv("SQL_USERNAME")
		cloudSQLPass   = os.Getenv("SQL_PASSWORD")
		dbNameCSQL     = os.Getenv("CLOUDSQL_DBNAME")
		//socket         = os.Getenv("CLOUDSQL_SOKET_PREFIX")
	)

	//dbURI := fmt.Sprintf("%s:%s@unix(%s/%s)/%s", userName, cloudSQLPass, socket, connectionName, dbNameCSQL)
	dbURI := fmt.Sprintf("%s:%s@unix(/cloudsql/%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", userName, cloudSQLPass, connectionName, dbNameCSQL)
	db, err := sql.Open("mysql", dbURI)
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

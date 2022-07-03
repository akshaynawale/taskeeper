package main

import (
	"fmt"
	tdb "taskdb"

	"github.com/golang/glog"
	"github.com/gomodule/redigo/redis"
)

type Server struct {
	Cache  redis.Conn // redis connection
	TaskDB *tdb.TaskeeperDB
}

func (s *Server) Init() error {
	var err error
	var conn redis.Conn
	// Initialize the redis connection to a redis instance running on your local machine
	if conn, err = redis.DialURL("redis://localhost"); err != nil {
		return fmt.Errorf("failed to connect redis server: %v", err)
	}
	// Assign the connection to the package level `cache` variable
	s.Cache = conn

	// connect to the mySql database
	var tDB tdb.TaskeeperDB
	if tDB, err = tdb.ConnectDb(); err != nil {
		return fmt.Errorf("failed to connect to the db: %v", err)
	}
	s.TaskDB = &tDB
	glog.Infof("successfully connected to the db: %v", tDB)

	return nil
}

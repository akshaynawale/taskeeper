package main

import (
	"fmt"
	"os"

	tdb "github.com/akshaynawale/taskeeper/taskdb"

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
	redisHost := os.Getenv("REDISHOST")
	redisPort := os.Getenv("REDISPORT")
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)
	// Initialize the redis connection to a redis instance running on your local machine
	if conn, err = redis.DialURL(redisAddr); err != nil {
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

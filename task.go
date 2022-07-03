package main

import (
	"net/http"

	"strconv"
	"text/template"

	"github.com/akshaynawale/taskeeper/taskdb"
	"github.com/golang/glog"
	"github.com/gomodule/redigo/redis"
)

type HomeTaskData struct {
	Tasks  []taskdb.Task
	Errors string
}

type HomeHandler struct {
	db    *taskdb.TaskeeperDB
	cache redis.Conn
}

func (hh HomeHandler) GetCacheConn() redis.Conn {
	return hh.cache
}

func (hh HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get the tasks for the user from the db
	// find the User
	var userID int
	var uID string
	var err error
	if uID = r.Header.Get("taskuserid"); uID == "" {
		glog.Errorf("failed to get userid from the http request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if userID, err = strconv.Atoi(uID); err != nil {
		glog.Errorf("failed to get userid from the http request, conversion failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// find the tasks associated with the user
	// convert string to integer
	var tasks []taskdb.Task
	tasks, err = taskdb.GetTasks(hh.db, userID)

	// fill the task info into the page
	var hometemp *template.Template
	if hometemp, err = template.ParseFiles("html/home2.html"); err != nil {
		glog.Errorf("failed to parse html template html/home2.html: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	htd := HomeTaskData{Tasks: tasks, Errors: ""}

	if err = hometemp.Execute(w, htd); err != nil {
		glog.Errorf("failed to write to http response writer: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}

type SaveTaskHandler struct {
	db    *taskdb.TaskeeperDB
	cache redis.Conn
}

func (sth SaveTaskHandler) GetCacheConn() redis.Conn {
	return sth.cache
}

func (sth SaveTaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	glog.Infof("inside the SaveTaskHandler")
	// get the userid from the request
	var uID string
	var userID int
	var err error
	if uID = r.Header.Get("taskuserid"); uID == "" {
		glog.Errorf("failed to get userid from the http request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	glog.Infof("taskuserid inside SaveTaskHandler: %s", uID)
	if userID, err = strconv.Atoi(uID); err != nil {
		glog.Errorf("failed to get userid from the http request, conversion failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Get the TaskTitle and body
	title := r.FormValue("tasktitle")
	body := r.FormValue("taskbody")
	// check if title and body is not empty
	if title == "" || body == "" {
		glog.Errorf("either title or body is empty for the task")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// save the task as a new task in the db
	if err = taskdb.CreateTask(sth.db, userID, title, body); err != nil {
		glog.Errorf("failed to write task into the db: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// get the new tasks from the db
	var tasks []taskdb.Task
	tasks, err = taskdb.GetTasks(sth.db, userID)
	// fill the task info into the page
	var hometemp *template.Template
	if hometemp, err = template.ParseFiles("html/home2.html"); err != nil {
		glog.Errorf("failed to parse html template html/home2.html: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	htd := HomeTaskData{Tasks: tasks, Errors: ""}

	if err = hometemp.Execute(w, htd); err != nil {
		glog.Errorf("failed to write to http response writer: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}

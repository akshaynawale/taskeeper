package main

import (
	"fmt"
	"net/http"

	"github.com/akshaynawale/taskeeper/taskdb"

	"github.com/golang/glog"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserHandler struct {
	db *taskdb.TaskeeperDB
}

func (h CreateUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	glog.Info("inside CreateUserHandler")

	// store all the request data into User struct
	u := taskdb.User{}
	glog.Infof("all values: %+v", r.Form)
	u.Username = r.FormValue("cuname")
	glog.Infof("username: %s", u.Username)

	// get password string and then generate hash for it
	pass := r.FormValue("cpasswd")
	hashedpass, err := bcrypt.GenerateFromPassword([]byte(pass), 8)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		glog.Error("failed to encrypt password")
		return
	}
	u.Password = fmt.Sprintf("%s", hashedpass)
	glog.Infof("password: %s", u.Password)

	// check if the user already exsits
	db_user, err := taskdb.GetUserByName(h.db, u.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If user already exsits
	if db_user != nil {
		fmt.Fprintf(w, fmt.Sprintf("user already exsits with Username: %s try with different username", u.Username))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// If user is not present then add it to the database
	err = taskdb.CreateUser(h.db, u.Username, u.Password)
	if err != nil {
		glog.Infof("failed to create user for username: %s with error: %v", u.Username, err)
		fmt.Fprintf(w, fmt.Sprintf("failed to create user Username: %s with error: %v", u.Username, err))
		return
	}
	// If we reach here means user got created
	fmt.Fprintf(w, fmt.Sprintf("created user sucessfully: %s now you can login with your username", u.Username))
}

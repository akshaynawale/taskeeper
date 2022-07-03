package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	//"taskdb"
	"time"

	"github.com/akshaynawale/taskeeper/taskdb"
	"github.com/golang/glog"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const SessionExpTime = 300

// check authentication
func Auth(r *http.Request) (bool, error) {
	// implement this
	return true, nil
}

// LoginHandler will be used to login a user
type LoginHandler struct {
	db    *taskdb.TaskeeperDB
	cache redis.Conn
}

func (h LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// implement this
	u := &taskdb.User{}
	glog.Infof("all values from the request: %+v", r.Form)
	u.Username = r.FormValue("uname")
	u.Password = r.FormValue("passwd")
	dbUser, err := taskdb.GetUserByName(h.db, u.Username)
	if err != nil || dbUser == nil {
		errMsg := fmt.Sprintf("failed to get user from db with username: %s error: %v", u.Username, err)
		glog.Error(errMsg)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(u.Password)); err != nil {
		errMsg := fmt.Sprintf("wrong password for user: %s try again with correct username/password", u.Username)
		glog.Error(errMsg)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	exp := time.Now().Add(5 * time.Minute)
	sID := fmt.Sprintf("%s", uuid.New())

	// store the session id into redis
	// Set the token in the cache, along with the user whom it represents
	// The token has an expiry time of 300 seconds
	_, err = h.cache.Do("SETEX", sID, "300", dbUser.UserID)
	if err != nil {
		// If there is an error in setting the cache, return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// send the session id as cookie to the browser
	c1 := http.Cookie{Name: "session-id", Value: sID, Expires: exp}
	http.SetCookie(w, &c1)
	hh := HomeHandler{db: h.db, cache: h.cache}
	glog.Infof("inside the login handler setting the request header uid to %d", dbUser.UserID)
	r.Header.Set("taskuserid", fmt.Sprintf("%d", dbUser.UserID))
	hh.ServeHTTP(w, r)
	return
}

type TKHandler interface {
	http.Handler
	GetCacheConn() redis.Conn
}

// CheckAuth is a middleware used to authenticate the http request
func CheckAuth(next TKHandler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get the session id from the cookie in the request
		c, err := r.Cookie("session-id")
		if err != nil {
			if err == http.ErrNoCookie {
				glog.Error("auth failed: no session id found in request cookie")
			}
			glog.Errorf("unknown error occured when trying to check for session id in cookie: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sID := c.Value

		conn := next.GetCacheConn()
		// We then get the name of the user from our cache, where we set the session ID
		userID, err := redis.Int(conn.Do("GET", sID))
		if err != nil {
			// If there is an error fetching from cache, return an internal server error status
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if userID == 0 {
			// If the session token is not present in cache, return an unauthorized error
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		/*
			var userIDInt int
			if userIDInt, err = strconv.Atoi(userID); err != nil {
				glog.Errorf("failed to convert userid to int: %v", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		*/
		// Renew the expiry timer in redis
		_, err = conn.Do("EXPIRE", sID, SessionExpTime)
		if err != nil {
			glog.Errorf("failed to extend expiry time for session id in redis: %v", err)
			w.WriteHeader(http.StatusGatewayTimeout)
			w.Write([]byte(" session timeout please login again"))
			return
		}

		// set new timer in cookie
		// send the session id as cookie to the browser
		c.Expires = time.Now().Add(3 * time.Minute)
		http.SetCookie(w, c)
		// also set the username in the http request
		glog.Infof("request Header previous- %+v", r.Header)
		glog.Infof("Inside Check auth putting the userID as: %v", userID)
		r.Header.Set("taskuserid", fmt.Sprintf("%d", userID))
		glog.Infof("request Header - %+v", r.Header)
		next.ServeHTTP(w, r)
		return
	})

}

// ShowLoginPage shows the login/signup page
func ShowLoginPage(w http.ResponseWriter, r *http.Request) {
	var err error
	var data []byte
	// read the page
	if data, err = ioutil.ReadFile("html/index.html"); err != nil {
		glog.Errorf("failed to read file at html/index.html : %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err = w.Write(data); err != nil {
		glog.Errorf("failed to write to http response writer: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

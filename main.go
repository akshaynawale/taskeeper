package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	var err error
	s := Server{}
	if err = s.Init(); err != nil {
		glog.Errorf("failed to initiate Server: %v", err)
		return
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/", ShowLoginPage)
	// login handler
	mux.Handle("/login", LoginHandler{db: s.TaskDB, cache: s.Cache})
	mux.Handle("/signUp", CreateUserHandler{db: s.TaskDB})
	mux.Handle("/home", CheckAuth(HomeHandler{db: s.TaskDB, cache: s.Cache}))
	mux.Handle("/saveTask", CheckAuth(SaveTaskHandler{db: s.TaskDB, cache: s.Cache}))

	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      mux, // Pass our instance of gorilla/mux in.
	}

	glog.Fatal(srv.ListenAndServe())
}

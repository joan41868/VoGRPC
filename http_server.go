package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"sync"
)

var logger = log.New(os.Stdout, "[server]: ", 0)

type server struct {
	router   *mux.Router
	wsServer *wsServer
}

// NewServer constructs a new server
func NewServer(r *mux.Router) *server {
	return &server{
		router:   r,
		wsServer: NewWSServer(),
	}
}

// Init initializes the endpoints for the backend
func (s *server) Init() {
	router.Handle("/socket.io/", s.wsServer.GetPathHandler())
	s.wsServer.GetPathHandler().EnableCORS("*")
	s.wsServer.Init()
	/* static file server should be initialized last */
	//fs := http.FileServer(http.Dir("./static/"))
	//router.PathPrefix("/").Handler(http.StripPrefix("/", fs))
}

// Start works with a sync.WaitGroup in order to support concurrency in later releases
func (s *server) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	// start the server
	if port := os.Getenv("PORT"); port == "" {
		os.Setenv("PORT", "50515")
	}
	logger.Println("Listening on : " + os.Getenv("PORT"))
	err := http.ListenAndServe(":"+os.Getenv("PORT"), router)
	if err != nil {
		logger.Fatalln(err)
		return
	}
	wg.Done()
}

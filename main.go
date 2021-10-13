package main

import (
	"github.com/gorilla/mux"
	"sync"
)

var router *mux.Router

func main() {
	router = mux.NewRouter()
	wg := new(sync.WaitGroup)

	hs := NewServer(router)
	hs.Init()

	hs.Start(wg)
}

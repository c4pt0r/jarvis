package main

import (
	"flag"
	"net/http"

	"github.com/ngaut/log"
)

var (
	SessionManager *SessionMgr
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()

	// init session manager, set global session manager
	SessionManager = NewSessionMgr()
	go SessionManager.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r)
	})
	log.Info("Websocket Server is running on", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

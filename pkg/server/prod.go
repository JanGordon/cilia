package server

import (
	"fmt"
	"net/http"
)

func Prod(port int) {
	go fileWatcher()
	server := &http.Server{Addr: fmt.Sprintf(":%v", port)}
	http.HandleFunc("/", handler)
	done := make(chan bool)
	go func() {
		server.ListenAndServe()
	}()
	fmt.Printf("ready to connect at http://localhost:%v\n", port)

	<-done
}

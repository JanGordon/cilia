package server

import (
	"fmt"
	"net/http"
)

func Prod(port int) {
	server := &http.Server{Addr: fmt.Sprintf(":%v", port)}
	http.HandleFunc("/", handler)
	done := make(chan bool)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()
	fmt.Printf("ready to connect at http://localhost:%v\n", port)

	<-done
}

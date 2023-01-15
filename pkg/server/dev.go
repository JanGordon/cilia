package server

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"

	"github.com/JanGordon/cilia-framework/pkg/global"
	"github.com/JanGordon/cilia-framework/pkg/ssr"
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
)

var reloadIndicator = make(chan string)
var upgrader = websocket.Upgrader{}
var wConn = 0
var reloadCount = 0

func Dev(port int) {
	go fileWatcher()
	server := &http.Server{Addr: fmt.Sprintf(":%v", port)}
	http.HandleFunc("/ws", wsUpgrader)
	http.HandleFunc("/", handler)
	done := make(chan bool)
	go func() {
		server.ListenAndServe()
		done <- true
	}()
	fmt.Printf("ready to connect at http://localhost:%v\n", port)

	<-done
}

func handler(w http.ResponseWriter, r *http.Request) {
	p := "." + r.URL.Path
	http.ServeFile(w, r, p)
}

func fileWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	filepath.WalkDir(global.ProjectRoot, func(path string, info fs.DirEntry, err error) error {
		watcher.Add(path)
		return nil
	})
	defer watcher.Close()

	//server

	done := make(chan bool)
	go func() {
		defer close(done)
		for {

			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// fmt.Println(event.Name, event.Op)
				if filepath.Ext(event.Name) != ".out" {
					// if html only that should be reloaded on page (with js)
					if filepath.Ext(event.Name) == ".html" {
						ssr.Compile()
						reloadIndicator <- "reloadhtml"
						fmt.Println("Realoding html page")

					} else {
						fmt.Println("Realoding entire page")
						ssr.Compile()
						reloadIndicator <- "reload"
					}

				}
				// dir, filename := filepath.Split(path)
				// if filepath.Ext(path) == ".html" && filename != "out.html" {
				// 	fmt.Println("Rebuilding", path)
				// 	compile.BuildPage(compile.ReplaceComponentWithHTML(path), dir+"out.html", false)
				// }
				//
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}

	}()

	<-done
}

func wsUpgrader(w http.ResponseWriter, r *http.Request) {
	// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade failed: ", err)
		return
	}
	defer conn.Close()
	defer func() { wConn = 0 }()
	// Continuosly read and write message
	for {
		r := <-reloadIndicator
		if r == "reload" || r == "reloadhtml" {
			// reloadIndicator <- false
			message := []byte(r)
			fmt.Println(r)
			err = conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("write failed:", err)
				break
			}
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}
			if string(message) == "reload successful" {
				reloadCount++
				// fmt.Printf("\033[0;0H")
				fmt.Printf("built and reloaded successfully x%v\n", reloadCount)
			} else {
				fmt.Println("client reload failed")
				break
			}
		} else {

		}

	}
	fmt.Println("client disconnected")
}

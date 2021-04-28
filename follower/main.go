package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

// Echo the data received on the WebSocket.
func EchoServer(ws *websocket.Conn) {

	for {
		fmt.Printf("Writing data to ws...")
		ws.Write([]byte("Hello"))
		time.Sleep(1 * time.Second)
	}
	//io.Copy(ws, ws)
}

func main() {

	log.Println("Starting webserver")
	// Serve static files...
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	http.Handle("/ws", websocket.Handler(EchoServer))

	log.Println("Listening on :3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

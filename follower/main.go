package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

type shipdata struct {
	lat float64
	lon float64
}

var ships = make(map[string]*shipdata)

func ShipServer(ws *websocket.Conn) {

	for {
		// Send all ships...
		for id, ship := range ships {
			data := fmt.Sprintf("%s %f %f", id, ship.lat, ship.lon)
			fmt.Printf("SEND %s\n", data)
			ws.Write([]byte(data))
		}

		// Move all ships...
		MAX_MVT := .1
		for _, ship := range ships {
			ship.lat += rand.Float64()*MAX_MVT - rand.Float64()*MAX_MVT
			ship.lon += rand.Float64()*MAX_MVT - rand.Float64()*MAX_MVT

		}

		fmt.Printf("Waiting...")
		time.Sleep(1 * time.Second)
	}
}

func main() {

	lat := 52.214016
	lon := 0.964676

	NUM_SHIPS := 12
	MAX_V := 6.0

	for i := 0; i < NUM_SHIPS; i++ {
		// Create some lat/lon
		id := fmt.Sprintf("sh%d", i)
		ships[id] = &shipdata{
			lat: lat + rand.Float64()*MAX_V - rand.Float64()*MAX_V,
			lon: lon + rand.Float64()*MAX_V - rand.Float64()*MAX_V,
		}
	}

	log.Println("Starting webserver")
	// Serve static files...
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	http.Handle("/ws", websocket.Handler(ShipServer))

	log.Println("Listening on :3000...")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

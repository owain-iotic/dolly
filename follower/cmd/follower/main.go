package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

type shipdata struct {
	lat float64
	lon float64
}

var ships = make(map[string]*shipdata)

var status_index = 0
var statuses = []string{
	"/",
	"-",
	"\\",
	"|",
}

func ShipServer(ws *websocket.Conn) {

	// Read the data from file...
	file, err := os.Open("ship_data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		csvline := scanner.Text()
		// Split it
		bits := strings.Split(csvline, ",")
		// MMSI,BaseDateTime,LAT,LON,SOG,COG,Heading,VesselName,IMO,CallSign,VesselType,Status,Length,Width,Draft,Cargo,TranscieverClass
		lat, _ := strconv.ParseFloat(bits[2], 64)
		lon, _ := strconv.ParseFloat(bits[3], 64)
		id := bits[9]
		timestamp := bits[1]

		// Send an update...

		status_index++
		if status_index == len(statuses) {
			status_index = 0
		}

		status := fmt.Sprintf("IOTICS %s Epic OMG Ship Example Replaying %s", statuses[status_index], timestamp)

		data := fmt.Sprintf("%s,%f,%f,%s", id, lat, lon, status)
		fmt.Printf("SEND %s %s\n", timestamp, data)
		ws.Write([]byte(data))

		time.Sleep(100 * time.Millisecond)
	}
	/*
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
	*/
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

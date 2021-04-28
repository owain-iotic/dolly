package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type shipdata struct {
	min_lon float64
	max_lon float64
	min_lat float64
	max_lat float64
}

func main() {

	var infile string

	flag.StringVar(&infile, "file", "data.csv", "Data file")

	flag.Parse()

	// Read the data from file, and do some analysis so we can pick good ships

	ships := make(map[string]*shipdata)

	file, err := os.Open(infile)
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

		ship, ok := ships[id]
		if ok {
			ship.min_lon = math.Min(ship.min_lon, lon)
			ship.max_lon = math.Max(ship.max_lon, lon)
			ship.min_lat = math.Min(ship.min_lat, lat)
			ship.max_lat = math.Max(ship.max_lat, lat)
		} else {
			ships[id] = &shipdata{
				min_lon: lon,
				max_lon: lon,
				min_lat: lat,
				max_lat: lat,
			}
		}
	}

	// Output stats...

	for id, ship := range ships {
		diff_lon := ship.max_lon - ship.min_lon
		diff_lat := ship.max_lat - ship.min_lat
		dist := diff_lon + diff_lat // Not correct, but good enough for jazz
		fmt.Printf("SHIP,%s,%f\n", id, dist)
	}
}

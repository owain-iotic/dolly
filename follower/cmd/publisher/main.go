package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/owain-iotic/dolly/follower/client"
	"github.com/owain-iotic/dolly/follower/common"
)

const (
	ssl  = true
	host = "plateng.iotics.space"

	authToken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJFUzI1NiJ9.eyJpc3MiOiJkaWQ6aW90aWNzOmlvdEhCQ21wUHZUUVJySndXZFhNNTZhMTltclhLd0g0NmZGTCNhZ2VudC0wIiwiYXVkIjoiaHR0cHM6Ly9kaWQucHJkLmlvdGljcy5jb20iLCJzdWIiOiJkaWQ6aW90aWNzOmlvdENkdWpWQ3ZCNllQQ1JGa1VNTnpjSnVNMVdkUUZhcHBpVyIsImlhdCI6MTYxOTY4MzA3MiwiZXhwIjoxNjE5NzExOTAyfQ.zdiXHK39scpHJjwL3EOeSKGMtjroculC6XemPmjWLZ5KBtS_X2kLfwEXiRF_43zm9DHeB0K-oz4PrPxIrzc4iw"
)

func main() {

	var speed float64

	flag.Float64Var(&speed, "speed", 1000, "Speed to play the events")

	flag.Parse()

	common.LoadShipConfig("selected_ships.txt", "twins.txt")

	scheme := "ws"
	if ssl {
		scheme = "wss"
	}

	url := fmt.Sprintf("%s://%s/ws", scheme, host)

	cli := client.NewIoticsStompClient()

	err := cli.Connect(url, authToken)
	if err != nil {
		panic(err)
	}

	replay_start_date, err := time.Parse(time.RFC3339, "2020-06-01T00:05:19+00:00")

	start_time := time.Now()

	if err != nil {
		panic(err)
	}

	// Do publishes to the ship twins...
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

		date, err := time.Parse(time.RFC3339, fmt.Sprintf("%s+00:00", timestamp))
		if err != nil {
			panic(err)
		}

		diff := date.Sub(replay_start_date)

		for {
			diff_replay := time.Since(start_time) * time.Duration(speed)
			if diff_replay > diff {
				break
			}
			time.Sleep(50 * time.Millisecond)
		}

		twinId, ok := common.Shiptwins[id]
		if !ok {
			panic("Don't have a twin for this ship")
		}

		s := fmt.Sprintf("{\"lat\": %f, \"lon\":%f, \"when\":\"%s\"}", lat, lon, timestamp)
		b64val := base64.StdEncoding.EncodeToString([]byte(s))
		data := fmt.Sprintf("{\"sample\": {\"data\": \"%s\", \"mime\": \"application/json\", \"occurredAt\": \"%s\"}}", b64val, time.Now().Format(time.RFC3339Nano))

		fmt.Printf("Posting update %s...\n", data)
		err = cli.PostFeedData(twinId, "shiplocation", data)

		if err != nil {
			panic(err)
		}
		// Publish it...
		fmt.Printf("DATA %s\n", csvline)
	}
}

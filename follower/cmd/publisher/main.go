package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/owain-iotic/dolly/follower/client"
)

const (
	ssl  = true
	host = "plateng.iotics.space"

	authToken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJFUzI1NiJ9.eyJpc3MiOiJkaWQ6aW90aWNzOmlvdEhCQ21wUHZUUVJySndXZFhNNTZhMTltclhLd0g0NmZGTCNhZ2VudC0wIiwiYXVkIjoiaHR0cHM6Ly9kaWQucHJkLmlvdGljcy5jb20iLCJzdWIiOiJkaWQ6aW90aWNzOmlvdENkdWpWQ3ZCNllQQ1JGa1VNTnpjSnVNMVdkUUZhcHBpVyIsImlhdCI6MTYxOTYxOTg1OSwiZXhwIjoxNjE5NjQ4Njg5fQ.cIxLrZ6vdcMPIdRlA6oHa-FBwC8yvIO1-R0oiFUo2gczWH4NHLu07DBj3rnY2hJi0HByTlZxsV0q_W8cL-6SnQ"
)

var shiptwins = map[string]string{
	"3FHS3":   "did:iotics:iotQYhr2yNaMoAT16uuhKRDR5yUoLMVjQ8jw",
	"3FMK3":   "did:iotics:iotT94QxhtDvaw4gzd1qhymbNcjG3CwCwbtr",
	"9HA3481": "did:iotics:iotUSE8zLAo5VnqXar5WM8VnVE9TjVKVnEGL",
	"9VPY3":   "did:iotics:iotFLjrTjAVFnxnBBmzLkd2VqEc9EPg9kXtc",
	"C6FX6":   "did:iotics:iotStFerXzr6FNVjtk4KFVpwzMszd6fCpSWD",
	"C6QM8":   "did:iotics:iotMEe8p4de64aSGsnC7cDtfkdMek2fpFGUP",
	"C6WW7":   "did:iotics:iotEr28AzvJ9XijVrJTUQC6fM4DqbKxgpsAk",
	"D5MJ2":   "did:iotics:iotJ81ZqgWDeYUpy6yV3UtC81wuaK62AyDC3",
	"D5WA5":   "did:iotics:iotNtKY6kpdWoDCymGBon85eUE1NFKqtE7wt",
	"KEBO":    "did:iotics:iotSdaaoRwX5WX9RC8u3KBQsBX7xrSs6bEJf",
	"LAHY6":   "did:iotics:iotDLBijG1znV7B5XBMz2rA17UyKpgy2KcG5",
	"VRGA3":   "did:iotics:iotJUa9oH6rYdc1jZpZfKJuxmuZQcsGBUgzF",
}

func main() {
	scheme := "ws"
	if ssl {
		scheme = "wss"
	}

	url := fmt.Sprintf("%s://%s/ws", scheme, host)

	cli := client.NewIoticsStompClient()

	cli.Connect(url, authToken)

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

	REPLAY_SPEED := 1000

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
			diff_replay := time.Since(start_time) * time.Duration(REPLAY_SPEED)
			if diff_replay > diff {
				break
			}
			time.Sleep(50 * time.Millisecond)
		}

		twinId, ok := shiptwins[id]
		if !ok {
			panic("Don't have a twin for this ship")
		}

		s := fmt.Sprintf("{\"lat\": %f, \"lon\":%f, \"when\":\"%s\"}", lat, lon, timestamp)
		b64val := base64.StdEncoding.EncodeToString([]byte(s))
		data := fmt.Sprintf("{\"sample\": {\"data\": \"%s\", \"mime\": \"application/json\", \"occurredAt\": \"%s\"}}", b64val, time.Now().Format(time.RFC3339Nano))

		fmt.Printf("Posting update %s...\n", data)
		cli.PostFeedData(twinId, "shiplocation", data)

		// Publish it...
		fmt.Printf("DATA %s\n", csvline)
	}
}

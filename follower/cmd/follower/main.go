package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/owain-iotic/dolly/follower/client"
	"golang.org/x/net/websocket"
)

const (
	ssl  = true
	host = "plateng.iotics.space"

	authToken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJFUzI1NiJ9.eyJpc3MiOiJkaWQ6aW90aWNzOmlvdEhCQ21wUHZUUVJySndXZFhNNTZhMTltclhLd0g0NmZGTCNhZ2VudC0wIiwiYXVkIjoiaHR0cHM6Ly9kaWQucHJkLmlvdGljcy5jb20iLCJzdWIiOiJkaWQ6aW90aWNzOmlvdENkdWpWQ3ZCNllQQ1JGa1VNTnpjSnVNMVdkUUZhcHBpVyIsImlhdCI6MTYxOTYxOTg1OSwiZXhwIjoxNjE5NjQ4Njg5fQ.cIxLrZ6vdcMPIdRlA6oHa-FBwC8yvIO1-R0oiFUo2gczWH4NHLu07DBj3rnY2hJi0HByTlZxsV0q_W8cL-6SnQ"

	followerTwinId = "did:iotics:iotTmRTTzh9LGuqPNkZgjQ3Pj6w8fBfovxfJ"
)

var status_index = 0
var statuses = []string{
	"/",
	"-",
	"\\",
	"|",
}

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

type shipdata struct {
	lat  float64
	lon  float64
	id   string
	did  string
	when string
}

// Channel to write to client UI
var updates = make(chan *shipdata, 1024)

func ShipServer(ws *websocket.Conn) {

	for {
		shipdata := <-updates
		status_index++
		if status_index == len(statuses) {
			status_index = 0
		}

		status := fmt.Sprintf("IOTICS %s Epic OMG Ship Example Replaying %s", statuses[status_index], shipdata.when)

		data := fmt.Sprintf("%s,%f,%f,%s", shipdata.id+"-"+shipdata.did, shipdata.lat, shipdata.lon, status)
		fmt.Printf("SEND %s %s\n", shipdata.when, data)
		ws.Write([]byte(data))

		time.Sleep(100 * time.Millisecond)
	}
}

func main() {

	scheme := "ws"
	if ssl {
		scheme = "wss"
	}

	url := fmt.Sprintf("%s://%s/ws", scheme, host)

	cli := client.NewIoticsStompClient()

	cli.Connect(url, authToken)

	// Subscribe to the twins...
	for id, did := range shiptwins {
		dest := fmt.Sprintf("/qapi/twins/%s/interests/twins/%s/feeds/%s", followerTwinId, did, "shiplocation")

		ch, err := cli.Subscribe(dest)
		if err != nil {
			panic(err)
		}

		go func(myid string, mydid string) {
			for {
				m, ok := <-ch

				if !ok {
					panic("Error on sub")
				}

				// Get the data...
				var result map[string]interface{}
				json.Unmarshal(m.Body, &result)

				// Now find what we want...
				feedData := result["feedData"].(map[string]interface{})
				dp := feedData["data"].(string)
				val, _ := base64.StdEncoding.DecodeString(dp)

				type location struct {
					Lat  float64 `json:"lat"`
					Lon  float64 `json:"lon"`
					When string  `json:"when"`
				}

				var loc location
				json.Unmarshal(val, &loc)

				// Now we have loc

				fmt.Printf("MESSAGE %s %s %f %f %s\n", myid, mydid, loc.Lat, loc.Lon, loc.When)

				ship := &shipdata{
					lat:  loc.Lat,
					lon:  loc.Lon,
					id:   myid,
					did:  mydid,
					when: loc.When,
				}

				updates <- ship

			}

		}(id, did)
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

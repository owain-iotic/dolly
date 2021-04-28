package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/owain-iotic/dolly/follower/client"
	"github.com/owain-iotic/dolly/follower/common"
	"golang.org/x/net/websocket"
)

const (
	ssl  = true
	host = "plateng.iotics.space"

	authToken = "eyJ0eXAiOiJKV1QiLCJhbGciOiJFUzI1NiJ9.eyJpc3MiOiJkaWQ6aW90aWNzOmlvdEhCQ21wUHZUUVJySndXZFhNNTZhMTltclhLd0g0NmZGTCNhZ2VudC0wIiwiYXVkIjoiaHR0cHM6Ly9kaWQucHJkLmlvdGljcy5jb20iLCJzdWIiOiJkaWQ6aW90aWNzOmlvdENkdWpWQ3ZCNllQQ1JGa1VNTnpjSnVNMVdkUUZhcHBpVyIsImlhdCI6MTYxOTYxOTg1OSwiZXhwIjoxNjE5NjQ4Njg5fQ.cIxLrZ6vdcMPIdRlA6oHa-FBwC8yvIO1-R0oiFUo2gczWH4NHLu07DBj3rnY2hJi0HByTlZxsV0q_W8cL-6SnQ"

	followerTwinId = "did:iotics:iotTmRTTzh9LGuqPNkZgjQ3Pj6w8fBfovxfJ"
)

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

		status := fmt.Sprintf("Replaying %s", shipdata.when)

		data := fmt.Sprintf("%s,%f,%f,%s", shipdata.id+" "+shipdata.did, shipdata.lat, shipdata.lon, status)
		fmt.Printf("SEND %s %s\n", shipdata.when, data)
		ws.Write([]byte(data))
	}
}

func main() {
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

	fmt.Printf("Stomp connected\n")

	// Subscribe to the twins...
	for id, did := range common.Shiptwins {
		dest := fmt.Sprintf("/qapi/twins/%s/interests/twins/%s/feeds/%s", followerTwinId, did, "shiplocation")

		ch, err := cli.Subscribe(dest)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Subscribed to %s\n", dest)

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
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

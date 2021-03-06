package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	//"math/rand"

	log "github.com/sirupsen/logrus"

	//"sync/atomic"
	//"math/rand"
	"os"
	"time"
	// "github.com/go-stomp/stomp/v3"
	// "github.com/google/uuid"
	// "golang.org/x/net/websocket"
)

//{\"visibility\":\"PUBLIC\"}

func main() {
	config := Config{
		Resolver:     "https://did.prd.iotics.com",
		AgentSeed:    "4137d728017dd491d26ad48d64c9a24ab4e8c898292957d63acbf87f7a507dab",
		AgentKeyName: "agent-0",
		AgentDid:     "did:iotics:iotWDjh2FcRfHxwCj7WB8mn2GCoKabVew999",
		Host:         "plateng.iotics.space",
		AuthToken:    "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJub3QtY3VycmVudGx5LXVzZWQuZXhhbXBsZS5jb20iLCJleHAiOjE2MTk2OTc2NzMsImlhdCI6MTYxOTY4NDA3MywiaXNzIjoiZGlkOmlvdGljczppb3RXRGpoMkZjUmZIeHdDajdXQjhtbjJHQ29LYWJWZXc5OTkjYWdlbnQtMCIsInN1YiI6ImRpZDppb3RpY3M6aW90Uml6NmFUeUpCaVJGUkJObWprckthUHBaeHltN0IzUnV0In0.GPu_iFx_GJ1vxkxXxNDey9YPGvnXufLvKJauay-oe6v5wJZ1iVyll-1xVOY99InRr3OSU8AGXj8CFcl1LoMt1w",
	}

	//export DID_USER_SEED=89489442ac079dbab0f4b5c5542b35c7bbd8bea23eb7dcb9e81dab3e0d194832
	//export DID_USER_KEYNAME=fn-host
	//export DID_USER_DID=did:iotics:iotRiz6aTyJBiRFRBNmjkrKaPpZxym7B3Rut

	ship, err := NewShip(config)
	if err != nil {
		panic(err)
	}
	ship.TwinName = "lollypop4"

	ship.Load("did:iotics:iotMBLRg33y6XQwd9pTpNDJxTwX9JMcG7hYb")

	// err = ship.Create()
	// if err != nil {
	// 	panic(err)
	// }
	// ships := host.search()
	// ships := ship.Search("some query param")
	//ship.Init(twinid)

	go ship.FeedAis.Follow(onFeedAisEvent)

	for {
		// ship.FeedAis.Vallabel = rand.Intn(12)
		// ship.FeedAis.Publish()
		time.Sleep(1000 * time.Millisecond)
	}
}

type onEventFeedAis func(feed FeedAis) string

func onFeedAisEvent(x FeedAis) string {
	msg := fmt.Sprintf("val label %v", x.Vallabel)
	fmt.Println(msg)
	return msg
}

type Host struct {
	config Config
	stomp  *IoticsStompClient
}

func NewHost(config Config) Host {
	return Host{
		config: config,
	}
}

func (h *Host) Search() error {
	ssl := true
	scheme := "ws"
	if ssl {
		scheme = "wss"
	}

	url := fmt.Sprintf("%s://%s/ws", scheme, h.config.Host)

	h.stomp = NewIoticsStompClient()

	err := h.stomp.Connect(url, h.config.AuthToken)
	if err != nil {
		return err
	}

	query := "{\"filter\": {\"text\": \"tfl-bikes\"}, \"responseType\": \"FULL\"}"
	id, ch, err := h.stomp.Search("GLOBAL", query)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Search id %s\n", id)

	for {
		select {
		case res := <-ch:
			fmt.Printf("Result: %s\n\n", res.Body)

			twinsResp := &TwinsResp{}

			err = json.Unmarshal(res.Body, &twinsResp)

			fmt.Printf("%v\n", twinsResp)
			fmt.Printf("%d twins...\n", len(twinsResp.Twins))
		}
	}

}

type NewFeed struct {
	Labels    Labels     `json:"labels"`
	Comments  Comments   `json:"comments"`
	Values    FeedValues `json:"values"`
	Tags      FeedTags   `json:"tags"`
	StoreLast bool       `json:"storeLast"`
}

type NewTwin struct {
	NewVisibility Visibility `json:"newVisibility"`
	Location      Location   `json:"location"`
	Labels        Labels     `json:"labels"`
	Comments      Comments   `json:"comments"`
	Tags          Tags       `json:"tags"`
	Properties    Properties `json:"properties"`
	//Feeds         []Feedinfo
}

type FeedTags struct {
	Added []string `json:"added"`
}

type FeedValues struct {
	Added []ValuesAdded `json:"added"`
}

type Visibility struct {
	Visibility string `json:"visibility"`
}

type Location struct {
	Location Latlon `json:"location"`
}

type Comments struct {
	Added []CommentsAdded `json:"added"`
}

type Tags struct {
	Added []string `json:"added"`
}

type Labels struct {
	Added []LabelsAdded `json:"added"`
}

type Properties struct {
	Added []PropertiesAdded `json:"added"`
}

type PropertiesAdded struct {
	Key   string             `json:"key"`
	Value StringLiteralValue `json:"stringLiteralValue"`
}

type ValuesAdded struct {
	Label    string `json:"label"`
	Comment  string `json:"comment"`
	Unit     string `json:"unit"`
	DataType string `json:"dataType"`
}

type StringLiteralValue struct {
	Value string `json:"value"`
}

type TagValue struct {
	Value string `json:"value"`
}

type CommentsAdded struct {
	Lang  string `json:"lang"`
	Value string `json:"value"`
}

type LabelsAdded struct {
	Lang  string `json:"lang"`
	Value string `json:"value"`
}

type TagsAdded struct {
	Tag string
}

type Idvalue struct {
	Value string
}

type Twindata struct {
	Visibility string
	Id         Idvalue
}

type TwinsResp struct {
	Twins []Twindata
}

// For describe twin
type TwinResp struct {
	Twin   Twindata
	Result Twinresult
}

type Latlon struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Label struct {
	Lang  string `json:"lang"`
	Value string `json:"value"`
}

type Feedinfo struct {
	FeedId    Idvalue
	StoreLast bool
	Labels    []interface{}
}

type Propinfo struct {
	Key                string
	StringLiteralValue Idvalue
}

type Twinresult struct {
	Location   Latlon
	Labels     []Label
	Comments   []Label
	Feeds      []Feedinfo
	Tags       []string
	Properties []Propinfo
}

type FeedResp struct {
	Feed   Feeddata
	Result FeedResult `json:result`
}

// For describe twin
type Feeddata struct {
	ID     Idvalue `json:id`
	TwinID Idvalue `json:twinID`
}

type Values struct {
	Label    string
	Comment  string
	unit     string
	dataType string
}

type FeedResult struct {
	Labels    []Label
	Comments  []Label
	Values    []Values
	Tags      []string
	StoreLast bool
}

type Ship struct {
	TwinName string
	TwinDid  string
	config   Config
	did      *Did
	FeedAis  FeedAis
}

type Config struct {
	Resolver     string
	AgentSeed    string
	AgentKeyName string
	AgentDid     string
	Host         string
	AuthToken    string
}

func NewShip(config Config) (*Ship, error) {
	os.Setenv("RESOLVER", config.Resolver)
	var rtn Ship
	did, err := NewDidFromConfig(config)
	if err != nil {
		return &rtn, err
	}
	rtn.config = config
	rtn.did = did
	rtn.FeedAis, err = NewFeedAis(&rtn.TwinDid, config)
	if err != nil {
		return &rtn, err
	}
	return &rtn, nil
}

func (s *Ship) Load(twinID string) error {
	s.TwinDid = twinID
	feedAis, err := NewFeedAis(&s.TwinDid, s.config)
	if err != nil {
		return err
	}

	s.FeedAis = feedAis

	return nil
}

func (s *Ship) valid() error {
	if s.TwinName == "" {
		return fmt.Errorf("Must set the twin name")
	}
	return nil
}

func (s *Ship) createTwinDid() error {
	didBytes, err := s.did.CreateTwinDid(s.TwinName)
	s.TwinDid = string(didBytes)
	if err != nil {
		return err
	}
	log.Info(string(s.TwinDid))
	return nil
}

func (s *Ship) createTwin() error {

	httpClient := NewHttpClient(s.config.Host, true, s.config.AuthToken)

	didBytes, err := s.did.CreateTwinDid(s.TwinName)
	s.TwinDid = string(didBytes)
	if err != nil {
		return err
	}
	log.Info(string(s.TwinDid))

	err = httpClient.CreateTwin(s.TwinDid)
	if err != nil {
		return err
	}

	return nil
}

func (s *Ship) createFeed() error {
	httpClient := NewHttpClient(s.config.Host, true, s.config.AuthToken)
	feedId := "ais"
	storeLast := true
	err := httpClient.CreateFeed(s.TwinDid, feedId, storeLast)
	if err != nil {
		return err
	}
	return nil
}

func (s *Ship) getTwinData() NewTwin {
	rtn := NewTwin{
		NewVisibility: Visibility{
			"PUBLIC",
		},
		Comments: Comments{
			Added: []CommentsAdded{
				{
					Lang:  "en",
					Value: "comment value",
				},
			},
		},
		Properties: Properties{
			Added: []PropertiesAdded{
				{
					Key: "http://data.iotics.com/Ship",
					Value: StringLiteralValue{
						"http://data.iotics.com/Shippymcshipface",
					},
				},
			},
		},
		Labels: Labels{
			Added: []LabelsAdded{
				{
					Lang:  "en",
					Value: "label",
				},
			},
		},
		Location: Location{
			Latlon{
				Lat: 53.34,
				Lon: -6.2603,
			},
		},
		Tags: Tags{
			[]string{
				"cat_weather",
			},
		},
	}
	return rtn
}

func (s *Ship) updateTwin() error {

	httpClient := NewHttpClient(s.config.Host, true, s.config.AuthToken)

	log.Info(string(s.TwinDid))

	err := httpClient.UpdateTwin(s.TwinDid, s.getTwinData())
	if err != nil {
		return err
	}

	return nil
}

func (s *Ship) getFeedData() NewFeed {

	//Need to support multiple feeds

	rtn := NewFeed{
		StoreLast: true,
		Values: FeedValues{
			Added: []ValuesAdded{
				{
					Label:    "vallabel",
					Comment:  "some comment or other",
					Unit:     "http://purl.obolibrary.org/obo/bannana",
					DataType: "integer",
				},
			},
		},
		Comments: Comments{
			Added: []CommentsAdded{
				{
					Lang:  "en",
					Value: "feed comment value",
				},
			},
		},
		Labels: Labels{
			Added: []LabelsAdded{
				{
					Lang:  "en",
					Value: "feed label",
				},
			},
		},
		Tags: FeedTags{
			[]string{
				"cat_weather",
			},
		},
	}
	return rtn

}

func (s *Ship) updateFeed(feedID string) error {

	httpClient := NewHttpClient(s.config.Host, true, s.config.AuthToken)

	log.Info(string(s.TwinDid))

	err := httpClient.UpdateFeed(s.TwinDid, feedID, s.getFeedData())
	if err != nil {
		return err
	}

	return nil
}

func (s *Ship) Create() error {
	err := s.valid()
	if err != nil {
		return err
	}

	err = s.createTwinDid()
	if err != nil {
		return err
	}

	err = s.createTwin()
	if err != nil {
		return err
	}

	err = s.updateTwin()
	if err != nil {
		return err
	}

	err = s.createFeed()
	if err != nil {
		return err
	}

	err = s.updateFeed("ais")
	if err != nil {
		return err
	}

	return nil
}

type FeedAis struct {
	config   Config             `json:"-"`
	stomp    *IoticsStompClient `json:"-"`
	TwinID   *string            `json:"-"`
	FeedID   string             `json:"-"`
	Vallabel int
}

func NewFeedAis(twinID *string, config Config) (FeedAis, error) {
	rtn := FeedAis{
		config: config,
		TwinID: twinID,
		FeedID: "ais",
	}

	ssl := true
	scheme := "ws"
	if ssl {
		scheme = "wss"
	}

	url := fmt.Sprintf("%s://%s/ws", scheme, config.Host)
	rtn.stomp = NewIoticsStompClient()
	err := rtn.stomp.Connect(url, config.AuthToken)
	if err != nil {
		return rtn, err
	}
	return rtn, nil
}

func (t *FeedAis) Publish() error {

	twinId := t.TwinID
	feedId := t.FeedID

	data, err := json.Marshal(t)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	b64Json := base64.StdEncoding.EncodeToString(data)
	if err != nil {
		return err
	}

	log.Info(b64Json)

	data1 := fmt.Sprintf("{\"sample\": {\"data\": \"%s\", \"mime\": \"application/json\", \"occurredAt\": \"%s\"}}", b64Json, time.Now().Format(time.RFC3339Nano))
	log.Infof("Posting update twinid: %s feedid %s data: %d...\n", *twinId, feedId, t.Vallabel)
	err = t.stomp.PostFeedData(*twinId, feedId, data1)
	if err != nil {
		return err
	}
	return nil

}

func (t *FeedAis) followerID() string {
	return "did:iotics:iotHL8fQvyXSm85KwSRmeWNhwr8PcLTGsdLY"
}

type followLoc struct {
	twinId string
	feedId string
}

func (t *FeedAis) Follow(fn onEventFeedAis) {
	////
	//locs []followLoc, failures chan followLoc, )
	var failures chan followLoc
	//var stat *int64
	loc := followLoc{
		twinId: *t.TwinID,
		feedId: t.FeedID,
	}
	dest := fmt.Sprintf("/qapi/twins/%s/interests/twins/%s/feeds/%s", t.followerID(), *t.TwinID, t.FeedID)

	ch, err := t.stomp.Subscribe(dest)
	if err != nil {
		panic(err)
	}

	// Read from ch...

	go func(myloc followLoc) {
		// atomic.AddInt64(stat, 1)
		// defer atomic.AddInt64(stat, -1)

		for {

			m, ok := <-ch

			if !ok {
				failures <- myloc
				return
			}

			// Get the data...
			var result map[string]interface{}
			json.Unmarshal(m.Body, &result)

			// Now find what we want...
			feedData := result["feedData"].(map[string]interface{})
			dp := feedData["data"].(string)
			val, _ := base64.StdEncoding.DecodeString(dp)

			if err := json.Unmarshal(val, &t); err != nil {
				panic(err)
			}

			// data := NewData()
			// bits, err := data.FromBase64Json(string(dp))
			// if err != nil {
			// 	fmt.Printf("error getting bits %v\n", err)
			// }

			// t.Vallabel = int(bits["vallabel"].(float64))

			fn(*t)
			//fmt.Printf("MESSAGE %s %s %s\n", myloc.twinId, myloc.feedId, val)
			//time.Sleep(1 * time.Second)
		}
	}(loc)
}

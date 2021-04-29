package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/rand"

	//"math/rand"
	"os"
	"time"
	// "github.com/go-stomp/stomp/v3"
	// "github.com/google/uuid"
	// "golang.org/x/net/websocket"
)

//{\"visibility\":\"PUBLIC\"}

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

type Data struct {
	items map[string]interface{}
}

func NewData() Data {
	return Data{
		items: make(map[string]interface{}),
	}
}

func (d *Data) Add(label string, value interface{}) {
	d.items[label] = value
}

func (d *Data) ToJson() ([]byte, error) {
	rtn, err := json.Marshal(d.items)
	if err != nil {
		return []byte{}, err
	}

	return rtn, nil
}

func (d *Data) ToBase64Json() (string, error) {
	json, err := d.ToJson()
	if err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(json)
	return encoded, nil
}

func (d *Data) FromBase64Json(json string) (map[string]interface{}, error) {
	encoded, err := base64.StdEncoding.DecodeString(json)
	if err != nil {
		return d.items, err
	}
	return d.FromJson(encoded)
}

func (d *Data) FromJson(item []byte) (map[string]interface{}, error) {
	err := json.Unmarshal(item, &d.items)
	if err != nil {
		return d.items, err
	}

	return d.items, nil
}

func main() {
	config := Config{
		Resolver:     "https://did.prd.iotics.com",
		AgentSeed:    "4137d728017dd491d26ad48d64c9a24ab4e8c898292957d63acbf87f7a507dab",
		AgentKeyName: "agent-0",
		AgentDid:     "did:iotics:iotWDjh2FcRfHxwCj7WB8mn2GCoKabVew999",
		Host:         "plateng.iotics.space",
		AuthToken:    "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJub3QtY3VycmVudGx5LXVzZWQuZXhhbXBsZS5jb20iLCJleHAiOjE2MTk2MzE5NTQsImlhdCI6MTYxOTYxODM1NCwiaXNzIjoiZGlkOmlvdGljczppb3RXRGpoMkZjUmZIeHdDajdXQjhtbjJHQ29LYWJWZXc5OTkjYWdlbnQtMCIsInN1YiI6ImRpZDppb3RpY3M6aW90Uml6NmFUeUpCaVJGUkJObWprckthUHBaeHltN0IzUnV0In0.d_bLs1ydSA2pZ2U4IUYcKOjZwJ7gmFQJpRxBV1CMhef3Ofo75eJX2qOR5jyw1zGg5PEjhzTttfbcw4s1qw_zyw",
	}

	//export DID_USER_SEED=89489442ac079dbab0f4b5c5542b35c7bbd8bea23eb7dcb9e81dab3e0d194832
	//export DID_USER_KEYNAME=fn-host
	//export DID_USER_DID=did:iotics:iotRiz6aTyJBiRFRBNmjkrKaPpZxym7B3Rut

	ship, err := NewShip(config)
	if err != nil {
		panic(err)
	}

	ship.TwinName = "lollypop4"

	err = ship.Create()
	if err != nil {
		panic(err)
	}
	
	for {
		ship.FeedAis.Vallabel = rand.Intn(12)
		ship.FeedAis.Publish()
		time.Sleep(1000 * time.Millisecond)
	}
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
	config   Config
	stomp    *IoticsStompClient
	TwinID   *string
	FeedID   string
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

	data := NewData()
	//for _, v := range feed.Data {
	data.Add("vallabel", t.Vallabel)
	//}
	b64Json, err := data.ToBase64Json()
	if err != nil {
		return err
	}

	log.Info(b64Json)

	data1 := fmt.Sprintf("{\"sample\": {\"data\": \"%s\", \"mime\": \"application/json\", \"occurredAt\": \"%s\"}}", b64Json, time.Now().Format(time.RFC3339Nano))
	log.Infof("Posting update twinid: %s feedid %s data: %s...\n", *twinId, feedId, data1)
	err = t.stomp.PostFeedData(*twinId, feedId, data1)
	if err != nil {
		return err
	}
	return nil

}

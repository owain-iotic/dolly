package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-stomp/stomp/v3"
	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

// IoticsStompClient
type IoticsStompClient struct {
	ws            *websocket.Conn
	stomp         *stomp.Conn
	lock          sync.Mutex
	subscriptions map[string]*stomp.Subscription

	clientRef   string
	clientAppId string
}

func NewIoticsStompClient() *IoticsStompClient {
	return &IoticsStompClient{
		subscriptions: make(map[string]*stomp.Subscription),
		clientRef:     uuid.NewString(),
		clientAppId:   uuid.NewString(),
	}
}

// Perform a search...
func (isc *IoticsStompClient) Search(scope string, query string) (string, chan *stomp.Message, error) {
	uu := uuid.NewString()
	ch := make(chan *stomp.Message, 1)

	// Subscribe to search results...

	sresults, err := isc.SubscribeWithId("/qapi/searches/results", uu)
	if err != nil {
		return "", nil, err
	}

	fmt.Println("Subbed to search results...")

	// Search result reader and dispatcher...
	go func() {
		for {
			res, ok := <-sresults
			if !ok {
				return
			} else {
				// Check it's meant for us...
				if res.Header.Get("Iotics-TransactionRef") == uu {
					ch <- res
				}
			}
		}
	}()

	err = isc.stomp.Send("/qapi/searches/dispatches", "", []byte(query),
		stomp.SendOpt.Header("Iotics-ClientAppId", isc.clientAppId),
		stomp.SendOpt.Header("Iotics-ClientRef", isc.clientRef),
		stomp.SendOpt.Header("Iotics-TransactionRef", uu),
		stomp.SendOpt.Header("scope", scope),

		//		stomp.SendOpt.Header("limit", "100"),
		//		stomp.SendOpt.Header("offset", "0"),
	)
	if err != nil {
		return "", nil, err
	}

	return uu, ch, nil
}

func (isc *IoticsStompClient) PostFeedData(twinId string, feedId string, data string) error {
	uu := uuid.NewString()

	dest := fmt.Sprintf("/qapi/twins/%s/feeds/%s", twinId, feedId)
	
	err := isc.stomp.Send(dest, "", []byte(data),
		stomp.SendOpt.Header("Iotics-ClientAppId", isc.clientAppId),
		stomp.SendOpt.Header("Iotics-ClientRef", isc.clientRef),
		stomp.SendOpt.Header("Iotics-TransactionRef", uu),
	)
	return err
}

// Subscribe subscribes to a topic
func (isc *IoticsStompClient) Subscribe(dest string) (chan *stomp.Message, error) {
	uu := uuid.NewString()
	return isc.SubscribeWithId(dest, uu)
}

// Subscribe subscribes to a topic
func (isc *IoticsStompClient) SubscribeWithId(dest string, id string) (chan *stomp.Message, error) {

	subsr, err := isc.stomp.Subscribe(dest, stomp.AckAuto,
		stomp.SubscribeOpt.Header("Iotics-ClientAppId", isc.clientAppId),
		stomp.SubscribeOpt.Header("Iotics-ClientRef", isc.clientRef),
		stomp.SubscribeOpt.Header("Iotics-TransactionRef", id),
		stomp.SubscribeOpt.Id(id),
	)
	if err != nil {
		return nil, err
	}

	// Store it so we know what's going on...
	isc.lock.Lock()
	isc.subscriptions[id] = subsr
	isc.lock.Unlock()

	ch := make(chan *stomp.Message, 1)

	// Start a reader for it...
	go func(s *stomp.Subscription) {
		for {
			m, err := s.Read()
			if err == nil {
				ch <- m
			} else {
				// Remove the subscription and quit...
				isc.lock.Lock()
				delete(isc.subscriptions, id)
				isc.lock.Unlock()
				close(ch)
				return
			}
		}
	}(subsr)

	return ch, nil
}

// Close closes the stomp connection
func (isc *IoticsStompClient) Close() error {
	err := isc.stomp.Disconnect()
	if err != nil {
		return err
	}
	err = isc.ws.Close()
	return err
}

// Connect connects to a stomp server.
func (isc *IoticsStompClient) Connect(url string, authToken string) error {
	origin := "http://localhost"

	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		return err
	}
	isc.ws = ws

	stompconn, err := stomp.Connect(isc.ws,
		stomp.ConnOpt.UseStomp,
		stomp.ConnOpt.Login("", authToken),
		stomp.ConnOpt.HeartBeat(1*time.Second, 1*time.Second),
		stomp.ConnOpt.AcceptVersion(stomp.V11),
		stomp.ConnOpt.AcceptVersion(stomp.V12),
		stomp.ConnOpt.Host("dragon"))

	if err != nil {
		return err
	}
	isc.stomp = stompconn

	return nil
}

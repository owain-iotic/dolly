package main

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type HttpClient struct {
	useSSL    bool
	hostname  string
	proto     string
	hosturi   string
	authToken string
}

func NewHttpClient(hostname string, useSSL bool, authToken string) HttpClient {
	rtn := HttpClient{}
	rtn.proto = "http://"
	if useSSL {
		rtn.proto = "https://"
	}
	rtn.authToken = authToken
	rtn.hosturi = fmt.Sprintf("%s%s", rtn.proto, hostname)
	return rtn
}

func (h *HttpClient) CreateTwin(id string) error {
	client := &http.Client{}

	type Idvalue struct {
		Value string `json:"value"`
	}
	type twinInfo struct {
		TwinId Idvalue `json:"twinId"`
	}

	ti := &twinInfo{
		TwinId: Idvalue{Value: id},
	}

	data, err := json.Marshal(ti)

	body := strings.NewReader(string(data))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/qapi/twins", h.hosturi), body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Iotics-ClientAppId", "go-iotics-client")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", h.authToken))
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// eg {"twin": {"id": {"value": "did:iotics:iotJw95r16B5xjpc9PXtHt9VvLN2osXMyibN"}, "visibility": "PRIVATE"}, "alreadyCreated": false}
	bodybytes, err := ioutil.ReadAll(resp.Body)
	log.Infof("Create twin response %s %s\n", resp.Status, string(bodybytes))

	return nil
}

func (h *HttpClient) CreateFeed(twinID string, feedID string, storeLast bool) error {
	client := &http.Client{}

	type Idvalue struct {
		Value string `json:"value"`
	}
	type feedInfo struct {
		FeedId    Idvalue `json:"feedId"`
		StoreLast bool    `json:"storeLast"`
	}

	ti := &feedInfo{
		StoreLast: storeLast,
		FeedId:    Idvalue{Value: feedID},
	}

	data, err := json.Marshal(ti)

	body := strings.NewReader(string(data))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/qapi/twins/%s/feeds", h.hosturi, twinID), body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Iotics-ClientAppId", "go-iotics-client")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", h.authToken))
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// eg {"twin": {"id": {"value": "did:iotics:iotJw95r16B5xjpc9PXtHt9VvLN2osXMyibN"}, "visibility": "PRIVATE"}, "alreadyCreated": false}
	bodybytes, err := ioutil.ReadAll(resp.Body)
	log.Infof("Create twin response %s %s\n", resp.Status, string(bodybytes))

	return nil
}

func (h *HttpClient) UpdateTwin(id string, ti NewTwin) error {
	client := &http.Client{}

	data, err := json.Marshal(ti)

	body := strings.NewReader(string(data))

	fmt.Println(body)

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/qapi/twins/%s", h.hosturi, id), body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Iotics-ClientAppId", "go-iotics-client")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", h.authToken))
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	bodybytes, err := ioutil.ReadAll(resp.Body)
	log.Infof("Update twin response %s %s\n", resp.Status, string(bodybytes))

	return nil
}

func (h *HttpClient) UpdateFeed(id string, feedID string, feed NewFeed) error {
	client := &http.Client{}

	data, err := json.Marshal(feed)

	body := strings.NewReader(string(data))

	fmt.Println(body)

	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/qapi/twins/%s/feeds/%s", h.hosturi, id, feedID), body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Iotics-ClientAppId", "go-iotics-client")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", h.authToken))
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	bodybytes, err := ioutil.ReadAll(resp.Body)
	log.Infof("Update twin response %s %s\n", resp.Status, string(bodybytes))

	return nil
}

// DescribeTwin on a host
func (h HttpClient) DescribeTwin(twinID string) (*TwinResp, error) {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	res := &TwinResp{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/qapi/twins/%s", h.hosturi, twinID), nil)
	if err != nil {
		return res, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Iotics-ClientAppId", "twingen")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", h.authToken))
	resp, err := client.Do(req)
	if err != nil {
		return res, err
	}

	if resp.StatusCode != 200 {
		log.Infof("Stupid server can't cope %s\n", resp.Status)
		return res, errors.New(fmt.Sprintf("Rest server is giving us error %s", resp.Status))
	}

	bodybytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(bodybytes, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

// DescribeTwin on a host
func (h HttpClient) DescribeFeed(twinID string, feedID string) (*FeedResp, error) {
	client := &http.Client{}

	res := &FeedResp{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/qapi/twins/%s/feeds/%s", h.hosturi, twinID, feedID), nil)
	if err != nil {
		return res, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Iotics-ClientAppId", "twingen")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", h.authToken))
	resp, err := client.Do(req)
	if err != nil {
		return res, err
	}

	if resp.StatusCode != 200 {
		log.Infof("Stupid server can't cope %s\n", resp.Status)
		return res, errors.New(fmt.Sprintf("Rest server is giving us error %s", resp.Status))
	}

	bodybytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(bodybytes, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

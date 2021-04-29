package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	//"strings"
)

type HttpClient struct {
	useSSL    bool
	hostname  string
	proto     string
	hosturi   string
	authToken string
	twinID    string
}

func NewHttpClient(hostname string, useSSL bool, authToken string, twinId string) HttpClient {
	rtn := HttpClient{}
	rtn.proto = "http://"
	if useSSL {
		rtn.proto = "https://"
	}
	rtn.authToken = authToken
	rtn.hosturi = fmt.Sprintf("%s%s", rtn.proto, hostname)
	rtn.twinID = twinId
	return rtn
}

// DescribeTwin on a host
func (h HttpClient) DescribeTwin() (*TwinResp, error) {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	res := &TwinResp{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/qapi/twins/%s", h.hosturi, h.twinID), nil)
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
		fmt.Printf("Stupid server can't cope %s\n", resp.Status)
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
func (h HttpClient) DescribeFeed(feedID string) (*FeedResp, error) {
	client := &http.Client{}

	res := &FeedResp{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/qapi/twins/%s/feeds/%s", h.hosturi, h.twinID, feedID), nil)
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
		fmt.Printf("Stupid server can't cope %s\n", resp.Status)
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

package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// CreateTwin on a host
func CreateTwin(ssl bool, host string, authToken string, id string) error {
	client := &http.Client{}

	proto := "http://"
	if ssl {
		proto = "https://"
	}

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

	url := fmt.Sprintf("%s%s/qapi/twins", proto, host)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Encoding", "identity")

	uu := uuid.NewString()

	req.Header.Add("Iotics-ClientAppId", "go-iotics-client")
	req.Header.Add("Iotics-ClientRef", "go-iotics-client")
	req.Header.Add("Iotics-TransactionId", uu)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authToken))
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return errors.New(fmt.Sprintf("Twin not created: %s", resp.Status))
	}

	// eg {"twin": {"id": {"value": "did:iotics:iotJw95r16B5xjpc9PXtHt9VvLN2osXMyibN"}, "visibility": "PRIVATE"}, "alreadyCreated": false}
	//bodybytes, err := ioutil.ReadAll(resp.Body)
	//	fmt.Printf("Create twin response %s %s\n", resp.Status, string(bodybytes))

	// TODO: Return something, alreadyCreated is useful...

	return nil
}

// UpdateTwin on a host
func UpdateTwin(ssl bool, host string, authToken string, id string, data string) error {
	client := &http.Client{}

	proto := "http://"
	if ssl {
		proto = "https://"
	}

	body := strings.NewReader(data)
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s%s/qapi/twins/%s", proto, host, id), body)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Iotics-ClientAppId", "go-iotics-client")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authToken))
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// eg {"twin": {"id": {"value": "did:iotics:iotJw95r16B5xjpc9PXtHt9VvLN2osXMyibN"}, "visibility": "PRIVATE"}, "alreadyCreated": false}
	bodybytes, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("Update twin response %s %s\n", resp.Status, string(bodybytes))

	return nil
}

// CreateFeed on a twin
func CreateFeed(ssl bool, host string, authToken string, twinid string, feedid string) error {
	client := &http.Client{}

	proto := "http://"
	if ssl {
		proto = "https://"
	}

	type Idvalue struct {
		Value string `json:"value"`
	}
	type feedInfo struct {
		FeedId Idvalue `json:"feedId"`
	}

	fi := &feedInfo{
		FeedId: Idvalue{Value: feedid},
	}

	data, err := json.Marshal(fi)

	body := strings.NewReader(string(data))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s/qapi/twins/%s/feeds", proto, host, twinid), body)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Iotics-ClientAppId", "go-iotics-client")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authToken))
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// eg {"feed": {"id": {"value": "random_temperature_feed"}, "twinId": {"value": "did:iotics:iotJw95r16B5xjpc9PXtHt9VvLN2osXMyibN"}}, "alreadyCreated": false}
	bodybytes, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("Create feed response %s %s\n", resp.Status, string(bodybytes))

	return nil
}

// ListAllTwins on a host
func ListAllTwins(ssl bool, host string, authToken string) ([]Twindata, error) {
	client := &http.Client{}

	proto := "http://"
	if ssl {
		proto = "https://"
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s/qapi/twins", proto, host), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Iotics-ClientAppId", "go-iotics-client")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authToken))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Stupid server can't cope %s\n", resp.Status)
		return nil, errors.New(fmt.Sprintf("Rest server is giving us error %s", resp.Status))
	}

	bodybytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res := &TwinsResp{}

	err = json.Unmarshal(bodybytes, &res)
	if err != nil {
		return nil, err
	}

	return res.Twins, nil
}

// DescribeTwin on a host
func DescribeTwin(ssl bool, host string, authToken string, id string) (*TwinResp, error) {
	client := &http.Client{}

	proto := "http://"
	if ssl {
		proto = "https://"
	}

	res := &TwinResp{}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s/qapi/twins/%s", proto, host, id), nil)
	if err != nil {
		return res, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Iotics-ClientAppId", "op")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authToken))
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

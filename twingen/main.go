package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
)

type Twinner struct {
	TwinID     string
	authToken  string
	httpClient *HttpClient
	Twin       *TwinResp
	Feeds      []*FeedResp
}

func (t *Twinner) LoadFeeds() error {
	for _, v := range t.Twin.Result.Feeds {
		feed, err := t.httpClient.DescribeFeed(v.FeedId.Value)
		if err != nil {
			return err
		}
		t.Feeds = append(t.Feeds, feed)
	}
	return nil
}

func NewTwinner(twinID string, authToken string, httpClient *HttpClient) Twinner {
	return Twinner{
		TwinID:     twinID,
		authToken:  authToken,
		httpClient: httpClient,
	}
}

// Load loads the source twin up so we can use it to generate a code template
func (t *Twinner) Load() error {
	var err error
	t.Twin, err = t.httpClient.DescribeTwin()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Loaded twin %v\n", t.Twin)

	err = t.LoadFeeds()
	if err != nil {
		return err
	}

	for i, feed := range t.Feeds {
		fmt.Printf(" - Feed[%d] %v\n", i, feed)
	}

	return nil
}

// loadFile loads file
func (t *Twinner) loadFile(name string) (string, error) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Generate generates the code from template and using the twin data...
func (t *Twinner) Generate() error {
	//td := Todo{"Test templates", "Let's test a template to see the magic."}
	fmt.Println(t.Twin.Result.Labels[0].Value)
	fileContent, err := t.loadFile("./templates/test.go.tmpl")
	if err != nil {
		panic(err)
	}
	tem, err := template.New("main").Parse(fileContent)
	if err != nil {
		return err
	}

	err = tem.Execute(os.Stdout, t)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	fmt.Println("Twingen starting...")
	twinID := "did:iotics:iotCya2zUUN4TGbtp8FHwLPjewsNbJpffHpr"
	authToken := "eyJ0eXAiOiJKV1QiLCJhbGciOiJFUzI1NiJ9.eyJpc3MiOiJkaWQ6aW90aWNzOmlvdEhCQ21wUHZUUVJySndXZFhNNTZhMTltclhLd0g0NmZGTCNhZ2VudC0wIiwiYXVkIjoiaHR0cHM6Ly9kaWQucHJkLmlvdGljcy5jb20iLCJzdWIiOiJkaWQ6aW90aWNzOmlvdENkdWpWQ3ZCNllQQ1JGa1VNTnpjSnVNMVdkUUZhcHBpVyIsImlhdCI6MTYxOTY4MzA3MiwiZXhwIjoxNjE5NzExOTAyfQ.zdiXHK39scpHJjwL3EOeSKGMtjroculC6XemPmjWLZ5KBtS_X2kLfwEXiRF_43zm9DHeB0K-oz4PrPxIrzc4iw"
	httpClient := NewHttpClient("plateng.iotics.space", true, authToken, twinID)

	twinner := NewTwinner(twinID, authToken, &httpClient)
	fmt.Printf("Loading twin %s...\n", twinID)
	twinner.Load()

	fmt.Println("Generating code...")
	err := twinner.Generate()
	fmt.Println(err)
}

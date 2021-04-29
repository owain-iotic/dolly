package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
)

const (
	twin_template = "./templates/test.go.tmpl"
)

type Twinner struct {
	TwinID     string
	authToken  string
	httpClient *HttpClient
	Twin       *TwinResp
	Feeds      []*FeedResp
}

// Create a new Twinner to copy an existing twin into a code template
func NewTwinner(twinID string, authToken string, httpClient *HttpClient) Twinner {
	return Twinner{
		TwinID:     twinID,
		authToken:  authToken,
		httpClient: httpClient,
	}
}

// LoadFeeds runs describe feed on each one so we get details
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

// Load loads the source twin up so we can use it to generate a code template
func (t *Twinner) Load() error {
	var err error
	t.Twin, err = t.httpClient.DescribeTwin()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Loaded twin %v\n", t.Twin)
	js, err := json.Marshal(t.Twin)
	fmt.Printf("Twin JSON %s\n", js)

	err = t.LoadFeeds()
	if err != nil {
		return err
	}

	for i, feed := range t.Feeds {
		fmt.Printf(" - Feed[%d] %v\n", i, feed)
		js, err = json.Marshal(feed)
		fmt.Printf("Feed JSON %s\n", js)

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

	// Load up the template
	fileContent, err := t.loadFile(twin_template)
	if err != nil {
		panic(err)
	}
	// Create a template from the file
	tem, err := template.New("main").Parse(fileContent)
	if err != nil {
		return err
	}

	// Execute the template using Twinner as the data source...
	err = tem.Execute(os.Stdout, t)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	fmt.Println("Twingen starting...")
	twinID := "did:iotics:iotCya2zUUN4TGbtp8FHwLPjewsNbJpffHpr"
	authToken := "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJub3QtY3VycmVudGx5LXVzZWQuZXhhbXBsZS5jb20iLCJleHAiOjE2MTk2OTc2NzMsImlhdCI6MTYxOTY4NDA3MywiaXNzIjoiZGlkOmlvdGljczppb3RXRGpoMkZjUmZIeHdDajdXQjhtbjJHQ29LYWJWZXc5OTkjYWdlbnQtMCIsInN1YiI6ImRpZDppb3RpY3M6aW90Uml6NmFUeUpCaVJGUkJObWprckthUHBaeHltN0IzUnV0In0.GPu_iFx_GJ1vxkxXxNDey9YPGvnXufLvKJauay-oe6v5wJZ1iVyll-1xVOY99InRr3OSU8AGXj8CFcl1LoMt1w"
	httpClient := NewHttpClient("plateng.iotics.space", true, authToken, twinID)

	twinner := NewTwinner(twinID, authToken, &httpClient)
	fmt.Printf("Loading twin %s...\n", twinID)
	twinner.Load()

	fmt.Println("Generating code...")
	err := twinner.Generate()
	fmt.Println(err)
}

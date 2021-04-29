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

func (t *Twinner) Load() error {
	var err error
	t.Twin, err = t.httpClient.DescribeTwin()
	if err != nil {
		fmt.Println(err)
	}

	err = t.LoadFeeds()
	if err != nil {
		return err
	}
	fmt.Println(t.Feeds[0].Feed.TwinID)
	return nil
}

func (t *Twinner) loadFile(name string) (string, error) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

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
	fmt.Println("hi")
	twinID := "did:iotics:iotCya2zUUN4TGbtp8FHwLPjewsNbJpffHpr"
	authToken := "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJub3QtY3VycmVudGx5LXVzZWQuZXhhbXBsZS5jb20iLCJleHAiOjE2MTk1MjM3MDgsImlhdCI6MTYxOTUxMDEwOCwiaXNzIjoiZGlkOmlvdGljczppb3RXRGpoMkZjUmZIeHdDajdXQjhtbjJHQ29LYWJWZXc5OTkjYWdlbnQtMCIsInN1YiI6ImRpZDppb3RpY3M6aW90Uml6NmFUeUpCaVJGUkJObWprckthUHBaeHltN0IzUnV0In0.1Qv_xujI8joRFX13wciFIzQUvK70iALGZN7Rb6sff7WldSU8YyS2CFHlu69XB60ulPDCHy67QV-xtuX1Zkku2w"
	httpClient := NewHttpClient("plateng.iotics.space", true, authToken, twinID)

	twinner := NewTwinner(twinID, authToken, &httpClient)
	twinner.Load()
	err := twinner.Generate()
	fmt.Println(err)
}

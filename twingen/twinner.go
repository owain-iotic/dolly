package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

type Twinner struct {
	TwinID       string
	authToken    string
	httpClient   *HttpClient
	Twin         *TwinResp
	Feeds        []*FeedResp
	outputFolder string
}

// Create a new Twinner to copy an existing twin into a code template
func NewTwinner(twinID string, authToken string, httpClient *HttpClient, outputFolder string) Twinner {
	return Twinner{
		TwinID:       twinID,
		authToken:    authToken,
		httpClient:   httpClient,
		outputFolder: outputFolder,
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
		return err
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

type TemplateModel struct {
	TwinName      string
	LCaseTwinName string
	Visibility    string
	Feeds         []FeedTmpl
	Comments      []CommentsTmpl
	Properties    []PropertiesTmpl
	Labels        []LabelsTmpl
	Tags          []string
}
type PropertiesTmpl struct {
	Key   string
	Value string
}
type CommentsTmpl struct {
	Lang  string
	Value string
}

type LabelsTmpl struct {
	Lang  string
	Value string
}
type FeedTmpl struct {
	FeedName            string
	FeedID              string
	FeedStructName      string
	FeedStructNameLCase string
	FeedValues          []FeedValuesTmpl
	FeedComments        []FeedCommentsTmpl
	FeedLabels          []FeedLabelsTmpl
	FeedTags            []string
}

type FeedCommentsTmpl struct {
	Lang  string
	Value string
}

type FeedLabelsTmpl struct {
	Lang  string
	Value string
}
type FeedValuesTmpl struct {
	Label      string
	Comment    string
	Unit       string
	DataType   string
	GoDataType string
	LCaseLabel string
}

func (t *Twinner) GetFeedByName(name string) *FeedResp {
	for _, v := range t.Feeds {
		if v.Feed.ID.Value == name {
			return v
		}
	}
	return nil
}

// Generate generates the code from template and using the twin data...
func (t *Twinner) Generate() error {

	if _, err := os.Stat(t.outputFolder); os.IsNotExist(err) {
		os.Mkdir(t.outputFolder, 0755)
	}

	c := NewTemplateModel()
	c.TwinName = strings.Title(t.Twin.Result.Tags[0])
	c.LCaseTwinName = strings.ToLower(t.Twin.Result.Tags[0])
	c.Visibility = t.Twin.Twin.Visibility
	c.Feeds = []FeedTmpl{}
	c.Comments = []CommentsTmpl{}
	c.Properties = []PropertiesTmpl{}
	c.Labels = []LabelsTmpl{}
	c.Tags = []string{}
	for _, v := range t.Twin.Result.Tags {
		c.Tags = append(c.Tags, v)
	}
	for _, v := range t.Twin.Result.Labels {
		c.Labels = append(c.Labels, LabelsTmpl{
			Value: v.Value,
			Lang:  v.Lang,
		})
	}
	for _, v := range t.Twin.Result.Comments {
		c.Comments = append(c.Comments, CommentsTmpl{
			Value: v.Value,
			Lang:  v.Lang,
		})
	}
	for _, v := range t.Twin.Result.Properties {
		c.Properties = append(c.Properties, PropertiesTmpl{
			Value: v.StringLiteralValue.Value,
			Key:   v.Key,
		})
	}
	for _, v := range t.Twin.Result.Feeds {

		feed := t.GetFeedByName(v.FeedId.Value)
		fv := []FeedValuesTmpl{}
		fc := []FeedCommentsTmpl{}
		fl := []FeedLabelsTmpl{}
		ft := []string{}
		for _, v := range feed.Result.Values {
			i := FeedValuesTmpl{
				DataType: v.DataType,
				Comment:  v.Comment,
				Unit:     v.Unit,
				Label:    v.Label,
			}
			// workout types here
			i.GoDataType = v.DataType
			switch v.DataType {
			case "decimal":
				i.GoDataType = "float32"
			case "float":
				i.GoDataType = "float32"
			case "int":
				i.GoDataType = "int"
			case "integer":
				i.GoDataType = "int"
			case "string":
				i.GoDataType = "string"
			case "boolean":
				i.GoDataType = "bool"
			default:
				i.GoDataType = fmt.Sprintf("not found datatype %s ", v.DataType)
			}
			i.LCaseLabel = strings.ToLower(v.Label)
			fv = append(fv, i)
		}
		for _, v := range feed.Result.Comments {
			fc = append(fc, FeedCommentsTmpl{
				Lang:  v.Lang,
				Value: v.Value,
			})
		}
		for _, v := range feed.Result.Labels {
			fl = append(fl, FeedLabelsTmpl{
				Lang:  v.Lang,
				Value: v.Value,
			})
		}
		for _, v := range feed.Result.Tags {
			ft = append(ft, v)
		}
		c.Feeds = append(c.Feeds, FeedTmpl{
			FeedName:            strings.Title(v.FeedId.Value),
			FeedID:              v.FeedId.Value,
			FeedStructName:      fmt.Sprintf("Feed%s", strings.Title(v.FeedId.Value)),
			FeedStructNameLCase: fmt.Sprintf("Feed%s", strings.ToLower(v.FeedId.Value)),
			FeedValues:          fv,
			FeedTags:            ft,
			FeedComments:        fc,
			FeedLabels:          fl,
		})
	}

	err := t.render("./templates/library.go.tmpl", "library", fmt.Sprintf("%s/library.go", t.outputFolder), c)
	if err != nil {
		return err
	}
	err = t.render("./templates/main.go.tmpl", "main", fmt.Sprintf("%s/main.go", t.outputFolder), c)
	if err != nil {
		return err
	}
	err = t.render("./templates/go.mod.tmpl", "mod", fmt.Sprintf("%s/go.mod", t.outputFolder), c)
	if err != nil {
		return err
	}
	return nil
}

func (t *Twinner) render(temp string, name string, path string, c interface{}) error {
	fileContent, err := t.loadFile(temp)
	if err != nil {
		return (err)
	}
	// Create a template from the file
	tem, err := template.New(name).Parse(fileContent)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		fmt.Println("create file: ", err)
		return err
	}

	// Execute the template using Twinner as the data source...
	err = tem.Execute(f, c)
	if err != nil {
		return err
	}
	return nil
}

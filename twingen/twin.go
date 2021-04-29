package main

import (
// "encoding/json"
// "errors"
// "fmt"
)

// For /twins
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
	Lat float64
	Lon float64
}

type Label struct {
	Lang  string
	Value string
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
	Label    string `json:"label"`
	Comment  string `json:"comment"`
	Unit     string `json:"unit"`
	DataType string `json:"dataType"`
}

type FeedResult struct {
	Labels    []Label
	Comments  []Label
	Values    []Values
	Tags      []string
	StoreLast bool
}



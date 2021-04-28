package client

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

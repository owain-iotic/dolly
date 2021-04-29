package main

import (
	"flag"
	"fmt"
)

const (

)

func main() {

	twinidPtr := flag.String("twinid", "did:iotics:iotTtKdce1CKpUYQPb9AQ2TggDgesKmHD6aF", "did of twin to clone")
	authPtr := flag.String("jwt", "ey....", "jwt")
	hostPtr := flag.String("host", "plateng.iotics.space", "iotic host")
	outFolderPtr := flag.String("output", "twin", "output folder")

	flag.Parse()

	fmt.Println("Twingen starting...")
	twinID := *twinidPtr
	authToken := *authPtr
	httpClient := NewHttpClient(*hostPtr, true, authToken, twinID)

	twinner := NewTwinner(twinID, authToken, &httpClient, *outFolderPtr)
	fmt.Printf("Loading twin %s...\n", twinID)
	err := twinner.Load()
	if err != nil {
		panic(err)
	}

	fmt.Println("Generating code...")
	err = twinner.Generate()
	fmt.Println(err)
}

func NewTemplateModel() TemplateModel {
	return TemplateModel{}
}

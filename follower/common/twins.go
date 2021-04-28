package common

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var Shiptwins = make(map[string]string)

func LoadShipConfig(shipfile string, twinfile string) {

	shipids := make([]string, 0)
	shfile, err := os.Open(shipfile)
	if err != nil {
		log.Fatal(err)
	}
	defer shfile.Close()

	shscanner := bufio.NewScanner(shfile)
	for shscanner.Scan() {
		id := shscanner.Text()
		shipids = append(shipids, id)
	}

	fmt.Printf("Loaded %d ship ids...\n", len(shipids))

	twinids := make([]string, 0)
	twfile, err := os.Open(twinfile)
	if err != nil {
		log.Fatal(err)
	}
	defer twfile.Close()

	twscanner := bufio.NewScanner(twfile)
	for twscanner.Scan() {
		id := twscanner.Text()
		twinids = append(twinids, id)
	}

	fmt.Printf("Loaded %d twin ids...\n", len(twinids))

	for i := 0; i < len(shipids); i++ {
		// Create the mapping...
		id := shipids[i]
		did := twinids[i]
		Shiptwins[id] = did
		fmt.Printf("Mapped %s to %s\n", id, did)
	}
}

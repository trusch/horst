package links

import (
	"fmt"
	"log"
	"testing"
)

func TestLinkMap(t *testing.T) {
	links := make(LinkMap)
	links.Add(&Link{"A", "a", "B", "a"})
	links.Add(&Link{"B", "a", "C", "a"})
	links.Add(&Link{"C", "a", "D", "a"})
	jsonData, err := links.ToJSON()
	if err != nil {
		log.Print(err)
		t.Fail()
	}
	err = links.FromJSON(jsonData)
	if err != nil {
		log.Print(err)
		t.Fail()
	}
	fmt.Println(links)
}

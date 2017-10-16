package main

import (
	"flag"
	"log"

	"github.com/trusch/horst/components/twittersource"
	"github.com/trusch/horst/runner"
)

var listenAddr = flag.String("listen", ":80", "listen address")
var etcdAddr = flag.String("etcd", "etcd:2379", "etcd address")
var id = flag.String("id", "twittersource", "id of this instance")

func main() {
	flag.Parse()
	c, err := twittersource.New()
	if err != nil {
		log.Fatal(err)
	}
	r, err := runner.New(*id, c, *etcdAddr, *listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	if err = r.Start(); err != nil {
		log.Fatal(err)
	}
	log.Print("started twitter data source")
	defer r.Stop()
	select {}
}

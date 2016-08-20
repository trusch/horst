package main

import (
	"github.com/trusch/horst/config"
	"github.com/trusch/horst/runner"
	"log"

	_ "github.com/trusch/horst/processors/logger"

	_ "github.com/trusch/horst/processors/projector"

	_ "github.com/trusch/horst/processors/twittersource"
)

func main() {
	cfg, err := config.New("config.json")
	if err != nil {
		log.Fatal(err)
	}
	runner, err := runner.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	runner.Run()
	select {}
}

package main

import (
	"log"

	"github.com/trusch/horst"
	"github.com/trusch/horst/config"
	"github.com/trusch/horst/links"
	"github.com/trusch/horst/manager"
	_ "github.com/trusch/horst/processors/blackhole"
	_ "github.com/trusch/horst/processors/duplicator"
	_ "github.com/trusch/horst/processors/logger"
	_ "github.com/trusch/horst/processors/naturals"
	_ "github.com/trusch/horst/processors/projector"
	_ "github.com/trusch/horst/processors/roundrobin"
	_ "github.com/trusch/horst/processors/twittersource"
	"github.com/trusch/horst/registry"
)

type horstCli struct {
	config  config.Config
	links   links.LinkMap
	manager horst.ProcessorManager
}

func (cli *horstCli) Start() {
	for proc, cfg := range cli.config {
		cli.LoadProcessor(cfg.Class, proc, cfg.Config)
	}
}

func (cli *horstCli) LoadProcessor(className, id string, config interface{}) {
	proc, err := registry.Construct(className, id, config, cli.manager)
	if err != nil {
		log.Fatal(err)
	}
	cli.manager.AddProcessor(id, proc)
}

func (cli *horstCli) UnloadProcessor(id string) {
	cli.manager.DelProcessor(id)
}

func (cli *horstCli) UpdateLink(from, fromOut, to, toIn string) {
	cli.links.Add(from, fromOut, to, toIn)
}

func main() {
	cfg := config.Config{}
	err := cfg.Load("config.json")
	if err != nil {
		log.Fatal(err)
	}
	linkMap, err := cfg.GetLinkMap()
	if err != nil {
		log.Fatal(err)
	}
	cli := horstCli{cfg, linkMap, manager.New(linkMap)}
	cli.Start()
	select {}
}

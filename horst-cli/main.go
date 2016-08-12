package main

import (
	"log"
	"time"

	"github.com/trusch/horst"
	"github.com/trusch/horst/links"
	"github.com/trusch/horst/manager"
	_ "github.com/trusch/horst/processors/logger"
	_ "github.com/trusch/horst/processors/naturals"
	"github.com/trusch/horst/registry"
)

type horstCli struct {
	links   links.LinkMap
	manager horst.ProcessorManager
}

func (cli *horstCli) LoadProcessor(className, id string) {
	proc, err := registry.Construct(className, id, nil, cli.manager)
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
	linkMap := links.LinkMap{}
	linkMap.Add("naturals1", "out", "logger1", "default")
	cli := horstCli{linkMap, manager.New(linkMap)}

	cli.LoadProcessor("logger", "logger1")
	cli.LoadProcessor("naturals", "naturals1")

	time.Sleep(5 * time.Second)

	cli.UnloadProcessor("logger1")
	cli.LoadProcessor("logger", "logger2")
	cli.UpdateLink("naturals1", "out", "logger2", "default")

	select {}
}

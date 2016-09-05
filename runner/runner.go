package runner

import (
	"github.com/trusch/horst"
	"github.com/trusch/horst/config"
	"github.com/trusch/horst/links"
	"github.com/trusch/horst/manager"
	"github.com/trusch/horst/registry"
)

// Runner runs a set of processors given by a config
type Runner struct {
	config  config.Config
	links   links.LinkMap
	manager horst.ProcessorManager
}

// New creates a new runner from a given config file
func New(config config.Config) (*Runner, error) {
	linkMap, err := config.GetLinkMap()
	if err != nil {
		return nil, err
	}
	return &Runner{
		config:  config,
		links:   linkMap,
		manager: manager.New(linkMap),
	}, nil
}

// Run instanciates all processors
func (cli *Runner) Run() {
	for proc, cfg := range cli.config {
		cli.LoadProcessor(cfg.Class, proc)
	}
}

// LoadProcessor loads a procssor at runtime
func (cli *Runner) LoadProcessor(className, id string) error {
	proc, err := registry.Construct(className, id, cli.config[id].Config, cli.manager)
	if err != nil {
		return err
	}
	cli.manager.AddProcessor(id, proc)
	return nil
}

// UnloadProcessor unloads a procssor at runtime
func (cli *Runner) UnloadProcessor(id string) {
	cli.manager.DelProcessor(id)
}

// UpdateLink overwrites a link
func (cli *Runner) UpdateLink(from, fromOut, to, toIn string) {
	cli.links.Add(from, fromOut, to, toIn)
}

// Process inserts a doc into the pipeline
func (cli *Runner) Process(to, toIn string, doc interface{}) {
	cli.manager.Process(to, toIn, doc)
}

// UpdateConfig overwrites the config part of a single processor
func (cli *Runner) UpdateConfig(id string, config interface{}) error {
	cli.config[id].Config = config
	cli.UnloadProcessor(id)
	return cli.LoadProcessor(cli.config[id].Class, id)
}

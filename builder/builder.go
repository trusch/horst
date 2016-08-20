package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"text/template"

	"github.com/trusch/horst/config"
)

var configFile = flag.String("c", "config.json", "horst config file")
var outputFile = flag.String("o", "generatedRunner", "output binary and code name")

func main() {
	flag.Parse()
	cfg, err := config.New(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	classes := make(map[string]bool)
	for _, nodeCfg := range cfg {
		classes[nodeCfg.Class] = true
	}
	f, err := os.Create(*outputFile + ".go")
	if err != nil {
		log.Fatal(err)
	}
	err = runnerTemplate.Execute(f, classes)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
	cmd := exec.Command("go", "fmt", *outputFile+".go")
	cmd.Run()
	cmd = exec.Command("go", "build", "-o", *outputFile, *outputFile+".go")
	cmd.Run()
}

var runnerTemplateString = `package main

import (
  "log"
	"github.com/trusch/horst/config"
	"github.com/trusch/horst/runner"
  {{ range $key, $val := . }}
  _ "{{ $key }}"
  {{ end }}
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

`
var runnerTemplate = template.Must(template.New("").Parse(runnerTemplateString))

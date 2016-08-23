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

func run(args ...string) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

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
	for class := range classes {
		run("go", "get", "-v", class)
	}
	err = runnerTemplate.Execute(f, classes)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
	run("go", "fmt", *outputFile+".go")
	run("go", "build", "-o", *outputFile, *outputFile+".go")
}

var runnerTemplateString = `package main

import (
  "log"
  "flag"
	"github.com/trusch/horst/config"
	"github.com/trusch/horst/runner"
	"github.com/trusch/horst/server"
	{{ range $key, $val := . }}
  _ "{{ $key }}"
  {{ end }}
)

var configFile = flag.String("c", "config.json", "horst config file")
var addr = flag.String("a", ":5566", "the control servers address")

func main() {
  flag.Parse()
	cfg, err := config.New(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	runner, err := runner.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	runner.Run()
	_, err = server.New(runner, *addr)
	if err != nil {
		log.Fatal(err)
	}
	select {}
}

`
var runnerTemplate = template.Must(template.New("").Parse(runnerTemplateString))

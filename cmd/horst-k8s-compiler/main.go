package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/trusch/horst/config"
)

type Deployment struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Name   string            `json:"name"`
		Labels map[string]string `json:"labels"`
	} `json:"metadata"`
	Spec struct {
		Replicas int `json:"replicas"`
		Selector struct {
			MatchLabels map[string]string `json:"matchLabels"`
		} `json:"selector"`
		Template struct {
			Metadata struct {
				Labels map[string]string `json:"labels"`
			} `json:"metadata"`
			Spec struct {
				Containers []*ContainerSpec `json:"containers"`
			} `json:"spec"`
		} `json:"template"`
	} `json:"spec"`
}

type ContainerSpec struct {
	Name  string      `json:"name"`
	Image string      `json:"image"`
	Ports []*PortSpec `json:"ports"`
}

type PortSpec struct {
	ContainerPort int `json:"containerPort,omitempty"`
	Port          int `json:"port,omitempty"`
}

type Service struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Name string `json:"name"`
	} `json:"metadata"`
	Spec struct {
		Selector map[string]string `json:"selector"`
		Ports    []*PortSpec       `json:"ports"`
	}
}

var cfgFile = flag.String("config", "config.json", "config file")

func main() {
	flag.Parse()

	cfg := make(map[string]*config.ComponentConfig)
	bs, err := ioutil.ReadFile(*cfgFile)
	if err != nil {
		log.Fatal(err)
	}
	if err = json.Unmarshal(bs, &cfg); err != nil {
		log.Fatal(err)
	}

	for id, conf := range cfg {
		dep := &Deployment{}
		dep.APIVersion = "extensions/v1beta1"
		dep.Kind = "Deployment"
		dep.Metadata.Name = id
		dep.Metadata.Labels = make(map[string]string)
		dep.Metadata.Labels["app"] = id
		dep.Spec.Replicas = 1
		dep.Spec.Selector.MatchLabels = make(map[string]string)
		dep.Spec.Selector.MatchLabels["app"] = id
		dep.Spec.Template.Metadata.Labels = make(map[string]string)
		dep.Spec.Template.Metadata.Labels["app"] = id
		dep.Spec.Template.Spec.Containers = []*ContainerSpec{&ContainerSpec{
			Name:  id,
			Image: conf.Image,
			Ports: []*PortSpec{&PortSpec{ContainerPort: 80}},
		}}
		df, err := os.Create(id + "_deployment.json")
		if err != nil {
			log.Fatal(err)
		}
		defer df.Close()
		encoder := json.NewEncoder(df)
		if err = encoder.Encode(dep); err != nil {
			log.Fatal(err)
		}
		svc := &Service{}
		svc.APIVersion = "v1"
		svc.Kind = "Service"
		svc.Metadata.Name = id
		svc.Spec.Selector = make(map[string]string)
		svc.Spec.Selector["app"] = id
		svc.Spec.Ports = []*PortSpec{&PortSpec{Port: 80}}
		sf, err := os.Create(id + "_service.json")
		if err != nil {
			log.Fatal(err)
		}
		defer sf.Close()
		encoder = json.NewEncoder(sf)
		if err = encoder.Encode(svc); err != nil {
			log.Fatal(err)
		}
	}
}

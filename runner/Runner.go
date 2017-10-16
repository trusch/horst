package runner

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/trusch/horst/components"
)

// Runner wraps a Component in an HTTP server,
type Runner struct {
	ID        string
	component components.Component
	etcd      *clientv3.Client
	stop      chan struct{}
	server    *http.Server
}

// New creates a new runner instance
func New(id string, component components.Component, etcdURI, listenAddress string) (*Runner, error) {
	etcd, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{etcdURI},
		DialTimeout: 3 * time.Second,
	})

	if err != nil {
		return nil, err
	}
	r := &Runner{id, component, etcd, nil, &http.Server{Addr: listenAddress}}
	r.server.Handler = r
	return r, nil
}

// Start reads config and outputs, starts watching for changes and starts the http serve
func (r *Runner) Start() error {
	if err := r.loadConfig(); err != nil {
		return err
	}
	if err := r.loadOutputs(); err != nil {
		return err
	}
	go r.watch()
	go r.server.ListenAndServe()
	return nil
}

// Stop stops the runner
func (r *Runner) Stop() error {
	close(r.stop)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := r.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func (r *Runner) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var event interface{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&event); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if err := r.component.Process(req.URL.Path, event); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
}

func (r *Runner) loadConfig() error {
	resp, err := r.etcd.Get(context.Background(), "/horst/configs/"+r.ID)
	if err != nil {
		return err
	}
	if len(resp.Kvs) != 1 {
		return errors.New("no config object available")
	}
	cfg := make(map[string]interface{})
	if err := json.Unmarshal(resp.Kvs[0].Value, &cfg); err != nil {
		return err
	}
	return r.component.HandleConfigUpdate(cfg)
}

func (r *Runner) loadOutputs() error {
	resp, err := r.etcd.Get(context.Background(), "/horst/outputs/"+r.ID)
	if err != nil {
		return err
	}
	if len(resp.Kvs) != 1 {
		return errors.New("no output object available")
	}
	outputs := make(map[string]string)
	if err := json.Unmarshal(resp.Kvs[0].Value, &outputs); err != nil {
		return err
	}
	return r.component.HandleOutputUpdate(outputs)
}

func (r *Runner) watch() {
	r.stop = make(chan struct{})
	cfgChan := r.etcd.Watch(context.Background(), "/horst/configs/"+r.ID)
	outputsChan := r.etcd.Watch(context.Background(), "/horst/outputs/"+r.ID)
	for {
		select {
		case resp := <-cfgChan:
			{
				if len(resp.Events) == 1 && resp.Events[0].Type == mvccpb.PUT {
					cfg := make(map[string]interface{})
					if err := json.Unmarshal(resp.Events[0].Kv.Value, &cfg); err != nil {
						log.Print(err)
						continue
					}
					if err := r.component.HandleConfigUpdate(cfg); err != nil {
						log.Print(err)
						continue
					}
				}
			}
		case resp := <-outputsChan:
			{
				if len(resp.Events) == 1 && resp.Events[0].Type == mvccpb.PUT {
					outputs := make(map[string]string)
					if err := json.Unmarshal(resp.Events[0].Kv.Value, &outputs); err != nil {
						log.Print(err)
						continue
					}
					if err := r.component.HandleOutputUpdate(outputs); err != nil {
						log.Print(err)
						continue
					}
				}
			}
		case <-r.stop:
			{
				break
			}
		}
	}
}

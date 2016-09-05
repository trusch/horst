package server

import (
	"encoding/json"
	"log"
	"net"
	"time"

	"github.com/trusch/horst/runner"
)

// Server represents the server managing the runner
type Server struct {
	runner   *runner.Runner
	listener net.Listener
}

// A Message is the format of the servers protocol
type Message struct {
	Command string      `json:"cmd"`
	Payload interface{} `json:"payload"`
}

var errorString = "error"
var successString = "success"

func (server *Server) backend() {
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			log.Print(err)
			time.Sleep(1 * time.Second)
		}
		go server.handleConnection(conn)
	}
}

func (server *Server) processLoad(msg Message, encoder *json.Encoder) {
	if payload, ok := msg.Payload.(map[string]interface{}); ok {
		if id, ok := payload["id"].(string); ok {
			if class, ok := payload["class"].(string); ok {
				err := server.runner.LoadProcessor(class, id)
				if err != nil {
					msg.Command = errorString
					msg.Payload = err.Error()
					encoder.Encode(msg)
					return
				}
				msg.Command = successString
				encoder.Encode(msg)
				return
			}
		}
	}
	msg.Command = errorString
	msg.Payload = "need object with id and class as payload"
	encoder.Encode(msg)
}

func (server *Server) processUnload(msg Message, encoder *json.Encoder) {
	if payload, ok := msg.Payload.(map[string]interface{}); ok {
		if id, ok := payload["id"].(string); ok {
			server.runner.UnloadProcessor(id)
			msg.Command = successString
			encoder.Encode(msg)
			return
		}
	}
	msg.Command = errorString
	msg.Payload = "need object with id as payload"
	encoder.Encode(msg)
}

func (server *Server) processUpdateLink(msg Message, encoder *json.Encoder) {
	if payload, ok := msg.Payload.(map[string]interface{}); ok {
		if from, ok := payload["from"].(string); ok {
			if fromOut, ok := payload["fromOut"].(string); ok {
				if to, ok := payload["to"].(string); ok {
					if toIn, ok := payload["toIn"].(string); ok {
						server.runner.UpdateLink(from, fromOut, to, toIn)
						msg.Command = successString
						encoder.Encode(msg)
						return
					}
				}
			}
		}
	}
	msg.Command = errorString
	msg.Payload = "need object with from, fromOut, to and toIn as payload"
	encoder.Encode(msg)
}

func (server *Server) processDoc(msg Message, encoder *json.Encoder) {
	if payload, ok := msg.Payload.(map[string]interface{}); ok {
		if to, ok := payload["to"].(string); ok {
			if toIn, ok := payload["toIn"].(string); ok {
				if value, ok := payload["value"]; ok {
					server.runner.Process(to, toIn, value)
					msg.Command = successString
					encoder.Encode(msg)
					return
				}
			}
		}
	}
	msg.Command = errorString
	msg.Payload = "need object with to, toIn and value as payload"
	encoder.Encode(msg)
}

func (server *Server) processUpdateConfig(msg Message, encoder *json.Encoder) {
	if payload, ok := msg.Payload.(map[string]interface{}); ok {
		if id, ok := payload["id"].(string); ok {
			config := payload["config"]
			err := server.runner.UpdateConfig(id, config)
			if err != nil {
				msg.Command = errorString
				msg.Payload = err.Error()
				encoder.Encode(msg)
				return
			}
			msg.Command = successString
			encoder.Encode(msg)
			return
		}
	}
	msg.Command = errorString
	msg.Payload = "need object with id and config as payload"
	encoder.Encode(msg)
}

func (server *Server) handleConnection(conn net.Conn) {
	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)
	var msg Message
	for {
		if err := decoder.Decode(&msg); err == nil {
			switch msg.Command {
			case "load":
				{
					server.processLoad(msg, encoder)
				}
			case "unload":
				{
					server.processUnload(msg, encoder)
				}
			case "updateLink":
				{
					server.processUpdateLink(msg, encoder)
				}
			case "updateConfig":
				{
					server.processUpdateConfig(msg, encoder)
				}
			case "process":
				{
					server.processDoc(msg, encoder)
				}
			default:
				{
					msg.Command = errorString
					msg.Payload = "unknown cmd, need one of load, unload, updateLink, updateConfig, process"
					encoder.Encode(msg)
				}
			}
		} else {
			msg.Command = errorString
			msg.Payload = "malformed json: " + err.Error()
			encoder.Encode(msg)
		}
	}
}

// New creates a new server for a given runner and listening address
func New(runner *runner.Runner, addr string) (*Server, error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	server := &Server{runner, ln}
	go server.backend()
	return server, nil
}

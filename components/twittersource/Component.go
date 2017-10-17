package twittersource

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/mitchellh/mapstructure"
	"github.com/trusch/horst/components"
	"github.com/trusch/horst/components/base"
)

// Component is the most basic component we can build
type Component struct {
	base.Component
	config *twitterConfig
	client *twitter.Client
	stream *twitter.Stream
}

type twitterConfig struct {
	ConsumerKey    string   `json:"consumerKey"`
	ConsumerSecret string   `json:"consumerSecret"`
	AccessToken    string   `json:"accessToken"`
	AccessSecret   string   `json:"accessSecret"`
	Track          []string `json:"track"`
}

// New returns a new twittersource.Component
func New() (components.Component, error) {
	return &Component{}, nil
}

// HandleConfigUpdate gets called when new config for this component is available
func (c *Component) HandleConfigUpdate(config map[string]interface{}) error {
	log.Print("handle update")
	if c.stream != nil {
		c.stream.Stop()
	}
	cfg := &twitterConfig{}
	if err := mapstructure.Decode(config, cfg); err != nil {
		return err
	}
	c.config = cfg
	return c.setupNewStream()
}

// Process gets called when a new event for a specific input should be processed
func (c *Component) Process(inputID string, event interface{}) error {
	// this shouldnt process anything
	return nil
}

func (c *Component) setupNewStream() error {
	oauthConfig := oauth1.NewConfig(c.config.ConsumerKey, c.config.ConsumerSecret)
	token := oauth1.NewToken(c.config.AccessToken, c.config.AccessSecret)
	httpClient := oauthConfig.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	filterParams := &twitter.StreamFilterParams{
		Track:         c.config.Track,
		StallWarnings: twitter.Bool(true),
	}

	// Convenience Demux demultiplexed stream messages
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		fmt.Println(tweet.Text)
		data, _ := json.Marshal(tweet)
		var dataMap map[string]interface{}
		json.Unmarshal(data, &dataMap)
		for id := range c.Outputs {
			c.Emit(id, dataMap)
		}
	}
	demux.DM = func(dm *twitter.DirectMessage) {
		fmt.Println(dm.SenderID)
	}
	demux.Event = func(event *twitter.Event) {
		fmt.Printf("%#v\n", event)
	}

	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		return err
	}
	c.stream = stream
	go func() {
		demux.HandleChan(stream.Messages)
		log.Print("demux returned")
	}()
	return nil
}

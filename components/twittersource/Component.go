package twittersource

import (
	"encoding/json"
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
	if c.stream != nil {
		c.stream.Stop()
	}
	cfg := &twitterConfig{}
	if err := mapstructure.Decode(config, cfg); err != nil {
		return err
	}
	c.config = cfg
	if err := c.setupNewStream(); err != nil {
		return err
	}
	log.Print("stream setup completed, waiting for tweets...")
	go func() {
		for evt := range c.stream.Messages {
			if tweet, ok := evt.(*twitter.Tweet); ok {
				log.Printf("got tweet from %v: %v", tweet.User.Name, tweet.Text)
				data, _ := json.Marshal(tweet)
				var dataMap map[string]interface{}
				json.Unmarshal(data, &dataMap)
				for id := range c.Outputs {
					c.Emit(id, dataMap)
				}
			}
		}
	}()
	return nil
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
	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		return err
	}
	c.stream = stream
	return nil
}

package twittersource

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/trusch/horst"
	"github.com/trusch/horst/processors/base"
	"github.com/trusch/horst/registry"
)

type twittersourceType struct {
	base.Base
	consumerKey    string
	consumerSecret string
	accessToken    string
	accessSecret   string
	track          []string
	stream         *twitter.Stream
}

func (twittersource *twittersourceType) backend() {
	config := oauth1.NewConfig(twittersource.consumerKey, twittersource.consumerSecret)
	token := oauth1.NewToken(twittersource.accessToken, twittersource.accessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter Client
	client := twitter.NewClient(httpClient)

	// FILTER
	filterParams := &twitter.StreamFilterParams{
		Track:         twittersource.track,
		StallWarnings: twitter.Bool(true),
	}
	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	twittersource.stream = stream

	for evt := range stream.Messages {
		if tweet, ok := evt.(*twitter.Tweet); ok {
			data, _ := json.Marshal(tweet)
			var dataMap map[string]interface{}
			json.Unmarshal(data, &dataMap)
			twittersource.Manager.Emit(twittersource.ID, "out", dataMap)
		}
	}
}

func (twittersource *twittersourceType) Process(in string, data interface{}) {}

func (twittersource *twittersourceType) Stop() {
	twittersource.stream.Stop()
}

func (twittersource *twittersourceType) getKeysFromConfig() error {
	if cfgObject, ok := twittersource.Config.(map[string]interface{}); ok {
		if consumerKey, ok := cfgObject["consumerKey"].(string); ok {
			twittersource.consumerKey = consumerKey
		} else {
			return errors.New("malformed consumerKey")
		}
		if consumerSecret, ok := cfgObject["consumerSecret"].(string); ok {
			twittersource.consumerSecret = consumerSecret
		} else {
			return errors.New("malformed consumerSecret")
		}
		if accessToken, ok := cfgObject["accessToken"].(string); ok {
			twittersource.accessToken = accessToken
		} else {
			return errors.New("malformed accessToken")
		}
		if accessSecret, ok := cfgObject["accessSecret"].(string); ok {
			twittersource.accessSecret = accessSecret
		} else {
			return errors.New("malformed accessSecret")
		}
		if track, ok := cfgObject["track"].([]interface{}); ok {
			var trackItems []string
			for _, val := range track {
				if trackItem, ok := val.(string); ok {
					trackItems = append(trackItems, trackItem)
				} else {
					return fmt.Errorf("need track as []string, got %T inside of the array", val)

				}
			}
			twittersource.track = trackItems
		} else {
			return fmt.Errorf("need track as []string, got %T", cfgObject["track"])
		}
	}
	if twittersource.consumerKey == "" || twittersource.consumerSecret == "" || twittersource.accessToken == "" || twittersource.accessSecret == "" {
		return errors.New("need twitter credentials")
	}
	return nil
}

func init() {
	registry.Register("github.com/trusch/horst/processors/twittersource", func(id string, config interface{}, mgr horst.ProcessorManager) (horst.Processor, error) {
		twittersource := &twittersourceType{}
		twittersource.InitBase(id, config, mgr)
		err := twittersource.getKeysFromConfig()
		if err != nil {
			return nil, err
		}
		go twittersource.backend()
		return twittersource, nil
	})
}

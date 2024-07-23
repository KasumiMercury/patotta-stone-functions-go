package config

import (
	_ "embed"
	"encoding/json"
)

//go:embed target.json
var target string

type Config struct {
	Target Target
}

type Target struct {
	Channel []Channel `json:"channel"`
}

type Channel struct {
	Display   string `json:"display"`
	ChannelId string `json:"channelId"`
}

func NewConfig() (*Config, error) {
	c := &Config{}
	if err := c.loadTarget(target); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) loadTarget(j string) error {
	// load target channels from target.json
	t := &Target{}
	if err := json.Unmarshal([]byte(j), t); err != nil {
		return err
	}

	c.Target = *t

	return nil
}

func (c *Config) ChannelIDs() []string {
	cids := make([]string, 0, len(c.Target.Channel))
	for _, c := range c.Target.Channel {
		cids = append(cids, c.ChannelId)
	}

	return cids
}

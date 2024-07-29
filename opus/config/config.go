package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

//go:embed target.json
var target string

type Config struct {
	Target Target
	Api    Api
}

type Api struct {
	ApiKey string
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
		return nil, fmt.Errorf("failed to load target: %w", err)
	}
	if err := c.loadEnv(); err != nil {
		return nil, fmt.Errorf("failed to load env: %w", err)
	}

	return c, nil
}

func (c *Config) loadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	c.Api.ApiKey = os.Getenv("API_KEY")

	return nil
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

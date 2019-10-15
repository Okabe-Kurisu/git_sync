package main

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Config struct {
	Config *globalConfig `yaml:"common"`
	Repos  []*syncConfig `yaml:"repos"`
}

type globalConfig struct {
	IsSyncTags bool `yaml:"is_sync_tags"`
	IsForce    bool `yaml:"is_force"`
	Frequency  int  `yaml:"frequency"`
}

type syncConfig struct {
	globalConfig `yaml:",inline"`
	Origin       string `yaml:"origin"`
	OriginBranch string `yaml:"origin_branch"`
	Target       string `yaml:"target"`
	TargetBranch string `yaml:"target_branch"`
}



func (c *Config) GetConfig() (*Config, error) {
	yamlFile, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return nil, err
	}

	for _, s := range c.Repos {
		if s.Frequency == 0 {
			s.Frequency = c.Config.Frequency
		}
	}

	return c, nil
}

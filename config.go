package main

import (
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Config *globalConfig `yaml:"common"`
	Repos  []*syncConfig `yaml:"repos"`
	Auths  []*authConfig `yaml:"auths"`
}

type globalConfig struct {
	IsSyncTags bool   `yaml:"is_sync_tags"`
	IsForce    bool   `yaml:"is_force"`
	Frequency  string `yaml:"frequency"`
}

type syncConfig struct {
	globalConfig `yaml:",inline"`
	Origin       string `yaml:"origin"`
	OriginBranch string `yaml:"origin_branch"`
	OriginAuth   authConfig
	Target       string `yaml:"target"`
	TargetBranch string `yaml:"target_branch"`
	TargetAuth   authConfig
}

type authConfig struct {
	Group    string `yaml:"group"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func GetConfig() (*Config, error) {
	var c = &Config{}
	yamlFile, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return nil, err
	}

	for _, s := range c.Repos {
		if s.Frequency == "" {
			s.Frequency = c.Config.Frequency
		}

		if strings.Contains(s.Origin, "git@") || strings.Contains(s.Target, "git@") {
			panic("not supported git protocol")
		}

		for _, a := range Reverse(c.Auths) {
			if strings.Contains(s.Origin, a.Group) {
				s.OriginAuth = *a
			} else if strings.Contains(s.Target, a.Group) {
				s.TargetAuth = *a
			}
		}
	}

	return c, nil
}

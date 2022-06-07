package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ServerPort  string `yaml:"server_port"`
	DatabaseURL string `yaml:"database_url"`
}

func Load(file string) (*Config, error) {
	c := Config{}
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(bytes, &c); err != nil {
		return nil, err
	}
	return &c, err
}

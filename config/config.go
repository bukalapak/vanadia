package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	SchemeToEnv bool `yaml:"SchemeToEnv,omitempty"`
	HostToEnv   struct {
		Segments int `yaml:"segments,omitempty"`
	} `yaml:"HostToEnv,omitempty"`
	AuthTokenToEnv bool `yaml:"AuthTokenToEnv,omitempty"`
}

func FromFile(filename string) (Config, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}

	config := Config{}
	err = yaml.UnmarshalStrict(file, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

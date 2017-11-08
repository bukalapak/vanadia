package config

import (
	"io/ioutil"

	"github.com/bukalapak/vanadia/postman"
	"gopkg.in/yaml.v2"
)

var DefaultConfig = Config{}

type Config struct {
	SchemeToEnv string `yaml:"SchemeToEnv,omitempty"`
	HostToEnv   struct {
		Segments int `yaml:"segments,omitempty"`
	} `yaml:"HostToEnv,omitempty"`
	AuthTokenToEnv string           `yaml:"AuthTokenToEnv,omitempty"`
	GlobalHeaders  []postman.Header `yaml:"GlobalHeaders,omitempty"`
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

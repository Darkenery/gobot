package config

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"time"
)

type Cfg struct {
	Bot   botConfig `yaml:"bot"`
}

type botConfig struct {
	ApiConfig                       apiConfig                       `yaml:"api"`
	UpdateGetter                    updateGetterConfig              `yaml:"update_getter"`
	GenerateRandomTextCommandConfig generateRandomTextCommandConfig `yaml:"generate_random_text"`
}

type apiConfig struct {
	Key              string        `yaml:"key"`
	Url              string        `yaml:"url"`
	Timeout          time.Duration `yaml:"timeout"`
	KeepAlive        time.Duration `yaml:"keepalive"`
	HandshakeTimeout time.Duration `yaml:"handshake_timeout"`
}

type updateGetterConfig struct {
	Timeout        int      `yaml:"timeout"`
	Limit          int      `yaml:"limit"`
	AllowedUpdates []string `yaml:"allowed_updates"`
}

type generateRandomTextCommandConfig struct {
	WordLimit int `yaml:"word_limit"`
}

func LoadConfig(path string) (Cfg, error) {
	bytes, err := ioutil.ReadFile(path)

	cfg := Cfg{}

	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(bytes, &cfg)

	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

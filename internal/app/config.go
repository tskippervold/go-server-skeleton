package app

import (
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

/*
	Application configuration is defined in `config/` yaml files.
	It's possible to extend this with reading env vars as well.
	https://dev.to/ilyakaznacheev/a-clean-way-to-pass-configs-in-a-go-application-1g64
*/

type Config struct {
	Server struct {
		Port              string        `yaml:"port" envconfig:"PORT"`
		ConnectionTimeout time.Duration `yaml:"connectionTimeout" envconfig:"CONN_TIMEOUT"`
	} `yaml:"server"`

	Database struct {
		Host string `yaml:"host"`
		Port uint16 `yaml:"port"`
		Name string `yaml:"name"`
		User string `yaml:"user"`
		Pass string `yaml:"pass"`
	}
}

// LoadConfig reads specified yaml file and returns a `Config` struct.
func LoadConfig(path string) (Config, error) {
	if len(path) <= 0 {
		var cfg Config
		err := envconfig.Process("", &cfg)
		return cfg, err
	}

	f, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)

	return cfg, err
}

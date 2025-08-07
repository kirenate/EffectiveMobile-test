package settings

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

func NewConfig() (*Config, error) {
	config := &Config{}
	file, err := os.Open(".env/.yaml")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	d := yaml.NewDecoder(file)
	err = d.Decode(&config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

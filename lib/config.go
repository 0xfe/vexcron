package lib

import (
	"fmt"
	"io/ioutil"
	"log"
)

type Config struct {
	fileName string
	env      Env
	entries  []Entry
	stats    Stats
}

func NewConfig() *Config {
	return &Config{
		env:     make(Env),
		entries: make([]Entry, 0),
	}
}

func (cfg *Config) Load(fileName string) error {
	contents, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	cfg.fileName = fileName
	return cfg.Parse(string(contents))
}

func (cfg *Config) Parse(data string) error {
	log.Print("Starting config parser.")
	entries, env, stats, err := ParseConfig(data)
	if err != nil {
		return fmt.Errorf("Config error: %v", err)
	}

	cfg.entries = entries
	cfg.env = env
	cfg.stats = stats

	return nil
}

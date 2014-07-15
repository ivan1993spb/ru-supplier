package main

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	fname         string
	Host, Port    string
	FilterEnabled bool
}

func LoadConfig(fname string) (conf *Config, err error) {
	var file *os.File
	conf = new(Config)
	*conf = *defaultConfig
	file, err = os.Open(fname)
	if err != nil {
		if os.IsNotExist(err) {
			err = errors.New("config file was not found")
		}
	} else {
		defer file.Close()
		dec := json.NewDecoder(file)
		if err = dec.Decode(&conf); err == nil {
			if !conf.Valid() {
				*conf = *defaultConfig
				err = errors.New("invalid config")
			}
		}
	}
	conf.fname = fname
	return
}

func (c *Config) Save() error {
	file, err := os.Create(c.fname)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	return enc.Encode(c)
}

func (c *Config) Valid() bool {
	return len(c.Host)*len(c.Port) != 0
}

var defaultConfig = &Config{
	Host:          "proxy-zakupki-gov-ru.local",
	Port:          "80",
	FilterEnabled: true,
}

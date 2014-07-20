package main

import (
	"encoding/json"
	"errors"
	"os"
)

var ErrInvalidConfig = errors.New("invalid config")

// Config contains configurations
// If you want use ptogram with any rss client port must be 80
// (some rss clients require this)
type Config struct {
	fname         string
	Host, Port    string
	FilterEnabled bool
}

var defaultConfig = &Config{
	Host:          "proxy-zakupki-gov-ru.local",
	Port:          "80",
	FilterEnabled: true,
}

func LoadConfig(fname string) (conf *Config, err error) {
	if len(fname) == 0 {
		panic("config: invalid file name")
	}
	conf = new(Config)
	*conf = *defaultConfig
	var file *os.File
	file, err = os.Open(fname)
	if err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
	} else {
		defer file.Close()
		dec := json.NewDecoder(file)
		if err = dec.Decode(&conf); err == nil {
			if !conf.Valid() {
				*conf = *defaultConfig
				err = ErrInvalidConfig
			}
		}
	}
	conf.fname = fname
	return
}

func (c *Config) Save() error {
	if !c.Valid() {
		return ErrInvalidConfig
	}
	if c.LikeDefault() {
		os.Remove(c.fname)
		return nil
	}
	file, err := os.Create(c.fname)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	return enc.Encode(c)
}

func (c *Config) LikeDefault() bool {
	return c.Host == defaultConfig.Host &&
		c.Port == defaultConfig.Port &&
		c.FilterEnabled == defaultConfig.FilterEnabled
}

func (c *Config) Valid() bool {
	return len(c.Host)*len(c.Port) > 0
}

func (c *Config) HTTPHost() (host string) {
	host = c.Host
	if c.Port != "80" {
		host += ":" + c.Port
	}
	return
}

func (c *Config) SetFilterEnabled(flag bool) {
	c.FilterEnabled = flag
}

func (c *Config) SetHost(host string) {
	c.Host = host
}

func (c *Config) SetPort(port string) {
	c.Port = port
}

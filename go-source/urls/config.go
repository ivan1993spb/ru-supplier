package main

import (
	"encoding/json"
	"os"
)

const _CONFIG_FILE = "config.json"

const _LOCAL_PROXY_DEF_HOST = "proxy-zakupki-gov-ru.local"

// GetHttpHost tries load configs from file and returns default value
// on failure
func GetHttpHost() string {
	file, err := os.Open(_CONFIG_FILE)
	if err == nil {
		defer file.Close()

		var conf *struct {
			Host, Port string
		}

		if err = json.NewDecoder(file).Decode(&conf); err == nil {
			if len(conf.Host)*len(conf.Port) > 0 {
				var httpHost = conf.Host

				if conf.Port != "80" {
					httpHost += ":" + conf.Port
				}

				return httpHost
			}
		}
	}

	return _LOCAL_PROXY_DEF_HOST
}

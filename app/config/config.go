package config

import (
	"encoding/json"
	"log"
	"strings"

	atomic_value "github.com/aicacia/go-atomic-value"
	"github.com/aicacia/pickup/app/repository"
	"github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
)

var config atomic_value.AtomicValue[*ConfigST]

func Get() *ConfigST {
	return config.Load()
}

type ConfigST struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	URI  string `json:"uri"`
	UI   struct {
		URI string `json:"uri"`
	} `json:"ui"`
	Dashboard struct {
		Enabled bool `json:"enabled"`
	} `json:"dashboard"`
	OpenAPI struct {
		Enabled bool `json:"enabled"`
	} `json:"openapi"`
	JWT struct {
		Secret  string `json:"secret"`
		Expires struct {
			Seconds int `json:"seconds"`
		} `json:"expires"`
	} `json:"jwt"`
}

func InitConfig() error {
	configs, err := repository.GetConfigs()
	if err != nil {
		return err
	}
	configJSON := make(map[string]interface{})
	for _, config := range configs {
		var value interface{}
		err := json.Unmarshal([]byte(config.Value), &value)
		if err != nil {
			log.Printf("invalid json %s: %v", config.Key, err)
			continue
		}
		setKeyValue(configJSON, config.Key, value)
	}
	var c ConfigST
	err = mapstructure.Decode(configJSON, &c)
	if err != nil {
		return err
	}
	config.Store(&c)

	listener, err := repository.CreateListener("configs_channel")
	if err != nil {
		return err
	}
	go configListener(listener)
	return nil
}

func CloseConfigListener() error {
	configListenerSignal <- true
	return nil
}

var configListenerSignal = make(chan bool, 1)

func configListener(listener *pq.Listener) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered in configListener: %v\n", r)
		}
	}()
	for {
		select {
		case <-configListenerSignal:
			err := listener.Close()
			if err != nil {
				log.Printf("error closing config listener: %v\n", err)
				return
			} else {
				log.Printf("config listener closed\n")
			}
			return
		case notification := <-listener.Notify:
			var extra ExtraST
			err := json.Unmarshal([]byte(notification.Extra), &extra)
			if err != nil {
				log.Printf("invalid json %s: %v", notification.Extra, err)
			} else {
				updateConfig(extra.Key, extra.Value)
			}
		}
	}
}

type ExtraST struct {
	Table      string      `json:"table"`
	Key        string      `json:"key"`
	Value      interface{} `json:"value"`
	ActionType string      `json:"action_type"`
}

func updateConfig(key string, value interface{}) error {
	configJSON := make(map[string]interface{})
	setKeyValue(configJSON, key, value)

	c := *Get()
	err := mapstructure.Decode(configJSON, &c)
	if err != nil {
		return err
	}
	config.Store(&c)

	return nil
}

func setKeyValue(parent map[string]interface{}, key string, value interface{}) {
	entry := parent
	path := strings.Split(key, ".")
	for _, key := range path[:len(path)-1] {
		subEntry, ok := entry[key]
		if !ok {
			subEntry = make(map[string]interface{})
			entry[key] = subEntry
		}
		entry = subEntry.(map[string]interface{})
	}
	k := path[len(path)-1]
	entry[k] = value
	log.Printf("%s = %v\n", key, value)
}

package config

import (
	"errors"
	"sync"
)

//Config main storage
type Config struct {
	cfgStorage map[string]interface{}
	mutex      sync.Mutex
}

// Option add new options to config
type Option func(cfg *Config)

var onceCfg sync.Once
var cfgInstance *Config

// GetConfig get singleton instance of Config
func GetConfig() *Config {
	onceCfg.Do(func() {
		cfgInstance = new(Config)
		cfgInstance.init()
	})
	return cfgInstance
}

// Set a value to key
func (c *Config) Set(key string, value interface{}) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cfgStorage[key] = value
	return true
}

func (c *Config) get(key string) (interface{}, error) {
	if v, ok := c.cfgStorage[key]; ok {
		return v, nil
	}
	return nil, errors.New("This key does not exit")
}

func (c *Config) init() {
	c.cfgStorage = make(map[string]interface{})
}

// GetBool get boolean value from given key
func (c *Config) GetBool(key string) bool {
	v, err := c.get(key)
	if err == nil {
		if rv, ok := v.(bool); ok {
			return rv
		}
	}
	return false
}

// GetInt get int value from given key
func (c *Config) GetInt(key string) int {
	v, err := c.get(key)
	if err == nil {
		if rv, ok := v.(int); ok {
			return rv
		}
	}
	return 0
}

// GetUint get unsigned int value from given key
func (c *Config) GetUint(key string) uint {
	v, err := c.get(key)
	if err == nil {
		if rv, ok := v.(uint); ok {
			return rv
		}
	}
	return 0
}

// GetString get string value from given key
func (c *Config) GetString(key string) string {
	v, err := c.get(key)
	if err == nil {
		if rv, ok := v.(string); ok {
			return rv
		}
	}
	return ""
}

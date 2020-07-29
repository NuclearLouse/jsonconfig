package jsonconfig

import (
	"encoding/json"

	"os"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"
)

// Config managment structure
type Config struct {
	sync.RWMutex
	m map[string]interface{}
}

// ReadConfig reads a configuration file and returns the managment structure
func ReadConfig(configFile string) (*Config, error) {

	var c Config
	f, err := os.Open(configFile)
	if err != nil {
		return &c, errors.Wrap(err, "read file")
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&c.m); err != nil {
		return &c, errors.Wrap(err, "decode json")
	}
	return &c, nil
}

// Get main function of getting value
func (c *Config) Get(section string, key string) interface{} {
	c.RLock()
	tmp := c.m[section]
	var res interface{}
	if tmp != nil {
		res = tmp.(map[string]interface{})[key]
	} else {
		res = nil
	}
	c.RUnlock()
	return res
}

// GetAsString return string value
func (c *Config) GetAsString(section string, key string) string {
	val := c.Get(section, key)

	switch val.(type) {
	case string:
		return val.(string)
	case int:
		return strconv.Itoa(val.(int))
	case int32:
		return strconv.FormatInt(int64(val.(int32)), 10)
	case int64:
		return strconv.FormatInt(val.(int64), 10)
	case uint32:
		return strconv.FormatUint(uint64(val.(uint32)), 10)
	case uint64:
		return strconv.FormatUint(val.(uint64), 10)
	case bool:
		return strconv.FormatBool(val.(bool))
	case float32:
		return strconv.FormatFloat(float64(val.(float32)), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(val.(float64), 'f', -1, 64)
	case time.Duration:
		return val.(time.Duration).String()
	case time.Time:
		return val.(time.Time).Format(time.RFC3339)
	default:
		return ""
	}
}

// GetAsFloat64 return float64 value
func (c *Config) GetAsFloat64(section string, key string) (float64, error) {
	val := c.Get(section, key)

	switch val.(type) {
	case string:
		res, err := strconv.ParseFloat(val.(string), 64)
		return res, err
	case float32:
		return float64(val.(float32)), nil
	case float64:
		return val.(float64), nil
	default:
		return 0, errors.New("convert to float64 in: " + section + "." + key)
	}
}

// GetAsInt return integer value
func (c *Config) GetAsInt(section string, key string) (int, error) {
	val := c.Get(section, key)

	switch val.(type) {
	case string:
		res, err := strconv.Atoi(val.(string))
		return res, err
	case int:
		return val.(int), nil
	case int32:
		return int(val.(int32)), nil
	case int64:
		return int(val.(int64)), nil
	case uint32:
		return int(val.(uint32)), nil
	case uint64:
		return int(val.(uint64)), nil
	case bool:
		if val.(bool) {
			return 1, nil
		}
		return 0, nil
	case float32:
		return int(val.(float32)), nil
	case float64:
		return int(val.(float64)), nil
	default:
		return 0, errors.New("convert to int in: " + section + "." + key)
	}
}

// SetValue function of setting value
func (c *Config) SetValue(section string, key string, value interface{}) {
	c.Lock()
	tmp := c.m[section]
	if tmp == nil {
		tmp = make(map[string]interface{})
		c.m[section] = tmp
	}
	tmp.(map[string]interface{})[key] = value
	c.Unlock()
}

//  Copyright jean-françois PHILIPPE 2014-2018

package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config Represente une configuration
// les elements de configurations sont imbricables
type Config struct {
	values map[string]interface{}
	parent *Config
}

// New Create a new Conf
// As parameter, take a list of map[string]string that will be loaded as default
// or Strings "key=value" that will be parsed by Load
func New(defaults map[string]string) *Config {
	result := &Config{values: make(map[string]interface{})}
	for k, v := range defaults {
		result.values[k] = v
	}

	return result
}

// Sections get map of Configs
func (c *Config) Sections() map[string]*Config {
	result := make(map[string]*Config)
	for key, val := range c.values {
		switch entry := val.(type) {
		case *map[string]interface{}:
			result[key] = &Config{*entry, c}
		}
	}
	return result
}

// Section Acces a une sous section
// retourne une section "vide" si elle n existe pas.
func (c *Config) Section(key string) *Config {
	keys := strings.Split(key, ".")
	m := c.sectionA(keys, false)
	if m != nil {
		return &Config{*m, c}
	}
	return &Config{make(map[string]interface{}), c}
}

func (c *Config) sectionA(keys []string, create bool) *map[string]interface{} {
	vals := c.values
	for _, k := range keys {
		k := strings.TrimSpace(k)
		if k != "" {
			sub, ok := vals[k]
			if ok {
				// Entree presente, je verifie son type !
				switch entry := sub.(type) {
				case string:
					if create {
						entry := make(map[string]interface{})
						vals[k] = &entry
						vals = entry
					} else {
						return nil
					}
				case *map[string]interface{}:
					vals = *entry
				}
			} else {
				// Entree abssente
				if create {
					entry := make(map[string]interface{})
					vals[k] = &entry
					vals = entry
				} else {
					return nil
				}
			}
		}
	}

	return &vals
}

// Raw Acces a la valeur 'brute'
func (c *Config) Raw(key string) (raw string, exists bool) {
	keys := strings.Split(key, ".")
	section := keys[:len(keys)-1]
	name := keys[len(keys)-1]
	entries := c.sectionA(section, false)
	if entries != nil {
		item, found := (*entries)[name]
		if found {
			switch value := item.(type) {
			case string:
				return value, true
			default:
				return "", false
			}
		}
	}
	return "", false
}

// Find Recherche une valeur en remontant eventuellement dans les configs parents
func (c *Config) Find(key string) (raw string, exists bool) {
	keys := strings.Split(key, ".")
	section := keys[:len(keys)-1]
	name := keys[len(keys)-1]
	conf := c
	var entries *map[string]interface{}
	for conf != nil {
		entries = conf.sectionA(section, false)
		if entries != nil {
			item, found := (*entries)[name]
			if found {
				switch value := item.(type) {
				case string:
					return value, true
				default:
					return "", false
				}
			}
		}
		conf = conf.parent
	}
	return os.Getenv(key), false
}

// String Recupere une valeur sous forme de chaine.
// Indique la clef de la valeur a recuperer et une valeur par defaut si non definie.
func (c *Config) String(key string, deflt ...string) (string, error) {
	raw, ok := c.Raw(key)
	// If not exists, return default value
	if !ok {
		if len(deflt) > 0 {
			raw = deflt[0]
		} else {
			return "", errors.New("Key '" + key + "' does not exsists")
		}
	}
	return c.Eval(raw)
}

// Bool Recupere une valeur booleenne.
// Indique la clef de la val a recuperer et une valeur par defaut.
// La valeur par defaut peut etre precisee sous forme bool ou sous la forme d une chaine
// qui sera evaluee avant conversion.
func (c *Config) Bool(key string, deflt ...interface{}) (bool, error) {
	raw, ok := c.Raw(key)
	// If not exists, return default valuea or false
	if !ok {
		if len(deflt) > 0 {
			switch val := deflt[0].(type) {
			case bool:
				return val, nil
			case string:
				raw = val // Will be evaluated below !
			default:
				return false, errors.New("Key '" + key + "' does not exsists, default not valid")
			}
		} else {
			return false, errors.New("Key '" + key + "' does not exsists")
		}
	}

	raw, err := c.Eval(raw)
	if err == nil {
		return strconv.ParseBool(raw)
	}

	return false, err
}

// Int64 Recupere une valeur booleenne.
// Indique la clef de la val a recuperer et une valeur par defaut.
// La valeur par defaut peut etre precisee sous forme nombre ou sous la forme d une chaine
// qui sera evaluee avant conversion.
func (c *Config) Int64(key string, deflt ...interface{}) (int64, error) {
	raw, ok := c.Raw(key)
	// If not exists, return default valuea or false
	if !ok {
		if len(deflt) > 0 {
			switch val := deflt[0].(type) {
			case int:
				return int64(val), nil
			case uint:
				return int64(val), nil
			case int8:
				return int64(val), nil
			case uint8:
				return int64(val), nil
			case int16:
				return int64(val), nil
			case uint16:
				return int64(val), nil
			case int32:
				return int64(val), nil
			case uint32:
				return int64(val), nil
			case int64:
				return int64(val), nil
			case uint64:
				return int64(val), nil
			case string:
				raw = val // Will be evaluated below !
			default:
				return 0, errors.New("Key '" + key + "' does not exsists, default not valid")
			}
		} else {
			return 0, errors.New("Key '" + key + "' does not exsists")
		}
	}

	raw, err := c.Eval(raw)
	if err == nil {
		return strconv.ParseInt(raw, 0, 64)
	}

	return 0, err
}

// Uint64 Recupere une valeur booleenne.
// Indique la clef de la val a recuperer et une valeur par defaut.
// La valeur par defaut peut etre precisee sous forme nombre ou sous la forme d une chaine
// qui sera evaluee avant conversion.
func (c *Config) Uint64(key string, deflt ...interface{}) (uint64, error) {
	raw, ok := c.Raw(key)
	// If not exists, return default valuea or false
	if !ok {
		if len(deflt) > 0 {
			switch val := deflt[0].(type) {
			case int:
				return uint64(val), nil
			case uint:
				return uint64(val), nil
			case int8:
				return uint64(val), nil
			case uint8:
				return uint64(val), nil
			case int16:
				return uint64(val), nil
			case uint16:
				return uint64(val), nil
			case int32:
				return uint64(val), nil
			case uint32:
				return uint64(val), nil
			case int64:
				return uint64(val), nil
			case uint64:
				return uint64(val), nil
			case string:
				raw = val // Will be evaluated below !
			default:
				return 0, errors.New("Key '" + key + "' does not exsists, default not valid")
			}
		} else {
			return 0, errors.New("Key '" + key + "' does not exsists")
		}
	}

	raw, err := c.Eval(raw)
	if err == nil {
		return strconv.ParseUint(raw, 0, 64)
	}

	return 0, err
}

// Duration Return a Duration; Valid values are those that are conformed to time.ParseDuration.
func (c *Config) Duration(key string, deflt ...interface{}) (time.Duration, error) {
	raw, ok := c.Raw(key)
	// If not exists, return default valuea or false
	if !ok {
		if len(deflt) > 0 {
			switch val := deflt[0].(type) {
			case time.Duration:
				return val, nil
			case string:
				raw = val // Will be evaluated below !
			default:
				return 0, errors.New("Key '" + key + "' does not exsists, default not valid")
			}
		} else {
			return 0, errors.New("Key '" + key + "' does not exsists")
		}
	}

	raw, err := c.Eval(raw)
	if err == nil {
		return time.ParseDuration(raw)
	}

	return 0 * time.Second, err
}

/* vi:set fileencodings=utf-8 tabstop=4 ai: */

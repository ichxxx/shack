package config

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/spf13/viper"
)

type configManager struct {
	core    *viper.Viper
	configs map[string]config
	mutex   sync.Mutex
}

// Add adds a config will be loaded.
func Add(config config, section string) {
	config.bind(config, section)
	manager.configs[section] = config
}

// File specify a config file to load.
// Default file is `config.yml`.
func File(file string) *configManager {
	defaultFile = file
	return manager
}

// ParseTag specify the tag to parse the config file.
// Default tag is `config`.
func ParseTag(tag string) *configManager {
	defaultTag = tag
	return manager
}

// Load loads the previously added configs from the config file.
func Load() {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	manager.loadConfig()
	for _, c := range manager.configs {
		c.mapConfig()
	}
	for _, c := range manager.configs {
		c.Init()
	}
}

func Options(opts ...func(v *viper.Viper)) {
	for _, opt := range opts {
		opt(manager.core)
	}
}

func (cm *configManager) loadConfig() {
	cm.core = viper.New()
	cm.core.SetConfigFile(defaultFile)
	err := cm.core.ReadInConfig()
	if err != nil {
		fmt.Printf("load config err: %s\n", err)
	}
}

func (cm *configManager) getFieldValue(key string, rv reflect.Value) (value reflect.Value, err error) {
	switch rv.Kind() {
	case reflect.String:
		value = reflect.ValueOf(cm.core.GetString(key))
	case reflect.Slice:
		switch reflect.TypeOf(rv.Interface()).Elem().Kind() {
		case reflect.String:
			value = reflect.ValueOf(cm.core.GetStringSlice(key))
		case reflect.Int:
			value = reflect.ValueOf(cm.core.GetIntSlice(key))
		}
	case reflect.Bool:
		value = reflect.ValueOf(cm.core.GetBool(key))
	case reflect.Int:
		value = reflect.ValueOf(cm.core.GetInt(key))
	case reflect.Int64:
		value = reflect.ValueOf(cm.core.GetInt64(key))
	case reflect.Uint:
		value = reflect.ValueOf(cm.core.GetUint(key))
	case reflect.Uint64:
		value = reflect.ValueOf(cm.core.GetUint64(key))
	case reflect.Float64:
		value = reflect.ValueOf(cm.core.GetFloat64(key))
	default:
		err = errors.New("shack: parse config error, can't trans value to the specify type")
	}

	return
}

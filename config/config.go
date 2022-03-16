package config

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type configManager struct {
	file    string
	tag     string
	core    *viper.Viper
	configs map[string]config
	mutex   sync.Mutex
}

var Core *viper.Viper

func ParseFlag(name string) *configManager {
	return manager.ParseFlag(name)
}

func (cm *configManager) ParseFlag(name string) *configManager {
	flag.StringVar(&defaultFile, name, "-c", "配置文件路径")
	flag.Parse()
	return cm
}

// Add adds a config will be loaded.
func Add(config config, section string) {
	config.bind(config, section)
	manager.configs[section] = config
}

func (cm *configManager) Add(config config, section string) {
	config.bind(config, section)
	cm.configs[section] = config
}

// File specify a config file to load.
// Default file is `config.yaml`.
func File(file string) *configManager {
	return manager.File(file)
}

func (cm *configManager) File(file string) *configManager {
	cm.file = file
	return cm
}

// ParseTag specify the tag to parse the config file.
// Default tag is `config`.
func ParseTag(tag string) *configManager {
	return manager.ParseTag(tag)
}

func (cm *configManager) ParseTag(tag string) *configManager {
	cm.tag = tag
	return cm
}

// Load loads the previously added configs from the config file.
func Load() {
	manager.Load()
}

func (cm *configManager) Load() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	cm.loadConfig()
	for _, c := range cm.configs {
		c.mapConfig()
	}
	for _, c := range cm.configs {
		c.Init()
	}
}

func Options(opts ...func(v *viper.Viper)) *configManager {
	return manager.Options(opts...)
}

func (cm *configManager) Options(opts ...func(v *viper.Viper)) *configManager {
	for _, opt := range opts {
		opt(cm.core)
	}
	return cm
}

func (cm *configManager) loadConfig() {
	Core = viper.New()
	cm.core = Core
	cm.core.SetConfigFile(defaultFile)
	err := cm.core.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("shack config: load config err: %s", err))
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
	case reflect.Map:
		elem := reflect.TypeOf(rv.Interface()).Elem()
		switch elem.Kind() {
		case reflect.String:
			value = reflect.ValueOf(cm.core.GetStringMapString(key))
		case reflect.Slice:
			value = reflect.ValueOf(cm.core.GetStringMapStringSlice(key))
		case reflect.Interface:
			value = reflect.ValueOf(cm.core.GetStringMap(key))
		}
	case reflect.Interface:
		value = reflect.ValueOf(cm.core.Get(key))
	default:
		switch rv.Interface().(type) {
		case time.Time:
			value = reflect.ValueOf(cm.core.GetTime(key))
		default:
			err = errors.New("parse config error, can't trans value to the specify type")
		}
	}

	return
}

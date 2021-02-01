package shack

import (
	"fmt"
	"reflect"

	"github.com/spf13/viper"
)

var (
	configFile = "config.toml"
	mode       string
	Config     = &configManager{}
)

type BaseConfig struct {
	self    config
	section string
}

type config interface {
	Init()
	bind(config, string)
}

type configManager struct {
	Core    *viper.Viper
	configs []config
}


func(bc *BaseConfig) Init() {
	bc.mapConfig()
}


func(bc *BaseConfig) bind(config config, section string) {
	bc.self = config
	bc.section = section
}


// Add adds a config will be loaded.
func(cm *configManager) Add(config config, section string) {
	config.bind(config, section)
	cm.configs = append(cm.configs, config)
}


// File specify a toml file to load.
// Default file is `config.toml`.
func(cm *configManager) File(file string) *configManager {
	configFile = file
	return cm
}


// Load loads the previously added configs from the toml file.
func(cm *configManager) Load() {
	cm.loadConfig()

	for _, config := range Config.configs {
		config.Init()
	}
}


func(cm *configManager) loadConfig() {
	cm.Core = viper.New()
	cm.Core.SetConfigFile(configFile)
	err := cm.Core.ReadInConfig()
	if err != nil {
		fmt.Printf("load config err: %s\n", err)
	}

	mode = cm.Core.GetString("app.mode")
}


func(cm *configManager) getFieldValue(key string, rv reflect.Value) reflect.Value {
	switch rv.Kind() {
	case reflect.String:
		return reflect.ValueOf(cm.Core.GetString(key))
	case reflect.Slice:
		switch reflect.TypeOf(rv.Interface()).Elem().Kind() {
		case reflect.String:
			return reflect.ValueOf(cm.Core.GetStringSlice(key))
		case reflect.Int:
			return reflect.ValueOf(cm.Core.GetIntSlice(key))
		}
	case reflect.Bool:
		return reflect.ValueOf(cm.Core.GetBool(key))
	case reflect.Int:
		return reflect.ValueOf(cm.Core.GetInt(key))
	case reflect.Int64:
		return reflect.ValueOf(cm.Core.GetInt64(key))
	case reflect.Uint:
		return reflect.ValueOf(cm.Core.GetUint(key))
	case reflect.Uint64:
		return reflect.ValueOf(cm.Core.GetUint64(key))
	case reflect.Float64:
		return reflect.ValueOf(cm.Core.GetFloat64(key))
	}

	return reflect.ValueOf(cm.Core.GetString(key))
}


func(bc *BaseConfig) mapConfig() {
	p := reflect.ValueOf(bc.self)
	if p.Kind() != reflect.Ptr || p.IsNil() {
		panic("shack config: dst is not a pointer")
	}

	rv := p.Elem()
	if rv.Kind() != reflect.Struct {
		panic("shack config: dst is not a struct")
	}

	t := rv.Type()
	size := rv.NumField()
	if size == 0 {
		panic("shack config: dst struct doesn't have any fields")
	}

	for i := 0; i < size; i++ {
		structField := t.Field(i)
		if sFieldKind := structField.Type.Kind(); sFieldKind == reflect.Struct || sFieldKind == reflect.Map {
			continue
		}

		if mode != "release" {
			configField := fmt.Sprintf("%s.%s.%s", bc.section, mode, structField.Name)
			if Config.Core.IsSet(configField) {
				rv.Field(i).Set(
					Config.getFieldValue(configField, rv.Field(i)),
				)
				continue
			}
		}

		configField := fmt.Sprintf("%s.%s", bc.section, structField.Name)
		if Config.Core.IsSet(configField) {
			rv.Field(i).Set(
				Config.getFieldValue(configField, rv.Field(i)),
			)
		}
	}
}

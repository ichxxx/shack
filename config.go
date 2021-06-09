package shack

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/spf13/viper"
)

var (
	tagParse   = "config"
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
// Default file is `test_config.toml`.
func(cm *configManager) File(file string) *configManager {
	configFile = file
	return cm
}


// ParseTag specify the tag to parse the toml file.
// Default tag is `config`.
func(cm *configManager) ParseTag(tag string) *configManager {
	tagParse = tag
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

	mode = cm.Core.GetString("shack.mode")
}


func(cm *configManager) getFieldValue(key string, rv reflect.Value) (value reflect.Value, err error) {
	switch rv.Kind() {
	case reflect.String:
		value = reflect.ValueOf(cm.Core.GetString(key))
	case reflect.Slice:
		switch reflect.TypeOf(rv.Interface()).Elem().Kind() {
		case reflect.String:
			value = reflect.ValueOf(cm.Core.GetStringSlice(key))
		case reflect.Int:
			value = reflect.ValueOf(cm.Core.GetIntSlice(key))
		}
	case reflect.Bool:
		value = reflect.ValueOf(cm.Core.GetBool(key))
	case reflect.Int:
		value = reflect.ValueOf(cm.Core.GetInt(key))
	case reflect.Int64:
		value = reflect.ValueOf(cm.Core.GetInt64(key))
	case reflect.Uint:
		value = reflect.ValueOf(cm.Core.GetUint(key))
	case reflect.Uint64:
		value = reflect.ValueOf(cm.Core.GetUint64(key))
	case reflect.Float64:
		value = reflect.ValueOf(cm.Core.GetFloat64(key))
	default:
		err = errors.New("can't trans to the specify type")
	}

	return
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

		configField, _ := getConfigField(structField)
		if configField == "-" {
			continue
		}

		var prefix string
		if mode != "" {
			prefix = fmt.Sprintf("%s.%s.", bc.section, mode)
			if !Config.Core.IsSet(prefix + configField) {
				prefix = fmt.Sprintf("%s.", bc.section)
			}
		} else {
			prefix = fmt.Sprintf("%s.", bc.section)
		}

		configField = prefix + configField
		if Config.Core.IsSet(configField) {
			fieldValue, err := Config.getFieldValue(configField, rv.Field(i))
			if err == nil {
				rv.Field(i).Set(fieldValue)
			}
		}
	}
}


func getConfigField(structField reflect.StructField) (name string, option string) {
	tag := structField.Tag.Get(tagParse)
	name, option = parseTag(tag)

	if !isValidTag(name) {
		return structField.Name, ""
	}

	return
}


func parseTag(s string) (tag string, option string) {
	if idx := strings.Index(s, ","); idx != -1 {
		return tag[:idx], tag[idx+1:]
	}
	return s, ""
}


func isValidTag(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		switch {
		case strings.ContainsRune("!#$%&()*+-./:<=>?@[]^_{|}~ ", c):
			// Backslash and quote chars are reserved, but
			// otherwise any punctuation chars are allowed
			// in a tag name.
		case !unicode.IsLetter(c) && !unicode.IsDigit(c):
			return false
		}
	}
	return true
}

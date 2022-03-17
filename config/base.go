package config

import (
	"reflect"
	"strings"
	"unicode"
)

const (
	defaultTag  = "config"
	defaultFile = "config.yaml"
)

var (
	manager = &configManager{
		file:    defaultFile,
		tag:     defaultTag,
		configs: map[string]config{},
	}
)

type Base struct {
	config  config
	section string
}

type config interface {
	Init()
	Get(key string) interface{}
	Set(key string, value interface{}) error
	bind(config, string)
	mapConfig(tag string)
}

func (b *Base) Init() {}

func (b *Base) Get(key string) interface{} {
	return manager.core.Get(b.section + "." + key)
}

func (b *Base) Set(key string, value interface{}) error {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	manager.core.Set(b.section+"."+key, value)
	return manager.core.WriteConfig()
}

func (b *Base) bind(config config, section string) {
	b.config = config
	b.section = section
}

func (b *Base) mapConfig(tag string) {
	p := reflect.ValueOf(b.config)
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
		if sFieldKind := structField.Type.Kind(); sFieldKind == reflect.Struct {
			continue
		}

		configField, _ := getConfigField(structField, tag)
		if configField == "-" {
			continue
		}

		configField = b.section + "." + configField
		if manager.core.IsSet(configField) {
			fieldValue, err := manager.getFieldValue(configField, rv.Field(i))
			if err == nil {
				rv.Field(i).Set(fieldValue)
			}
		}
	}
}

func getConfigField(structField reflect.StructField, tag string) (name string, option string) {
	tagValue := structField.Tag.Get(tag)
	name, option = parseTagValue(tagValue)

	if !isValidTagValue(name) {
		return structField.Name, ""
	}

	return
}

func parseTagValue(s string) (tag string, option string) {
	if idx := strings.Index(s, ","); idx != -1 {
		return tag[:idx], tag[idx+1:]
	}
	return s, ""
}

func isValidTagValue(s string) bool {
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

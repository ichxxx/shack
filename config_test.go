package shack

import (
	"testing"
)


type testConfig struct {
	BaseConfig
	Name    string
	Port    int
	FooBar  []int  `config:"foo_bar"`
	Name2   string `json:"name"`
	Port2   int    `json:"port"`
	FooBar2 []int  `json:"foo_bar"`
}

var test = &testConfig{}


func TestParseConfig(t *testing.T) {
	Config.Add(test, "test")
	Config.File("test_config.toml").Load()
	if test.Name != "shack" {
		t.Errorf("expecting result:%s, got:%s", "shack", test.Name)
	}
	if test.Port != 8080 {
		t.Errorf("expecting result:%d, got:%d", 8080, test.Port)
	}
	if len(test.FooBar) != 2 || test.FooBar[0] != 1 || test.FooBar[1] != 2 {
		t.Errorf("expecting result:%v, got:%v", []int{1, 2}, test.FooBar)
	}
}


func TestParseConfigWithTag(t *testing.T) {
	Config.Add(test, "test")
	Config.File("test_config.toml").ParseTag("json").Load()
	if test.Name2 != "shack" {
		t.Errorf("expecting result:%s, got:%s", "shack", test.Name)
	}
	if test.Port2 != 8080 {
		t.Errorf("expecting result:%d, got:%d", 8080, test.Port)
	}
	if len(test.FooBar2) != 2 || test.FooBar2[0] != 1 || test.FooBar2[1] != 2 {
		t.Errorf("expecting result:%v, got:%v", []int{1, 2}, test.FooBar2)
	}
}
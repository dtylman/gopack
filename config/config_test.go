package config

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestLoad(t *testing.T) {
	conf, err := Load("pkg.config.json")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(conf)
}

func TestCreateConf(t *testing.T) {
	var conf PackageOptions
	conf.Files = make(map[string]string)
	conf.Folders = make(map[string]string)
	conf.Files["source"] = "target"
	conf.Folders["empty"] = ""
	conf.Folders["source"] = "target"

	data, _ := json.Marshal(conf)
	fmt.Println(string(data))
}

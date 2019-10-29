package main

import (
	"encoding/json"
	"io/ioutil"
)

//Scripts holds scripts config
type Scripts struct {
	PreInst    string `json:"pre_inst"`
	PostInst   string `json:"post_inst"`
	PreUnInst  string `json:"pre_uninst"`
	PostUnInst string `json:"post_uninst"`
}

//Config holds package configuration details
type Config struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Revision    string            `json:"revision"`
	Arch        string            `json:"arch"`
	Description string            `json:"description"`
	Homepage    string            `json:"homepage"`
	Depends     string            `json:"depends"`
	Section     string            `json:"section"`
	Maintainer  string            `json:"maintainer"`
	Folders     map[string]string `json:"folders"`
	Files       map[string]string `json:"files"`
	Script      Scripts           `json:"scripts"`
}

//Load loads configuration from file
func Load(fileName string) (*Config, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var conf Config
	err = json.Unmarshal(data, &conf)
	return &conf, err
}

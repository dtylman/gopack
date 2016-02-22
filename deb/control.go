package deb

import (
	"bytes"
	"fmt"
	"reflect"
)

//Control represents debian control structure
type Control struct {
	Package      string
	Version      string
	Source       string
	Depends      string
	Architecture string
	Maintainer   string
	Conflicts    string
	Section      string
	Homepage     string
	Description  string
}

//Bytes marshal control structure as bytes
func (c *Control) bytes() []byte {
	buff := new(bytes.Buffer)
	val := reflect.ValueOf(c).Elem()
	for i := 0; i < val.NumField(); i++ {
		value := val.Field(i).String()
		if value != "" {
			buff.WriteString(fmt.Sprintf("%s: %v\n", val.Type().Field(i).Name, value))
		}
	}
	return buff.Bytes()
}

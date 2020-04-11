package configuration

import (
	"encoding/json"
	"io/ioutil"
)

//Import configuration from file
func Import(filename string) (conf *Config, err error) {
	fileContent, err := ioutil.ReadFile(filename)
	if err == nil {
		bufConf := ConfigJSON{}
		if json.Unmarshal(fileContent, &bufConf) == nil {
			conf = bufConf.ToConfig()
		}
	}
	return
}

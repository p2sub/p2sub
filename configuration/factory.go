package configuration

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//Import configuration from file
func Import(filename string) *Config {
	file, err := os.Open(filename)
	defer file.Close()
	if err == nil {
		if data, err := ioutil.ReadAll(file); err == nil {
			confs := ConfigJSON{}
			if json.Unmarshal(data, &confs) == nil {
				return confs.ToConfig()
			}
		}
	}
	return nil
}

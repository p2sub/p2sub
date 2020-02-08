package configuration

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//ConfigItem confuguration structure
type ConfigItem struct {
	Name          string `json:"name"`
	Address       string `json:"address"`
	Nonce         string `json:"nonce"`
	Signature     string `json:"signature"`
	ConfigSerivce string `json:"configService"`
}

//Configs array of ConfigItem
type Configs []ConfigItem

//Import configuration from file
func Import(filename string) *Configs {
	file, err := os.Open(filename)
	defer file.Close()
	if err == nil {
		if data, err := ioutil.ReadAll(file); err == nil {
			config := &Configs{}
			if json.Unmarshal(data, config) == nil {
				return config
			}
		}
	}
	return nil
}

//Export configuration to file
func (c *Configs) Export(filename string) bool {
	file, err := os.Create(filename)
	defer file.Close()
	if err == nil {
		if data, err := json.Marshal(c); err == nil {
			if n, err := file.Write(data); n == len(data) && err == nil {
				return true
			}
		}
	}
	return false
}

//ToString convert configuration to string
func (c *Configs) ToString() string {
	if data, err := json.Marshal(c); err == nil {
		return string(data)
	}
	return "<nil>"
}

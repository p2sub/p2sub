package configuration

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/btcsuite/btcutil/base58"
	"github.com/p2sub/p2sub/address"
)

const (
	//NodeNotary issue certificate and configuration
	NodeNotary uint32 = iota
	//NodeMaster create blocks
	NodeMaster
	//NodeSlaver sync with blockchain
	NodeSlaver
	//NodeObserver only observe
	NodeObserver
)

var nodeStr = map[uint32]string{
	NodeNotary:   "notary",
	NodeMaster:   "master",
	NodeSlaver:   "slaver",
	NodeObserver: "observer",
}
var nodeUint = map[string]uint32{
	"notary":   NodeNotary,
	"master":   NodeMaster,
	"slaver":   NodeSlaver,
	"observer": NodeObserver,
}

//ConfigJSON confuguration structure
type ConfigJSON struct {
	Name          string `json:"name"`
	BindHost      string `json:"host"`
	BindPort      string `json:"port"`
	Seed          string `json:"seed"`
	Issuer        string `json:"issuer"`
	NodeType      string `json:"nodeType"`
	Nonce         string `json:"nonce"`
	Signature     string `json:"signature"`
	ConfigService string `json:"configService"`
}

//Config indexed in memory
type Config struct {
	Name          string
	BindPort      string
	BindHost      string
	NodeAddress   address.Address
	Issuer        []byte
	Nonce         uint32
	NodeType      uint32
	Signature     []byte
	ConfigService string
}

//Export configuration to file
func (conf *Config) Export(filename string) bool {
	file, err := os.Create(filename)
	defer file.Close()
	if err == nil {
		if data, err := json.MarshalIndent(conf.ToJSON(), "", "  "); err == nil {
			if n, err := file.Write(data); n == len(data) && err == nil {
				return true
			}
		}
	}
	return false
}

//ToJSON memory to ConfigJSON
func (conf *Config) ToJSON() *ConfigJSON {
	return &ConfigJSON{
		Name:          conf.Name,
		BindHost:      conf.BindHost,
		BindPort:      conf.BindPort,
		Seed:          base58.Encode(conf.NodeAddress.GetSeed()),
		NodeType:      nodeStr[conf.NodeType],
		Nonce:         strconv.FormatUint(uint64(conf.Nonce), 10),
		Signature:     base58.Encode(conf.Signature),
		Issuer:        base58.Encode(conf.Issuer),
		ConfigService: conf.ConfigService,
	}
}

//ToConfig json structure to memory config
func (confJSON *ConfigJSON) ToConfig() *Config {
	if nonce, err1 := strconv.ParseUint(confJSON.Nonce, 10, 32); err1 == nil {
		tmpConfig := &Config{
			Name:          confJSON.Name,
			BindHost:      confJSON.BindHost,
			BindPort:      confJSON.BindPort,
			NodeAddress:   *address.FromSeed(base58.Decode(confJSON.Seed)),
			Nonce:         uint32(nonce),
			NodeType:      nodeUint[confJSON.NodeType],
			Signature:     base58.Decode(confJSON.Signature),
			ConfigService: confJSON.ConfigService,
		}
		if tmpConfig.Issuer != nil {
			tmpConfig.Issuer = base58.Decode(confJSON.Issuer)
		}
		return tmpConfig

	}
	return nil
}

//ToString convert configuration to string
func (confJSON *ConfigJSON) ToString() string {
	if data, err := json.MarshalIndent(confJSON, "", "  "); err == nil {
		return string(data)
	}
	return "{}"
}

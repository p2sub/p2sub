package main

import (
	"errors"
	"flag"
	"os"
	"strings"
	"sync"

	"github.com/p2sub/p2sub/config"
	"github.com/p2sub/p2sub/logger"
	"go.uber.org/zap"
)

// P2SubConfig conf wrapper for P2Sub
type P2SubConfig struct {
	cfg *config.Config
}

// FlagConfig flags configuration
type FlagConfig struct {
	name        string
	dataType    string
	value       interface{}
	required    bool
	description string
}

var conf *P2SubConfig
var sugar *zap.SugaredLogger
var confOnce sync.Once

// GetP2SubConfig get singleton instance of Config
func GetP2SubConfig() *P2SubConfig {
	confOnce.Do(func() {
		conf = &P2SubConfig{cfg: config.New()}
	})
	return conf
}

// GetKeyFile get key file
func (p *P2SubConfig) GetKeyFile() string {
	return p.cfg.GetString("node::key_file")
}

// SetKeyFile set key file
func (p *P2SubConfig) SetKeyFile(keyFile string) bool {
	return p.cfg.Set("node::key_file", keyFile)
}

// GetBindPort get bind port of current node
func (p *P2SubConfig) GetBindPort() uint {
	return p.cfg.GetUint("node::bind_port")
}

// SetBindPort set bind port of current node
func (p *P2SubConfig) SetBindPort(bindPort uint) bool {
	return p.cfg.Set("node::bind_port", bindPort)
}

// GetBindHost get bind host
func (p *P2SubConfig) GetBindHost() string {
	return p.cfg.GetString("node::bind_host")
}

// SetBindHost set bind host
func (p *P2SubConfig) SetBindHost(bindHost string) bool {
	return p.cfg.Set("node::bind_host", bindHost)
}

// GetDirectConnect get direct connect node's identity
func (p *P2SubConfig) GetDirectConnect() string {
	return p.cfg.GetString("node::direct_connect")
}

// SetDirectConnect set direct connect node's identity
func (p *P2SubConfig) SetDirectConnect(nodeAddress string) bool {
	return p.cfg.Set("node::direct_connect", nodeAddress)
}

// GetDomain get domain of node discovery
func (p *P2SubConfig) GetDomain() string {
	return p.cfg.GetString("node::domain")
}

// SetDomain set domain of node discovery
func (p *P2SubConfig) SetDomain(domain string) bool {
	return p.cfg.Set("node::domain", domain)
}

func (f FlagConfig) valToBool() bool {
	if v, ok := f.value.(bool); ok {
		return v
	}
	return false
}

func (f FlagConfig) valToString() string {
	if v, ok := f.value.(string); ok {
		return v
	}
	return ""
}

func (f FlagConfig) valToInt() int {
	if v, ok := f.value.(int); ok {
		return v
	}
	return 0
}

func (f FlagConfig) valToUint() uint {
	if v, ok := f.value.(uint); ok {
		return v
	}
	return 0
}

func nameToFlag(name string) string {
	parts := strings.Split(name, "::")
	if len(parts) == 2 {
		// node::key-file
		return strings.ReplaceAll(parts[1], "_", "-")
	}
	panic(errors.New("Wrong format of flag name"))
}

// Init common components
func Init() {
	sugar = logger.GetSugarLogger()
	conf = GetP2SubConfig()

	// All flags configuration
	flagConfigs := []FlagConfig{
		{
			name:        "node::key_file",
			dataType:    "string",
			value:       "",
			description: "File name to save/load key configuration",
			required:    true,
		},
		{
			name:        "node::direct_connect",
			dataType:    "string",
			value:       "",
			description: "Direct connect to a given node",
		},
		{
			name:        "node::domain",
			dataType:    "string",
			value:       "P2Sub::alpha::0.0.1",
			description: "Rendezvous string used to discover same node",
		},
		{
			name:        "node::bind_port",
			dataType:    "uint",
			value:       0,
			description: "Bind port of current node",
			required:    true,
		},
		{
			name:        "node::bind_host",
			dataType:    "string",
			value:       "0.0.0.0",
			description: "Bind host of current node",
			required:    true,
		},
	}

	// Transform flag config to arguments
	for _, flagConf := range flagConfigs {
		switch flagConf.dataType {
		case "string":
			flag.String(nameToFlag(flagConf.name), flagConf.valToString(), flagConf.description)
			break
		case "bool":
			flag.Bool(nameToFlag(flagConf.name), flagConf.valToBool(), flagConf.description)
			break
		case "uint":
			flag.Uint(nameToFlag(flagConf.name), flagConf.valToUint(), flagConf.description)
			break
		case "int":
			flag.Int(nameToFlag(flagConf.name), flagConf.valToInt(), flagConf.description)
			break
		}
	}

	// Parse flags
	flag.Parse()

	isFlagOn := make(map[string]bool)

	flag.Visit(func(f *flag.Flag) {
		isFlagOn[f.Name] = true
	})

	//Save configuration
	for _, flagConf := range flagConfigs {
		if flagConf.required && !isFlagOn[nameToFlag(flagConf.name)] {
			flag.Usage()
			os.Exit(1)
		}
		rawValue := flag.Lookup(nameToFlag(flagConf.name)).Value.(flag.Getter).Get()
		if isFlagOn[nameToFlag(flagConf.name)] {
			sugar.Infof("Flag config: %s value: %v", flagConf.name, rawValue)
		}

		conf.cfg.Set(flagConf.name, rawValue)
	}
}

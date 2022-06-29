package utils

import (
	"encoding/json"
	"os"

	"github.com/spf13/viper"
)

const (
	IptablesCommand    string = "iptables"
	Ip6tablesCommand   string = "ip6tables"
	SaveIpv4Command    string = "iptables-save"
	SaveIpv6Command    string = "ip6tables-save"
	RestoreIpv4Command string = "iptables-restore"
	RestoreIpv6Command string = "ip6tables-restore"
)

type Config struct {
	Protocol      string `mapstructure:"PROTOCOL"`
	Host          string `mapstructure:"SCAN_HOST"`
	PortsAmount   int    `mapstructure:"PORTS_AMOUNT"`
	RulesFile     string `mapstructure:"RULES_FILE"`
	RulesBackupV4 string `mapstructure:"SAVE_IPV4_FILE"`
	RulesBackupV6 string `mapstructure:"SAVE_IPV6_FILE"`
}

type ConfigJson struct {
	Iface    string   `json:"interface,omitempty"`
	Protocol string   `json:"protocol,omitempty"`
	Port     int      `json:"port,omitempty"`
	Allow    []string `json:"allow,omitempty"`
	Chain    string   `json:"chain,omitempty"`
	Version  string   `json:"version,omitempty"`
}

func LoadConfigEnv(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	// Use this for export environment to your bash session
	// viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func LoadConfigJson(path string) ([]*ConfigJson, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var conf []*ConfigJson
	if err := json.NewDecoder(file).Decode(&conf); err != nil {
		return nil, err
	}

	return conf, nil
}

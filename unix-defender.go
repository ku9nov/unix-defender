package main

import (
	"flag"
	"log"
	"unix-defender/iptables"
	"unix-defender/scanner"
	"unix-defender/utils"
)

const (
	restoreIpv4Command = "iptables-restore"
	restoreIpv6Command = "ip6tables-restore"
)

var action string

func main() {
	// flag.BoolVar(&remove, "remove", false, "Remove the following firewall rules FOREVER (a very long time)!")
	configEnv, err := utils.LoadConfigEnv(".")
	if err != nil {
		log.Fatal("Cannot load environment config:", err)
	}
	flag.StringVar(&action, "action", "scan", "Use for choice scanning ports or manage iptables. It can be 'scan', 'manage' or 'interfaces'.")
	flag.Parse()
	switch action {
	case "manage":
		iptables.IpTables()
	case "interfaces":
		scanner.LocalAddresses()
	case "restore":
		iptables.RestoreRules(restoreIpv4Command, &configEnv.RulesBackupV4)
		iptables.RestoreRules(restoreIpv6Command, &configEnv.RulesBackupV6)
	default:
		scanner.ScanPorts()
	}
}

package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"
	"unix-defender/checker"
	"unix-defender/iptables"
	"unix-defender/scanner"
	"unix-defender/utils"
)

var action string

func main() {
	// flag.BoolVar(&remove, "remove", false, "Remove the following firewall rules FOREVER (a very long time)!")
	configEnv, err := utils.LoadConfigEnv(".")
	if err != nil {
		log.Fatal("Cannot load environment config:", err)
	}
	flag.StringVar(&action, "action", "default", "Use for choice scanning ports or manage iptables. It can be 'scan', 'manage' or 'interfaces'.")
	flag.Parse()
	switch action {
	case "manage":
		iptables.IpTables()
	case "interfaces":
		scanner.LocalAddresses()
	case "restore":
		iptables.RestoreRules(utils.RestoreIpv4Command, &configEnv.RulesBackupV4)
		iptables.RestoreRules(utils.RestoreIpv6Command, &configEnv.RulesBackupV6)
	case "scan":
		scanner.ScanPorts()
	default:
		fmt.Println("Nothing choiced, start cycle.")
	}

	rand.Seed(time.Now().UnixNano())
	fileNameIpv4 := utils.RandomString(44)
	fileNameIpv6 := utils.RandomString(66)
	utils.SigTerm(fileNameIpv4, fileNameIpv6)
	for {
		time.Sleep(10 * time.Second)
		checker.SaveRulesTmp(utils.SaveIpv4Command, fileNameIpv4, &configEnv.RulesBackupV4)
		checker.SaveRulesTmp(utils.SaveIpv6Command, fileNameIpv6, &configEnv.RulesBackupV6)
		fmt.Println("Tmp file is saved. Name: ", fileNameIpv4, "and", fileNameIpv6)
	}

}

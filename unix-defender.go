package main

import (
	"flag"
	"io/ioutil"
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
	rand.Seed(time.Now().UnixNano())
	configEnv, err := utils.LoadConfigEnv(utils.EnvFile)
	if err != nil {
		log.Fatal("Cannot load environment config:", err)
	}
	flag.StringVar(&action, "action", "default", "Use for choice scanning ports or manage iptables. It can be 'scan', 'manage' or 'interfaces'.")
	flag.Parse()
	if !configEnv.LoggingEnable {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
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
		log.Println("Unix-defender is successfully running in foreground mode.")
		utils.SendMessageToSlack(utils.StartMessage, utils.BlueColor)
		fileNameIpv4 := utils.RandomString(44)
		fileNameIpv6 := utils.RandomString(66)
		utils.SigTerm(fileNameIpv4, fileNameIpv6)
		for {
			time.Sleep(10 * time.Second)
			checker.SaveRulesTmp(utils.SaveIpv4Command, fileNameIpv4, &configEnv.RulesBackupV4)
			checker.SaveRulesTmp(utils.SaveIpv6Command, fileNameIpv6, &configEnv.RulesBackupV6)
			log.Println("Tmp file is saved. Name: ", fileNameIpv4, "and", fileNameIpv6)
		}
	}

}

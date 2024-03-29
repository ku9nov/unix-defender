package iptables

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"unix-defender/utils"
)

func clearRules(IPv string) error {
	if err := ipTablesManage(IPv, "-F"); err != nil {
		log.Fatal(err)
	}
	log.Println(IPv, "rules removed.")
	return nil
}

func saveRules(saveCommand string, fileName *string) error {
	err := os.Remove(filepath.Join(utils.MainDir, filepath.Base(*fileName)))
	if err != nil {
		_ = err
		//Do nothing.
	}
	file, err := os.OpenFile(filepath.Join(utils.MainDir, filepath.Base(*fileName)), os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal("File with saved rules not exists or cannot be created", err)
	}
	defer file.Close()
	cmd, err := exec.Command(saveCommand).CombinedOutput()
	if err != nil {
		log.Fatal("Can't save iptables rules to file.", string(cmd), err)
	}
	writer := bufio.NewWriter(file)
	fmt.Fprint(writer, string(cmd))
	writer.Flush()

	return nil
}

func RestoreRules(restoreCommand string, fileName *string) error {
	file := filepath.Join(utils.MainDir, filepath.Base(*fileName))
	cmd, err := exec.Command(restoreCommand, file).CombinedOutput()
	if err != nil {
		log.Fatal("Can't restore rules: ", string(cmd), err)
	}
	return nil

}

func process(rules *utils.ConfigJson) error {
	port := fmt.Sprint(rules.Port)

	if rules.Chain == "INPUT" {
		if rules.Allow[0] != "0/0" {

			if err := ipTablesManage(rules.Version, "-A", rules.Chain, "-i", rules.Iface, "-p", rules.Protocol, "--dport", port, "-j", "DROP"); err != nil {
				return err
			}
		}

		for _, a := range rules.Allow {
			if err := ipTablesManage(rules.Version, "-I", rules.Chain, "-i", rules.Iface, "-s", a, "-p", rules.Protocol, "--dport", port, "-j", "ACCEPT"); err != nil {
				return err
			}
		}
	} else {
		if rules.Allow[0] != "0/0" {
			if err := ipTablesManage(rules.Version, "-A", rules.Chain, "-o", rules.Iface, "-p", rules.Protocol, "--dport", port, "-j", "DROP"); err != nil {
				return err
			}
		}

		for _, a := range rules.Allow {
			if err := ipTablesManage(rules.Version, "-I", rules.Chain, "-o", rules.Iface, "-s", a, "-p", rules.Protocol, "--dport", port, "-j", "ACCEPT"); err != nil {
				return err
			}
		}
	}

	// if err := iptables("-A INPUT -j DROP"); err != nil { /**** return in future, drop all other input trafic ****\
	// 	return err
	// }
	return nil
}

func ipTablesManage(args ...string) error {
	arg := args[1:]
	if args[0] == "IPv4" {
		cmd := exec.Command(utils.IptablesCommand, arg...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			if bytes.Contains(out, []byte("This doesn't exist in IPv4Tables")) {
				return nil
			}
			return err
		}
	} else {
		cmd := exec.Command(utils.Ip6tablesCommand, arg...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			if bytes.Contains(out, []byte("This doesn't exist in IPv6Tables")) {
				return nil
			}
			return err
		}
	}

	return nil
}

func IpTables() {
	configEnv, err := utils.LoadConfigEnv(utils.EnvFile)
	if err != nil {
		log.Fatal("Cannot load environment config:", err)
	}
	path := filepath.Join(utils.MainDir, filepath.Base(configEnv.RulesFile))
	configs, err := utils.LoadConfigJson(path)
	if err != nil {
		log.Println("Cannot load Json config:", err)
		return
	}
	clearRules("IPv4")
	clearRules("IPv6")
	for _, rules := range configs {
		if err := process(rules); err != nil {
			log.Fatal(err)
		}
	}
	log.Println("Successfully added new iptables rules.")
	saveRules(utils.SaveIpv4Command, &configEnv.RulesBackupV4)
	saveRules(utils.SaveIpv6Command, &configEnv.RulesBackupV6)
	utils.SendMessageToSlack(utils.ReconfigureMessage, utils.GreenColor)
}

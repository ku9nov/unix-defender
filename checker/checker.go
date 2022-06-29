package checker

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
)

func assertEq(first []string, second []string) bool {
	return reflect.DeepEqual(first, second)
}

func scanFile(fileRead *os.File) []string {
	var file []string
	scan := bufio.NewScanner(fileRead)
	for scan.Scan() {
		currentLineText := scan.Text()
		file = append(file, currentLineText)
	}
	result := file[5 : len(file)-2]
	return result
}

func readFile(fileForRead string, tmp bool) []string {
	if tmp {
		fileRead, err := os.Open(filepath.Join("/tmp", filepath.Base(fileForRead)))
		if err != nil {
			panic(err)
		}
		return scanFile(fileRead)

	} else {
		fileRead, err := os.Open(fileForRead)
		if err != nil {
			panic(err)
		}
		return scanFile(fileRead)
	}

}

func compareFiles(fileName string, fileNameMain *string) {
	if assertEq(readFile(fileName, true), readFile(*fileNameMain, false)) {
		fmt.Println("Rules files are equal")
	} else {
		fmt.Println("ALARM, IPTABLES IS CORRUPTED")
	}
}

func SaveRulesTmp(saveCommand string, fileName string, fileNameMain *string) error {
	file, err := os.OpenFile(filepath.Join("/tmp", filepath.Base(fileName)), os.O_WRONLY|os.O_CREATE, 0600)
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
	compareFiles(fileName, fileNameMain)
	return nil
}

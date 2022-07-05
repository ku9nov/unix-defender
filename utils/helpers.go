package utils

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func SigTerm(ipV4file string, ipV6file string) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		test := []string{ipV4file, ipV6file}
		for _, s := range test {
			e := os.Remove(filepath.Join("/tmp", filepath.Base(s)))
			if e != nil {
				log.Fatal(e)
			}
			fmt.Println("Succesfully removed:", s)
		}
		SendMessageToSlack(DisabledMessage, RedColor)
		os.Exit(1)
	}()
}

func RandomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

package scanner

import (
	"fmt"
	"log"
	"net"
	"os"
	"sort"

	"unix-defender/utils"
)

func LocalAddresses() {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Print(fmt.Errorf("localAddresses: %+v\n ", err.Error()))
		return
	}
	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Print(fmt.Errorf("localAddresses: %+v\n ", err.Error()))
			continue
		}
		for _, a := range addrs {
			switch v := a.(type) {
			case *net.IPAddr:
				fmt.Printf("%v : %s (%s)\n", i.Name, v, v.IP.DefaultMask())

			case *net.IPNet:
				fmt.Printf("%v : %s [%v/%v]\n", i.Name, v, v.IP, v.Mask)
			}

		}
	}
	os.Exit(1)
}

func worker(ports chan int, protocol string, host string, results chan int) {
	for port := range ports {
		address := fmt.Sprintf("%v:%d", host, port)
		connection, err := net.Dial(protocol, address)
		if err != nil {
			results <- 0
			continue
		}
		connection.Close()
		results <- port
	}
}
func ScanPorts() {
	config, err := utils.LoadConfigEnv(utils.EnvFile)
	if err != nil {
		log.Fatal("Cannot load environment config:", err)
	}
	// this channel will receive ports to be scanned
	ports := make(chan int, 100)
	// this channel will receive results of scanning
	results := make(chan int)
	// create a slice to store the results so that they can be sorted later.
	var openports []int

	// create a pool of workers
	for i := 0; i < cap(ports); i++ {
		go worker(ports, config.Protocol, config.Host, results)
	}

	// send ports to be scanned
	go func() {
		for i := 1; i <= config.PortsAmount; i++ {
			ports <- i
		}
	}()

	for i := 0; i < config.PortsAmount; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	// After all the work has been completed, close the channels
	close(ports)
	close(results)
	// sort open port numbers
	sort.Ints(openports)
	for _, port := range openports {
		fmt.Printf("%d %v Open \n", port, config.Protocol)
	}
	// os.Exit(1)
}

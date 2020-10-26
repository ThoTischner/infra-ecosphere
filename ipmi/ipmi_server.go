package ipmi

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

import (
	"../utils"
	"../bmc"
)

var running bool = false

func DeserializeAndExecute(buf io.Reader, addr *net.UDPAddr, server *net.UDPConn) {
	RMCPDeserializeAndExecute(buf, addr, server)
}

func IPMIServerHandler(BMCIP string, BMCPORT string) {
	addr := fmt.Sprintf("%s:%s", BMCIP, BMCPORT)
	serverAddr, err := net.ResolveUDPAddr("udp", addr)
	utils.CheckError(err)

	server, err := net.ListenUDP("udp", serverAddr)
	utils.CheckError(err)
	defer server.Close()

	buf := make([]byte, 1024)
	for running {
		_, addr, _ := server.ReadFromUDP(buf)
		log.Println("Receive a UDP packet from ", addr.IP.String(), ":", addr.Port)

		bytebuf := bytes.NewBuffer(buf)
		DeserializeAndExecute(bytebuf, addr, server)
	}
}

func IPMIServerServiceRun() {
	signalChan := make(chan os.Signal, 1)
	exitChan := make(chan bool, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGINT)
	go func() {
		<- signalChan
		log.Println("Capture Interrupt from System, terminate this server.")
		running = false
		exitChan <- true
	}()

	running = true
	config := utils.LoadConfig("infra-ecosphere.cfg")
	if config.BmcNet == "true" {
		for ip, _ := range bmc.BMCs {
			port := bmc.BMCs[ip].Port
			go func(ip string, port string) {
				log.Println("Start BMC Listener for BMC foo ", ip)
				IPMIServerHandler(ip, port)
				log.Println("BMC Listener ", ip, " is terminated.")
			}(ip, string(port))
		}
	} else
	{
		log.Println("Dont start BMC loopback ips, cause BmcNet config set to false")
	}


	<- exitChan
	log.Println("Wait for Listener terminating...")
	time.Sleep(3 * time.Second)
}

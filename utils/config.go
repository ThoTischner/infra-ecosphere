package utils

import (
	"../bmc"
	"../vm"
	"../web"
	"encoding/json"
	"log"
	"net"
	"os"
)

type ConfigNode struct {
	BMCIP string
    BMCPORT string
	VMName string
}

type ConfigBMCUser struct {
	Username string
	Password string
}

type Configuration struct {
	Nodes		[]ConfigNode
	BMCUsers	[]ConfigBMCUser
	WebAPIPort	int
	BmcNet		string
}

func LoadConfig(configFile string) Configuration {
	file, opError := os.Open(configFile)
	if opError != nil {
		log.Println("Config: Failed to open config file ", configFile, ", ignore...")
		return Configuration{
			WebAPIPort: 9090,
			BmcNet: "false",
		}
	}

	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Fatalln("Config: Error: ", err)
	}

	// initialize BMCs and Instances
	for _, node := range configuration.Nodes {
		fakeNode := false
		if len(node.VMName) == 0 {
			fakeNode = true
		}
		instance := vm.AddInstnace(node.VMName, fakeNode)
		bmc.AddBMC(net.ParseIP(node.BMCIP), node.BMCPORT, instance)
	}

	for _, user := range configuration.BMCUsers {
		log.Printf("Config: Add BMC User %s\n", user.Username)
		bmc.AddBMCUser(user.Username, user.Password)
	}

	if configuration.WebAPIPort <= 1024 || configuration.WebAPIPort > 65535 {
		log.Fatalln("Web API Port value should be larger than 1024 and less than 65536.")
	} else {
		web.ListenPort = configuration.WebAPIPort
	}

	return configuration
}
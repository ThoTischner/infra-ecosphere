package main

import (
	"./ipmi"
	"./utils"
	"./web"
)

func main() {
	utils.LoadConfig("infra-ecosphere.cfg")
	go ipmi.IPMIServerServiceRun()
	web.WebAPIServiceRun()
}

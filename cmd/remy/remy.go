package main

import (
	"github.com/klauern/remy/cli"
	"github.com/klauern/remy/wls"
)

// Config is the base configuration used for all REST requests to the AdminServer
var Config wls.AdminServer

func main() {
	cli.Run(&Config)
}

package main

import (
	cfg "github.com/klauern/remy"
	"github.com/klauern/remy/cmd"
)

// Config is the base configuration used for all REST requests to the AdminServer
var Config cfg.AdminServer

func main() {
	cmd.Run(&Config)
}

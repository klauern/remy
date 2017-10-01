package main

import (
    "github.com/klauern/remy/cmd"
    cfg "github.com/klauern/remy"
)

// Config is the base configuration used for all REST requests to the AdminServer
var Config cfg.AdminServer

func main() {
    cmd.Run(&Config)
}

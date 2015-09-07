package main

import (
	//	"github.com/klauern/wlsrest/wls"
	"github.com/spf13/cobra"
	//	"fmt"
	"github.com/klauern/wlsrest/cli"
)

func main() {
	var WlsRestCmd = &cobra.Command{
		Use:   "wlsrest",
		Short: "Query a WebLogic Server's resources",
		Long:  "Query your underlying server's resources, including Datasources, Applications, Clusters, and Servers by using the WebLogic RESTful Management Extensions API",
		Run:   cli.WlsRestCmd,
	}

	var serversCmd = &cobra.Command{
		Use:   "servers [Server to query, blank for ALL]",
		Short: "Display Server information",
		Long:  "Show details on all servers under an AdminServer, or specify a specific one",
		Run:   cli.Servers,
	}

	var clustersCmd = &cobra.Command{
		Use:   "clusters [cluster to query, blank for ALL]",
		Short: "Query clusters under AdminServer",
		Long:  "Query the AdminServer for specific clusters, or leave blank for all clusters that this server owns",
		Run:   cli.Clusters,
	}
	//	env := wls.Environment{Password:"pass", Username:"user", ServerUrl:"http://localhost:8080"}
	//	fmt.Print(env)

	WlsRestCmd.AddCommand(serversCmd)
	WlsRestCmd.AddCommand(clustersCmd)

	WlsRestCmd.Execute()
}

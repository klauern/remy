package main

import (
	//	"github.com/klauern/wlsrest/wls"
	"github.com/spf13/cobra"
	//	"fmt"
	"github.com/klauern/remy/cli"
	"github.com/klauern/remy/wls"
	"github.com/spf13/viper"
)

var Config wls.AdminServer

func main() {

	// Base command for the application.
	var WlsRestCmd = &cobra.Command{
		Use:   "remy",
		Short: "Query a WebLogic Server's Rest Management Extention-enabled resources",
		Long:  "Query your underlying server's resources, including Datasources, Applications, Clusters, and Servers by using the WebLogic RESTful Management Extensions API",
	}

	// Request the Servers resource, optionally passing a specific [servername] instance to get that particular Server.
	var serversCmd = &cobra.Command{
		Use:   "servers [Server to query, blank for ALL]",
		Short: "Display Server information",
		Long:  "Show details on all servers under an AdminServer, or specify a specific one",
		Run:   cli.Servers,
	}

	// Request the Clusters resource, optionally passing a specific [clustername] to get a specific Cluster.
	var clustersCmd = &cobra.Command{
		Use:   "clusters [cluster to query, blank for ALL]",
		Short: "Query clusters under AdminServer",
		Long:  "Query the AdminServer for specific clusters, or leave blank for all clusters that this server owns",
		Run:   cli.Clusters,
	}

	// Datasource command, requesting all datasrouces.  Pass a secondary [datasourcename] to get a specific datasource.
	var datasourcesCmd = &cobra.Command{
		Use:   "datasources [datasources to query, blank for ALL]",
		Short: "Query datasources under AdminServer",
		Long:  "Query the AdminServer for specific datasources, or leave blank for all datasources that this server owns",
		Run:   cli.DataSources,
	}

	// Application list command.  Pass an optional [applicationname] to get a specific application instance details.
	var applicationsCmd = &cobra.Command{
		Use:   "applications [application to query, blank for ALL]",
		Short: "Query applications deployed under AdminServer",
		Long:  "Query the AdminServer for specific applications, or leave blank for all applications that this server knows about",
		Run:   cli.Applications,
	}

	// Generate a configuration setting file in your ~/ home or local directory.
	// When determined to be in the ~/home, it will be a ~/.wlsrest.cfg file.
	// When a local file, it will be a wlsrest.cfg file instead.
	var configureCmd = &cobra.Command{
		Use:   "config",
		Short: "Configure the credentials and server to default REST connections to",
		Long:  "Configure what Username, Password, and Admin Server:Port you want to send REST requests to when submitting calls on any of the other commands",
		Run:   cli.Configure,
	}

	// Add option to pass --full-format for all responses.  Single server, application, etc., requests will always return
	// full responses, but group-related queries will return shortened versions
	WlsRestCmd.PersistentFlags().BoolVarP(&cli.FullFormat, cli.FullFormatFlag, "f", false, "Return full format from REST server")

	// Allow specific AdminServer URL to be passed in to override local config files
	WlsRestCmd.PersistentFlags().StringVarP(&Config.AdminUrl, cli.AdminURLFlag, "s", "http://localhost:7001", "Url for the Admin Server")

	// Allow the Username property to be overridden locally on the command-line
	WlsRestCmd.PersistentFlags().StringVarP(&Config.Username, cli.UsernameFlag, "u", "weblogic", "Username with privileges to access AdminServer")

	// Allow the Password property to be overridden on the command-line
	WlsRestCmd.PersistentFlags().StringVarP(&Config.Password, cli.PasswordFlag, "p", "welcome1", "Password for the user")

	viper.BindPFlags(WlsRestCmd.PersistentFlags())

	WlsRestCmd.AddCommand(applicationsCmd, configureCmd, clustersCmd, datasourcesCmd, serversCmd)
	WlsRestCmd.Execute()
}

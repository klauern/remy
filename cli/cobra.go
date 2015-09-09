package cli

import (
	"fmt"
	"github.com/klauern/remy/wls"
	"github.com/spf13/cobra"
//	"os"
//	"os/user"
	"github.com/spf13/viper"
"os/user"
	"os"
)

const (
	CONFIG_FILE = "wlsrest"
)

// Determined whether to request fully-formatted responses from the REST endpoint.  For single-instance requests, this is always
// a full format, but for groups (servers, applications, clusters, etc.) the server defaults to a short-form response.
var FullFormat bool

func WlsRestCmd(cmd *cobra.Command, args []string) {
	fmt.Printf("Args passed to WlsRestCmd: %v\n", args)

}

func Servers(cmd *cobra.Command, args []string) {
	fmt.Printf("Args passed to Servers: %v\n", args)
	env, err := findConfiguration()
	if err != nil {
		panic(fmt.Sprintf("No configuration found.  Please call 'help config' to find out how to set this"))
	}
	if len(args) > 1 {
		servers, err := env.Servers(FullFormat)
		if err != nil {
			panic(fmt.Sprintf("Unable to get Servers: %v", err))
		}
		fmt.Printf("Servers:\n%+v", servers)
	}
}

func Clusters(cmd *cobra.Command, args []string) {
	fmt.Printf("Args passed to Clusters: %v\n", args)
}

func DataSources(cmd *cobra.Command, args []string) {
	fmt.Printf("Args passed to DataSources: %v\n", args)
}

func Applications(cmd *cobra.Command, args []string) {
	fmt.Printf("Args passed to Applications: %v\n", args)
}

// Generate a configuration file to store default credentials to use when making REST queries to an AdminServer
func Configure(cmd *cobra.Command, args []string) {
	cfg, err := findConfiguration()
	if err != nil {
		fmt.Printf("Error found: %v", err)
	}
	fmt.Printf("Current Working Directory: %v", cfg.ServerUrl)
}

// Finds a configuration setting for your login.  Looks for the following configuration file, processed in the following
// order:
//   - command-line flags --username, --password and --server <host:port>
//   - WLSREST_CONFIG (environment variable)
//   - wlsrest.cfg (in the current directory)
//   - .wlsrest.cfg (in the $HOME directory)
//
// This is borrowed lovingly from Ansible's similar setup for it's configuration (http://docs.ansible.com/ansible/intro_configuration.html)
func findConfiguration() (*wls.WlsAdminServer, error) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("WLS_")
	viper.BindEnv("USERNAME")
	viper.BindEnv("PASSWORD")
	viper.BindEnv("HOSTPORT")
	viper.SetConfigType("toml")
	viper.SetConfigName(CONFIG_FILE)
	u, err := user.Current()
	if err != nil {
		return nil, err
	}
	viper.AddConfigPath(u.HomeDir)
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	viper.AddConfigPath(cwd)
	return &wls.WlsAdminServer{"", "", ""}, nil
}

package cli

import (
	"fmt"
	"github.com/klauern/remy/wls"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/user"
)

const (
	// ConfigFile is the base file prefix for looking for configuration files.  wlsrest.toml, .wlsrest.toml are all valid filenames
	ConfigFile = "wlsrest"

	// AdminURLFlag is the flag for specifying/overriding the Administration Server url (http://localhost:7001)
	AdminURLFlag = "adminurl"

	// PasswordFlag is the flag for specifying/overriding the Password to log in to the AdminServer
	PasswordFlag = "password"

	// UsernameFlag is the flag for specifying/overriding the Username to log in to AdminServer with
	UsernameFlag = "username"

	// FullFormatFlag is the flag to override whether to request the fully-formatted dataset for a resource
	FullFormatFlag = "full-format"
)

// FullFormat determines whether to request fully-formatted responses from the REST endpoint.  For single-instance requests, this is always
// a full format, but for groups (servers, applications, clusters, etc.) the server defaults to a short-form response.
var FullFormat bool

// Servers takes a Viper Command and it's argument list, and calls the underlying wls.Servers service to retrieve server
// information.
func Servers(cmd *cobra.Command, args []string) {
	//	fmt.Printf("Args passed to Servers: %v\n", args)
	env, err := findConfiguration()
	if err != nil {
		panic(fmt.Sprintf("No configuration found.  Please call 'help config' to find out how to set this"))
	}
	if len(args) > 2 {
		panic(fmt.Sprintf("Too many arguments.  enter 'help servers' command to find out how to call this"))
	}
	if len(args) == 1 {
		server, err := env.Server(args[0])
		if err != nil {
			panic(fmt.Sprintf("Unable to get Servers: %v", err))
		}
		fmt.Printf("Server %v: %v", args[0], server)
	}
	if len(args) == 0 {
		fmt.Printf("Full Format? %+v", FullFormat)
		fmt.Printf("Environment: %+v", env)
		servers, err := env.Servers(FullFormat)
		if err != nil {
			panic(fmt.Sprintf("Unable to get Servers: %v", err))
		}
		fmt.Printf("Servers:\n%+v", servers)
	}
}

// Clusters takes a viper.Command object and arguments to call the AdminServer to retrieve Cluster information
func Clusters(cmd *cobra.Command, args []string) {
	fmt.Printf("Args passed to Clusters: %v\n", args)
	env, err := findConfiguration()
	if err != nil {
		panic(fmt.Sprintf("No configuration found.  Please call 'help config' to find out how to set this"))
	}
	if len(args) > 2 {
		panic(fmt.Sprintf("Too many arguments.  enter 'help clusters' command to find out how to call this"))
	}
	if len(args) == 1 {
		server, err := env.Cluster(args[0])
		if err != nil {
			panic(fmt.Sprintf("Unable to get Clusters: %v", err))
		}
		fmt.Printf("Cluster %v: %v", args[0], server)
	}
	if len(args) == 0 {
		fmt.Printf("Environment: %+v", env)
		clusters, err := env.Clusters(FullFormat)
		if err != nil {
			panic(fmt.Sprintf("Unable to get Clusters: %v", err))
		}
		fmt.Printf("Clusters:\n%+v", clusters)
	}
}

// DataSources is a command function to call out the wls.DataSources resource running on a remote AdminServer.
func DataSources(cmd *cobra.Command, args []string) {
	fmt.Printf("Args passed to DataSources: %v\n", args)
	env, err := findConfiguration()
	if err != nil {
		panic(fmt.Sprintf("No configuration found.  Please call 'help config' to find out how to set this"))
	}
	if len(args) > 2 {
		panic(fmt.Sprintf("Too many arguments.  enter 'help datasources' command to find out how to call this"))
	}
	if len(args) == 1 {
		datasource, err := env.DataSource(args[0])
		if err != nil {
			panic(fmt.Sprintf("Unable to get Datasource: %v", err))
		}
		fmt.Printf("Datasource %v: %v", args[0], datasource)
	}
	if len(args) == 0 {
		fmt.Printf("Environment: %+v", env)
		datasources, err := env.DataSources(FullFormat)
		if err != nil {
			panic(fmt.Sprintf("Unable to get Datasources: %v", err))
		}
		fmt.Printf("Datasources:\n%+v", datasources)
	}
}

// Applications is a Cobra command function to call out to the wls.Applications resource on a remote AdminServer.
func Applications(cmd *cobra.Command, args []string) {
	fmt.Printf("Args passed to Applications: %v\n", args)
	env, err := findConfiguration()
	if err != nil {
		panic(fmt.Sprintf("No configuration found.  Please call 'help config' to find out how to set this"))
	}
	if len(args) > 2 {
		panic(fmt.Sprintf("Too many arguments.  enter 'help applications' command to find out how to call this"))
	}
	if len(args) == 1 {
		application, err := env.Application(args[0])
		if err != nil {
			panic(fmt.Sprintf("Unable to get Application: %v", err))
		}
		fmt.Printf("Application %v: %v", args[0], application)
	}
	if len(args) == 0 {
		fmt.Printf("Environment: %+v", env)
		applications, err := env.Applications(FullFormat)
		if err != nil {
			panic(fmt.Sprintf("Unable to get Applications: %v", err))
		}
		fmt.Printf("Applications:\n%+v", applications)
	}
}

// Configure will generate a configuration file to store default credentials to use when making REST queries to an AdminServer
func Configure(cmd *cobra.Command, args []string) {
	cfg, err := findConfiguration()
	if err != nil {
		fmt.Printf("Error found: %v", err)
	}
	fmt.Printf("Current Working Directory: %v", cfg.AdminURL)
}

// findConfiguration finds a configuration setting for your login.  Looks for the following configuration file, processed in the following
// order:
//   - command-line flags --username, --password and --server <host:port>
//   - WLSREST_CONFIG (environment variable)
//   - wlsrest.cfg (in the current directory)
//   - .wlsrest.cfg (in the $HOME directory)
//
// This is borrowed lovingly from Ansible's similar setup for it's configuration (http://docs.ansible.com/ansible/intro_configuration.html)
func findConfiguration() (*wls.AdminServer, error) {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("WLS")
	viper.BindEnv(UsernameFlag)
	viper.BindEnv(PasswordFlag)
	viper.BindEnv(AdminURLFlag)
	viper.SetConfigType("toml")
	viper.SetConfigName(ConfigFile)
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

	server := &wls.AdminServer{}
	if viper.IsSet(UsernameFlag) {
		server.Username = viper.GetString(UsernameFlag)
	}
	if viper.IsSet(PasswordFlag) {
		server.Password = viper.GetString(PasswordFlag)
	}
	if viper.IsSet(AdminURLFlag) {
		server.AdminURL = viper.GetString(AdminURLFlag)
	}

	return server, nil
}

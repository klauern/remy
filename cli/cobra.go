package cli

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/klauern/remy/wls"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strings"
)

const (
	// ConfigFile is the base file prefix for looking for configuration files.  wlsrest.toml, .wlsrest.toml are all valid filenames
	ConfigFile = "wlsrest"

	// ConfigFileSuffix is used to append the suffix to the config file.  We'll be using TOML format.
	ConfigFileSuffix = ".toml"

	// AdminURLFlag is the flag for specifying/overriding the Administration Server url (http://localhost:7001)
	AdminURLFlag = "adminurl"

	// PasswordFlag is the flag for specifying/overriding the Password to log in to the AdminServer
	PasswordFlag = "password"

	// UsernameFlag is the flag for specifying/overriding the Username to log in to AdminServer with
	UsernameFlag = "username"

	// FullFormatFlag is the flag to override whether to request the fully-formatted dataset for a resource
	FullFormatFlag = "full-format"

	// EnvironmentSetFlag is used in the 'config' command to determine whether to set user, pass, serverUrl in the
	// environment variables.  These environment variables are prefixed with WLS_*.
	EnvironmentSetFlag = "environment"

	// LocalSetFlag is the flag used in the 'config' command for setting whether to generate/update the local directory's ./wlsrest.toml config
	// file.
	LocalSetFlag = "local"

	// HomeSetFlag is the flag used in the 'config' command to set whether to generate/update the ~/.wlsrest.toml configuration file
	HomeSetFlag = "home"
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

// Configure generates or updates a configuration file to store default credentials to use when making REST queries to an AdminServer
func Configure(cmd *cobra.Command, args []string) {
	cfg, err := findConfiguration()
	if err != nil {
		panic(fmt.Errorf("Not able to find configuration: %s \n", err))
	}

	if viper.GetBool(EnvironmentSetFlag) {
		fmt.Printf("Using the environment variables to set the %v, %v, and %v environment variables\n", "WLS_USERNAME", "WLS_PASSWORD", "WLS_ADMINURL")
		setEnvironmentVariables(cfg)
	}
	var buf bytes.Buffer
	enc := toml.NewEncoder(&buf)
	err = enc.Encode(cfg)
	if err != nil {
		panic(fmt.Errorf("Unable to encode wlsrest configuration: %s \n", err))
	}

	if viper.GetBool(LocalSetFlag) {
		fmt.Printf("Using the Local directory to set the ./wlsrest.toml file\n")
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile(path.Join(cwd, ConfigFile+".toml"), buf.Bytes(), 0644)
		if err != nil {
			panic(err)
		}
	}

	if viper.GetBool(HomeSetFlag) {
		fmt.Printf("Using the $HOME directory to generate the ~/.wlsrest.toml file\n")
		u, err := user.Current()
		if err != nil {
			panic(err)
		}
		home := u.HomeDir
		err = ioutil.WriteFile(path.Join(home, "."+ConfigFile+".toml"), buf.Bytes(), 0644)
		if err != nil {
			panic(err)
		}
	}
}

func setEnvironmentVariables(env *wls.AdminServer) {
	prefixEnv := "WLS_"

	if uerr := os.Setenv(strings.ToUpper(prefixEnv+UsernameFlag), env.Username); uerr != nil {
		panic(uerr)
	}
	if perr := os.Setenv(strings.ToUpper(prefixEnv+PasswordFlag), env.Password); perr != nil {
		panic(perr)
	}
	if aerr := os.Setenv(strings.ToUpper(prefixEnv+AdminURLFlag), env.AdminURL); aerr != nil {
		panic(aerr)
	}
}

// findConfiguration finds a configuration setting for your login.  Looks for the following configuration file, processed in the following
// order of precedence:
//   - command-line flags --username, --password and --server <host:port>
//   - WLS_USERNAME, WLS_PASSWORD, WLS_ADMINURL (environment variables)
//   - wlsrest.toml (in the current directory)
//   - .wlsrest.toml (in the $HOME directory)
//
// This is borrowed lovingly from Ansible's similar setup for it's configuration (http://docs.ansible.com/ansible/intro_configuration.html)
func findConfiguration() (*wls.AdminServer, error) {
	viper.SetEnvPrefix("WLS")
	viper.BindEnv(UsernameFlag)
	viper.BindEnv(AdminURLFlag)
	viper.BindEnv(PasswordFlag)

	viper.SetConfigType("toml")
	cwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("errof getting current working directory: %s\n", err))
		return nil, err
	}
	viper.SetConfigName(ConfigFile)
	viper.AddConfigPath(cwd)
	viper.ReadInConfig()

	viper.SetConfigName("." + ConfigFile)
	u, err := user.Current()
	if err != nil {
		return nil, err
	}
	fmt.Printf("User Home: %v\n", u.HomeDir)
	viper.AddConfigPath(u.HomeDir)
	viper.ReadInConfig()

	server := &wls.AdminServer{}
	server.Username = viper.GetString(UsernameFlag)
	server.Password = viper.GetString(PasswordFlag)
	server.AdminURL = viper.GetString(AdminURLFlag)

	viper.AutomaticEnv()

	return server, nil
}

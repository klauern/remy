package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"strings"

    wls "github.com/klauern/remy"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/pkg/errors"
)

const (
	remyVersion string = "0.2.1"

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

	// RemyKey is the key used to get or set the encryption key used in encrypting a password
	RemyKey = "remykey"

	// DefaultRemyKeyString is the default 32-byte string used to encrypt plain-text passwords in AES.
	DefaultRemyKeyString = "This is the default key.  Please"

	// EncryptedPrefix is what is prepended to any encrypted password.  If a password doesn't have this, it is treated
	// as plain text instead.
	EncryptedPrefix = "{AES}"

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

// FlagLocalConfig determines whether to generate/update a the config file in the local ./ directory or not
var FlagLocalConfig bool

// FlagHomeConfig determines whether to generate/update the $HOME ~/ folder's .wlstrest.cfg file or not
var FlagHomeConfig bool

// Servers takes a Viper Command and it's argument list, and calls the underlying wls.Servers service to retrieve server
// information.
func Servers(cmd *cobra.Command, args []string) {
	env := findConfiguration()
	if len(args) > 2 {
		panic(fmt.Sprintf("Too many arguments.  enter 'help servers' command to find out how to call this"))
	}
	if len(args) == 1 {
		fmt.Printf("Finding Server information for %v\n", args[0])
		server, err := env.Server(args[0])
		if err != nil {
			panic(fmt.Sprintf("Unable to get Servers: %v", err))
		}
		fmt.Printf("Server %v:\n%#v", args[0], server)
	}
	if len(args) == 0 {
		fmt.Printf("Finding all Servers\nUsing Full Format? %v\n", FullFormat)
		servers, err := env.Servers(FullFormat)
		if err != nil {
			panic(fmt.Sprintf("Unable to get Servers: %v", err))
		}
		for i := range servers {
			fmt.Printf("%#v\n", &servers[i])
		}
	}
}

// Clusters takes a viper.Command object and arguments to call the AdminServer to retrieve Cluster information
func Clusters(cmd *cobra.Command, args []string) {
	env := findConfiguration()
	if len(args) > 2 {
		panic(fmt.Sprintf("too many arguments.  enter 'help clusters' command to find out how to call this"))
	}
	if len(args) == 1 {
		fmt.Printf("Finding Cluster information for %v\n", args[0])
		cluster, err := env.Cluster(args[0])
		if err != nil {
			panic(fmt.Sprintf("unable to get Clusters: %v", err))
		}
		fmt.Printf("%#v\n", cluster)
	}
	if len(args) == 0 {
		fmt.Printf("Finding All Clusters\nUsing Full Format? %v\n", FullFormat)
		clusters, err := env.Clusters(FullFormat)
		if err != nil {
			panic(fmt.Sprintf("unable to get Clusters: %v", err))
		}
		for i := range clusters {
			fmt.Printf("%#v\n", &clusters[i])
		}
	}
}

// DataSources is a command function to call out the wls.DataSources resource running on a remote AdminServer.
func DataSources(cmd *cobra.Command, args []string) {
	env := findConfiguration()
	if len(args) > 2 {
		panic(fmt.Sprintf("Too many arguments.  enter 'help datasources' command to find out how to call this"))
	}
	if len(args) == 1 {
		fmt.Printf("Finding DataSource information for %v\n", args[0])
		datasource, err := env.DataSource(args[0])
		if err != nil {
			panic(fmt.Sprintf("Unable to get Datasource: %v", err))
		}
		fmt.Printf("Datasource %v: %v", args[0], datasource)
	}
	if len(args) == 0 {
		fmt.Printf("Finding all DataSources\nUsing Full Format? %v\n", FullFormat)
		datasources, err := env.DataSources(FullFormat)
		if err != nil {
			panic(fmt.Sprintf("Unable to get Datasources: %v\n", err))
		}
		fmt.Printf("Datasources:\n%+v", datasources)
	}
}

// Applications is a Cobra command function to call out to the wls.Applications resource on a remote AdminServer.
func Applications(cmd *cobra.Command, args []string) {
	env := findConfiguration()
	if len(args) > 2 {
		panic(fmt.Sprintf("Too many arguments.  enter 'help applications' command to find out how to call this"))
	}
	if len(args) == 1 {
		fmt.Printf("Finding application information for %v\n", args[0])
		application, err := env.Application(args[0])
		if err != nil {
			panic(fmt.Sprintf("Unable to get Application: %v", err))
		}
		fmt.Printf("%#v\n", application)
	}
	if len(args) == 0 {
		fmt.Printf("Finding All Applications\nUsing Full Format? %v\n", FullFormat)
		applications, err := env.Applications(FullFormat)
		if err != nil {
			panic(fmt.Sprintf("Unable to get Applications: %v\n", err))
		}
		for i := range applications {
			fmt.Printf("%#v", &applications[i])
		}
	}
}

// Configure generates or updates a configuration file to store default credentials to use when making REST queries to an AdminServer
func Configure(cmd *cobra.Command, args []string) {
	cfg := findConfiguration()

	// Encrypt the password before setting it in the config
	if !strings.Contains(cfg.Password, EncryptedPrefix) {
		encryptedPass := encrypt([]byte(viper.GetString(RemyKey)), cfg.Password)
		viper.Set(PasswordFlag, encryptedPass)
		cfg.Password = EncryptedPrefix + encryptedPass
	}

	if viper.GetBool(LocalSetFlag) {
		fmt.Printf("Using the Local directory to generate the wlsrest.toml file\n")
		cwd, err := os.Getwd()
		if err != nil {
			panic(fmt.Errorf("Unable to get current working directory: %v", err))
		}

		err = ioutil.WriteFile(path.Join(cwd, ConfigFile+".toml"), cfg.EncodeConfigFile().Bytes(), 0644)
		if err != nil {
			panic(fmt.Errorf("Unable to write file to %v: %v", cwd, err))
		}
	}

	if viper.GetBool(HomeSetFlag) {
		fmt.Printf("Using the $HOME directory to generate the .wlsrest.toml file\n")
		u, err := user.Current()
		if err != nil {
			panic(fmt.Errorf("Unable to get current User: %v", err))
		}
		home := u.HomeDir
		err = ioutil.WriteFile(path.Join(home, "."+ConfigFile+".toml"), cfg.EncodeConfigFile().Bytes(), 0644)
		if err != nil {
			panic(fmt.Errorf("Unable to write configuration file to ~/: %v", err))
		}
	}
}

// findConfiguration finds and retrieves a configuration setting for your login.  It looks for the configurations in the following locations,
// processed in the following order of precedence (higher to lower precedence):
//   - command-line flags --username, --password and --server <host:port>
//   - WLS_USERNAME, WLS_PASSWORD, WLS_ADMINURL (environment variables)
//   - wlsrest.toml (in the current directory)
//   - .wlsrest.toml (in the $HOME directory)
//
// This is borrowed lovingly from Ansible's similar setup for it's configuration (http://docs.ansible.com/ansible/intro_configuration.html)
func findConfiguration() *wls.AdminServer {
	// We only load TOML files currently
	viper.SetConfigType("toml")

	viper.SetDefault(RemyKey, DefaultRemyKeyString)

	// Add the ./wlsrest.toml configuration file first, it will override the next file to find
	if cwd, err := os.Getwd(); err == nil {
		viper.SetConfigName(ConfigFile)
		viper.AddConfigPath(cwd)
		if err = viper.ReadInConfig(); err != nil {
			panic(errors.WithMessage(err, "unable to read in configuration"))
		}


	}
	// Add the ~/.wlsrest.toml config next.  It will fall-through to find this one if there's not one in the current dir
	viper.SetConfigName("." + ConfigFile)
	if u, err := user.Current(); err == nil {
		viper.AddConfigPath(u.HomeDir)
		if err = viper.ReadInConfig(); err != nil {
			panic(errors.WithMessage(err, "unable to read in configuration"))
		}
	}

	// Define environment variables to look for.
	viper.SetEnvPrefix("WLS")
	viper.AutomaticEnv()

	// Finally, load the configuration pieces from Viper
	server := &wls.AdminServer{}
	server.Username = viper.GetString(UsernameFlag)
	if strings.Contains(viper.GetString(PasswordFlag), "{AES}") {
		server.Password = decrypt([]byte(viper.GetString(RemyKey)), viper.GetString(PasswordFlag)[len(EncryptedPrefix):])
	} else {
		server.Password = viper.GetString(PasswordFlag)
	}
	server.AdminURL = viper.GetString(AdminURLFlag)

	//	fmt.Printf("%+v\n", server)
	return server
}

// encrypt string to base64 crypto using AES
func encrypt(key []byte, text string) string {
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}

// decrypt from base64 to decrypted string
func decrypt(key []byte, cryptoText string) string {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}

// Run runs our `remy` application.
func Run(cfg *wls.AdminServer) {
	// Base command for the application.
	var WlsRestCmd = &cobra.Command{
		Use:   "remy",
		Short: "Query a WebLogic Domain's REST Management Extention-enabled resources",
		Long:  "Query a WebLogic Domain's resources, including Datasources, Applications, Clusters, and Servers by using the WebLogic RESTful Management Extensions API",
	}

	// Request the Servers resource, optionally passing a specific [servername] instance to get that particular Server.
	var serversCmd = &cobra.Command{
		Use:   "servers [Server to query, blank for ALL]",
		Short: "Display Server information",
		Long:  "Show details on all servers under an AdminServer, or specify a specific one",
		Run:   Servers,
	}

	// Request the Clusters resource, optionally passing a specific [clustername] to get a specific Cluster.
	var clustersCmd = &cobra.Command{
		Use:   "clusters [cluster to query, blank for ALL]",
		Short: "Query clusters under AdminServer",
		Long:  "Query the AdminServer for specific clusters, or leave blank for all clusters that this server owns",
		Run:   Clusters,
	}

	// Datasource command, requesting all datasrouces.  Pass a secondary [datasourcename] to get a specific datasource.
	var datasourcesCmd = &cobra.Command{
		Use:   "datasources [datasources to query, blank for ALL]",
		Short: "Query datasources under AdminServer",
		Long:  "Query the AdminServer for specific datasources, or leave blank for all datasources that this server owns",
		Run:   DataSources,
	}

	// Application list command.  Pass an optional [applicationname] to get a specific application instance details.
	var applicationsCmd = &cobra.Command{
		Use:   "applications [application to query, blank for ALL]",
		Short: "Query applications deployed under AdminServer",
		Long:  "Query the AdminServer for specific applications, or leave blank for all applications that this server knows about",
		Run:   Applications,
	}

	// Generate a configuration setting file in your ~/ home or local directory.
	// When determined to be in the ~/home, it will be a ~/.wlsrest.toml file.
	// When a local file, it will be a wlsrest.toml file instead.
	var configureCmd = &cobra.Command{
		Use:   "config",
		Short: "Configure the credentials and server to default REST connections to",
		Long:  "Configure what Username, Password, and Admin Server:Port you want to send REST requests to when submitting calls on any of the other commands",
		Run:   Configure,
	}

	// Version command displays the version of the application.
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Show the version of this command",
		Long:  "Display the version of this command",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("remy version %v\n", remyVersion)
		},
	}

	// Add option to pass --full-format for all responses.  Single server, application, etc., requests will always return
	// full responses, but group-related queries will return shortened versions
	WlsRestCmd.PersistentFlags().BoolVarP(&FullFormat, FullFormatFlag, "f", false, "Return full format from REST server")

	// Allow specific AdminServer URL to be passed in to override local config files
	WlsRestCmd.PersistentFlags().StringVarP(&cfg.AdminURL, AdminURLFlag, "s", "http://localhost:7001", "Url for the Admin Server")

	// Allow the Username property to be overridden locally on the command-line
	WlsRestCmd.PersistentFlags().StringVarP(&cfg.Username, UsernameFlag, "u", "weblogic", "Username with privileges to access AdminServer")

	// Allow the Password property to be overridden on the command-line
	WlsRestCmd.PersistentFlags().StringVarP(&cfg.Password, PasswordFlag, "p", "welcome1", "Password for the user")

	configureCmd.Flags().BoolVar(&FlagHomeConfig, HomeSetFlag, false, "Generate/Update the ~/$HOME config file")
	configureCmd.Flags().BoolVar(&FlagLocalConfig, LocalSetFlag, false, "Generate/Update the local directory's config file")

	if err := viper.BindPFlags(WlsRestCmd.PersistentFlags()); err != nil {
		panic(errors.WithMessage(err, "cannot bind flag for "+WlsRestCmd.Name()))
	}
	if err := viper.BindPFlags(configureCmd.Flags()); err != nil {
		panic(errors.WithMessage(err, "cannot bind flag for "+configureCmd.Name()))
	}


	WlsRestCmd.AddCommand(applicationsCmd, configureCmd, clustersCmd, datasourcesCmd, serversCmd, versionCmd)
	if err := WlsRestCmd.Execute(); err != nil {
		panic(errors.WithMessage(err, "error executing "+WlsRestCmd.Name()))
	}
}

package cli

import (
	"fmt"
	"github.com/spf13/cobra"
)

func WlsRestCmd(cmd *cobra.Command, args []string) {
	fmt.Printf("Args passed to WlsRestCmd: %v\n", args)
}

func Servers(cmd *cobra.Command, args []string) {
	fmt.Printf("Args passed to Servers: %v\n", args)
}

func Clusters(cmd *cobra.Command, args []string) {
	fmt.Printf("Args passed to Clusters: %v\n", args)

}

// Configure default credentials to use when making REST queries to an AdminServer
func Configure(cmd *cobra.Command, args []string) {

}

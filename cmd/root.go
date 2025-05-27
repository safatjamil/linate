package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "linate",
	Short: "A linux associate. Get information, check backup files, network connections, and many other things.",
	Long:  `A linux associate. Get information about the os/memory/cpu. Check backup files, network connections
and do many other things. Run linate -h for more information.
Visit https://github.com/safatjamil/linate for more information.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(netCmd)
	rootCmd.AddCommand(backUpCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.CompletionOptions.DisableDefaultCmd = true  
}

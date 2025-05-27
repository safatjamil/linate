package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var version = "1.0.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Linate version.",
	Long:  `Linate version. Visit https://github.com/safatjamil/linate for more information.`,
	Run: show_version_info,
}


func show_version_info(cmd *cobra.Command, args []string) {
	fmt.Printf("%s%s%s\n", colors["yellow"], version, colors["reset"])
}
package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"runtime"
)
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Information about the os, memory, disk usage and other metrics",
	Long:  `Information about the os, memory, disk usage and other metrics.Run linate info --help for more options.`,
	// Run: os_info,
}

var osCmd = &cobra.Command{
	Use:   "os",
	Short: "Information about the os",
	Long:  `Information about the os`,
	Run: os_info,
}

func init() {
	//osCmd.PersistentFlags().StringP("type", "-t", "", "os | memory | disk")
	// infoCmd.PersistentFlags().StringP("dir", "-t", "", "os | memory | disk")
	rootCmd.AddCommand(infoCmd)
	infoCmd.AddCommand(osCmd)
}


type osInfo struct {
	OperatingSystem string
	Architecture string
	Distribution string
}

func os_info(cmd *cobra.Command, args [] string) {
	os := runtime.GOOS
	arch := runtime.GOARCH

	fmt.Print("OS: ", os, "\nArchitecture: ", arch, "\n")
}
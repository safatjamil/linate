package cmd

import (
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/v4/host"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Information about the os, memory, disk usage and other metrics",
	Long:  `Information about the os, memory, disk usage and other metrics. Run linate info --help for more options.`,
	// Run: os_info,
}

var osCmd = &cobra.Command{
	Use:   "os",
	Short: "Information about the os",
	Long:  `Information about the os`,
	Run:   os_info,
}

func init() {
	//osCmd.PersistentFlags().StringP("type", "-t", "", "os | memory | disk")
	// infoCmd.PersistentFlags().StringP("dir", "-t", "", "os | memory | disk")
	rootCmd.AddCommand(infoCmd)
	infoCmd.AddCommand(osCmd)
}

type OsInfo struct {
	Architecture  string
	Distribution  string
	Version       string
	TotalMemory   string
	KernelVersion string
}

type platform struct{}

func os_info(cmd *cobra.Command, args []string) {
	var osinfo OsInfo
	osinfo.Distribution, _, osinfo.Version, _ = host.PlatformInformation()
	osinfo.KernelVersion, _ = host.KernelVersion()
	osinfo.Architecture = runtime.GOARCH

	fmt.Print("OS: ", osinfo.Distribution, " ", osinfo.Version, "\n")
	fmt.Print("Architecture: ", osinfo.Architecture, "\n")
	fmt.Print("Kernel: ", osinfo.KernelVersion, "\n")

}

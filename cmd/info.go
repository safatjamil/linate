package cmd

import (
	"fmt"
	// "math"
	"runtime"
	// "strconv"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/disk"
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
	KernelVersion string
	CPUCount      int
	CPUModel      string
	TotalMemory   uint64
	TotalDisk     uint64
}

// type InfoStat struct {
// 	CPU        int32    `json:"cpu"`
// 	VendorID   string   `json:"vendorId"`
// 	Family     string   `json:"family"`
// 	Model      string   `json:"model"`
// 	Stepping   int32    `json:"stepping"`
// 	PhysicalID string   `json:"physicalId"`
// 	CoreID     string   `json:"coreId"`
// 	Cores      int32    `json:"cores"`
// 	ModelName  string   `json:"modelName"`
// 	Mhz        float64  `json:"mhz"`
// 	CacheSize  int32    `json:"cacheSize"`
// 	Flags      []string `json:"flags"`
// 	Microcode  string   `json:"microcode"`
// }

type platform struct{}

func os_info(cmd *cobra.Command, args []string) {
	var osinfo OsInfo
	
	osinfo.Architecture = runtime.GOARCH
	osinfo.Distribution, _, osinfo.Version, _ = host.PlatformInformation()
	osinfo.KernelVersion, _ = host.KernelVersion()

	cpuinfo, _ := cpu.Info()
	osinfo.CPUCount, _ = cpu.Counts(true)
	osinfo.CPUModel = cpuinfo[0].ModelName
    
	buff, _ := mem.VirtualMemory()
	osinfo.TotalMemory = buff.Total / 1048576

	diskinfo, _ := disk.Usage("/home/shafat-jamil")
	osinfo.TotalDisk = diskinfo.Total / 1048576
	
	fmt.Print("OS: ", osinfo.Distribution, " ", osinfo.Version, "\n")
	fmt.Print("Architecture: ", osinfo.Architecture, "\n")
	fmt.Print("Kernel: ", osinfo.KernelVersion, "\n")
	fmt.Print("CPU(s): ", osinfo.CPUCount, "\n")
	fmt.Print("CPU Model: ", osinfo.CPUModel, "\n")
	fmt.Print("Total Memory: ", osinfo.TotalMemory, " MB\n")
    fmt.Print("Disk info: ", osinfo.TotalDisk, " MB\n")
}

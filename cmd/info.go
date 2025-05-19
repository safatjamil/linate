package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/shirou/gopsutil/process"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/spf13/cobra"
)

func init() {
	infoCmd.AddCommand(osCmd)
	infoCmd.AddCommand(memCmd)
	infoCmd.AddCommand(loadCmd)
	infoCmd.AddCommand(processCmd)
	processCmd.Flags().StringP("type", "t", "", "Get process information by memory and CPU usage. Available options are mem and cpu.")
	processCmd.MarkFlagRequired("type")
	netCmd.AddCommand(netDetailsCmd)
}

var colors = map[string]string{
	"black":   "\033[30m",
	"red":     "\033[31m",
	"green":   "\033[32m",
	"yellow":  "\033[33m",
	"blue":    "\033[34m",
	"magenta": "\033[35m",
	"cyan":    "\033[36m",
	"white":   "\033[37m",
	"reset":   "\033[0m",
}

var file = map[string]string{
	"memory": "/proc/meminfo",
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Information about the os, memory, disk usage and other metrics",
	Long:  `Information about the os, memory, disk usage and other metrics. Run linate info --help for more options.`,
}

var osCmd = &cobra.Command{
	Use:   "os",
	Short: "Information about the os",
	Long:  `Information about distribution, architecture, CPU, memory, and disk.`,
	Run:   os_info,
}

var memCmd = &cobra.Command{
	Use:   "memory",
	Short: "Information about the memory",
	Long:  `Memory available, memory free, buffers, cached, swap, and others.`,
	Run:   memory_info,
}

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Information about the system's load",
	Long:  `load1, load5, and load15`,
	Run:   load_info,
}

var processCmd = &cobra.Command{
	Use:   "process",
	Short: "Information about the system processes. Get CPU usage and memory usage.",
	Long:  `Information about the system processes. Get CPU usage and memory usage.`,
	Run:   process_info,
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

type MemoryInfo struct {
	MemTotal     int
	MemFree      int
	MemAvailable int
	Buffers      int
	Cached       int
	SwapCached   int
	Active       int
	Inactive     int
	SwapTotal    int
	SwapFree     int
}

type ProcessInfo struct {
	PID         int32
	User        string
	Name        string
	MemoryUsage float32
	CPUUsage    float64
}

type NetDetails struct {
	InterfaceName string
	MacAddress    string
	IpAddress     string
}

func os_info(cmd *cobra.Command, args []string) {
	osinfo := OsInfo{}
	osinfo.Architecture = runtime.GOARCH
	osinfo.Distribution, _, osinfo.Version, _ = host.PlatformInformation()
	osinfo.KernelVersion, _ = host.KernelVersion()

	cpuinfo, _ := cpu.Info()
	osinfo.CPUCount, _ = cpu.Counts(true)
	osinfo.CPUModel = cpuinfo[0].ModelName

	buff, _ := mem.VirtualMemory()
	osinfo.TotalMemory = buff.Total / 1048576

	diskinfo, _ := disk.Usage("/")
	osinfo.TotalDisk = diskinfo.Total / 1048576

	title := [7]string{"OS", "Architecture", "Kernel", "CPU(s)", "CPU Mpdel", "Total Memoy", "Disk Size"}
	text_color := colors["yellow"]
	reset_color := colors["reset"]

	fmt.Printf("%-20s %s%s %s%s\n", title[0], text_color, osinfo.Distribution, osinfo.Version, reset_color)
	fmt.Printf("%-20s %s%s%s\n", title[1], text_color, osinfo.Architecture, reset_color)
	fmt.Printf("%-20s %s%s%s\n", title[2], text_color, osinfo.KernelVersion, reset_color)
	fmt.Printf("%-20s %s%d%s\n", title[3], text_color, osinfo.CPUCount, reset_color)
	fmt.Printf("%-20s %s%s%s\n", title[4], text_color, osinfo.CPUModel, reset_color)
	fmt.Printf("%-20s %s%d MB%s\n", title[5], text_color, osinfo.TotalMemory, reset_color)
	fmt.Printf("%-20s %s%d MB%s\n", title[6], text_color, osinfo.TotalDisk, reset_color)
}

func memory_info(cmd *cobra.Command, args []string) {
	var memory MemoryInfo
	memory = GetMemoryInfo()

	title := [10]string{"Memory Total", "Memory Free", "Memory Available", "Buffers", "Cached", "SwapCached", "Active", "Inactive", "SwapTotal", "SwapFree"}
	text_color := colors["yellow"]
	reset_color := colors["reset"]

	fmt.Printf("%-20s %s%d MB%s\n", title[0], text_color, memory.MemTotal, reset_color)
	fmt.Printf("%-20s %s%d MB%s\n", title[1], text_color, memory.MemFree, reset_color)
	fmt.Printf("%-20s %s%d MB%s\n", title[2], text_color, memory.MemAvailable, reset_color)
	fmt.Printf("%-20s %s%d MB%s\n", title[3], text_color, memory.Buffers, reset_color)
	fmt.Printf("%-20s %s%d MB%s\n", title[4], text_color, memory.Cached, reset_color)
	fmt.Printf("%-20s %s%d MB%s\n", title[5], text_color, memory.SwapCached, reset_color)
	fmt.Printf("%-20s %s%d MB%s\n", title[6], text_color, memory.Active, reset_color)
	fmt.Printf("%-20s %s%d MB%s\n", title[7], text_color, memory.Inactive, reset_color)
	fmt.Printf("%-20s %s%d MB%s\n", title[8], text_color, memory.SwapTotal, reset_color)
	fmt.Printf("%-20s %s%d MB%s\n", title[9], text_color, memory.SwapFree, reset_color)
}

func GetMemoryInfo() MemoryInfo {
	f, e := os.Open(file["memory"])
	if e != nil {
		log.Fatal("Pleasde check the permission of the file /proc/meminfo. It must have read permission for 'others'", e)
		defer f.Close()
		os.Exit(1)
	}
	scanner := bufio.NewScanner(f)
	res := MemoryInfo{}
	for scanner.Scan() {
		key, value := parseLineForMemory(scanner.Text())
		if key == "MemTotal" {
			res.MemTotal = value / 1024
		} else if key == "MemFree" {
			res.MemFree = value / 1024
		} else if key == "MemAvailable" {
			res.MemAvailable = value / 1024
		} else if key == "Buffers" {
			res.Buffers = value / 1024
		} else if key == "Cached" {
			res.Cached = value / 1024
		} else if key == "SwapCached" {
			res.SwapCached = value / 1024
		} else if key == "Active" {
			res.Active = value / 1024
		} else if key == "Inactive" {
			res.Inactive = value / 1024
		} else if key == "SwapTotal" {
			res.SwapTotal = value / 1024
		} else if key == "SwapFree" {
			res.SwapFree = value / 1024
		}
	}
	return res
}

func parseLineForMemory(raw string) (key string, value int) {
	res := strings.Split(strings.ReplaceAll(raw[:len(raw)-2], " ", ""), ":")
	return res[0], toInt(res[1])
}

func load_info(cmd *cobra.Command, args []string) {
	l, e := load.Avg()
	if e != nil {
		log.Fatal("Can not read load information", e)
		os.Exit(1)
	}
	title := [3]string{"Load1", "Load5", "Load15"}
	text_color := colors["yellow"]
	reset_color := colors["reset"]
	fmt.Printf("%-15s %s%f%s\n", title[0], text_color, l.Load1, reset_color)
	fmt.Printf("%-15s %s%f%s\n", title[1], text_color, l.Load5, reset_color)
	fmt.Printf("%-15s %s%f%s\n", title[2], text_color, l.Load15, reset_color)
}

func process_info(cmd *cobra.Command, args []string) {
	typ, _ := cmd.Flags().GetString("type")
	var display_type string
	if typ != "mem" && typ != "cpu" {
		log.Fatal("Incorrect value for the flag --type. Available options are mem and cpu.")
		os.Exit(1)
	}
	processes, e := process.Processes()
	if e != nil {
		log.Fatal("Can not read information about the processes", e)
		os.Exit(1)
	}

	// Load the data into proc
	var proc = make([]ProcessInfo, len(processes))
	counter := 0
	for _, p := range processes {
		process := ProcessInfo{}
		process.PID = p.Pid
		process.User, _ = p.Username()
		process.Name, _ = p.Name()
		process.CPUUsage, _ = p.CPUPercent()
		process.MemoryUsage, _ = p.MemoryPercent()
		// process.CreationTime, _ = p.CreateTime()
		proc[counter] = process
		counter += 1
	}

	// Sort by memory usage
	if typ == "mem" {
		display_type = "Memory Usage(%)"
		sort.Slice(proc, func(i, j int) bool {
			return proc[i].MemoryUsage > proc[j].MemoryUsage
		})
	} else if typ == "cpu" {
		display_type = "CPU Usage(%)"
		sort.Slice(proc, func(i, j int) bool {
			return proc[i].CPUUsage > proc[j].CPUUsage
		})
	}
	viewLength := 15
	if viewLength > len(proc) {
		viewLength = len(proc)
	}

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Process ID", "Name", "User", display_type)
	if typ == "mem" {
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
		for i := 0; i < viewLength; i++ {
			memuse := fmt.Sprintf("%.2f", proc[i].MemoryUsage)
			tbl.AddRow(proc[i].PID, proc[i].Name, proc[i].User, memuse)
		}
	} else if typ == "cpu" {
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
		for i := 0; i < viewLength; i++ {
			cpuuse := fmt.Sprintf("%.2f", proc[i].CPUUsage)
			tbl.AddRow(proc[i].PID, proc[i].Name, proc[i].User, cpuuse)
		}
	}

	tbl.Print()
}

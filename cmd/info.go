package cmd

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"

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
	infoCmd.AddCommand(usersCmd)
	processCmd.Flags().StringP("sort", "s", "", "Get process information by memory and CPU usage. Available options are mem and cpu.")
	processCmd.MarkFlagRequired("sort")
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
	Short: "Information about the os, memory, disk and other metrics.",
	Long:  `Information about the os, memory, disk usage and other metrics. Run linate info --help for more options.`,
}

var osCmd = &cobra.Command{
	Use:   "os",
	Short: "Information about the os.",
	Long:  `Information about distribution, architecture, CPU, memory, and disk.`,
	Run:   os_info,
}

var memCmd = &cobra.Command{
	Use:   "memory",
	Short: "Information about the memory.",
	Long:  `Memory available, memory free, buffers, cached, swap, and others.`,
	Run:   memory_info,
}

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Information about the system's load.",
	Long:  `load1, load5, and load15.`,
	Run:   load_info,
}

var processCmd = &cobra.Command{
	Use:   "process",
	Short: "Information about the system processes. Get processes by CPU and memory usage.",
	Long:  `Information about the system processes. Get processes by CPU and memory usage.`,
	Run:   process_info,
}

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Information about system users, last login time and their privilege.",
	Long:  `Information about system users, last login time and their privilege.`,
	Run:   users_info,
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
	// Read the /proc/meminfo file
	f, e := os.Open(file["memory"])
	if e != nil {
		exitWithError("Pleasde check the permission of the file /proc/meminfo. It must have read permission for 'others'")
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
		exitWithError("Can not read load information")
	}
	title := [3]string{"Load1", "Load5", "Load15"}
	text_color := colors["yellow"]
	reset_color := colors["reset"]
	fmt.Printf("%-15s %s%f%s\n", title[0], text_color, l.Load1, reset_color)
	fmt.Printf("%-15s %s%f%s\n", title[1], text_color, l.Load5, reset_color)
	fmt.Printf("%-15s %s%f%s\n", title[2], text_color, l.Load15, reset_color)
}

type ProcessInfo struct {
	PID          int32
	User         string
	Name         string
	MemoryUsage  float32
	CPUUsage     float64
	CreationTime string
}


func process_info(cmd *cobra.Command, args []string) {
	srt, _ := cmd.Flags().GetString("sort")
	var display_type string

	if srt != "mem" && srt != "cpu" && srt != "longrun" {
		exitWithError("Incorrect value for the flag --sort. Available options are mem, cpu, and longrun.\n")
	}
	processes, e := process.Processes()
	if e != nil {
		exitWithError("Can not read information about the processes")
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
		cT, _ := p.CreateTime()
		procTime := time.UnixMilli(cT)
		process.CreationTime = fmt.Sprintf("%v-%v-%v %v:%v", procTime.Year(), procTime.Month(), procTime.Day(), procTime.Hour(), procTime.Minute())
		proc[counter] = process
		counter += 1
	}

	// Sort by memory usage
	if srt == "mem" {
		display_type = "Memory Usage(%)"
		sort.Slice(proc, func(i, j int) bool {
			return proc[i].MemoryUsage > proc[j].MemoryUsage
		})
	} else if srt == "cpu" {
		display_type = "CPU Usage(%)"
		sort.Slice(proc, func(i, j int) bool {
			return proc[i].CPUUsage > proc[j].CPUUsage
		})
	} else if srt == "longrun" {
		display_type = "CPU(%) | Memory(%)"
	}
    
	// Display length
	viewLength := 30
	if viewLength > len(proc) {
		viewLength = len(proc)
	}

	// Create the table view
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("Process ID", "Name", "User", display_type, "Started")
	if srt == "mem" {
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
		for i := 0; i < viewLength; i++ {
			memuse := fmt.Sprintf("%.2f", proc[i].MemoryUsage)
			tbl.AddRow(proc[i].PID, proc[i].Name, proc[i].User, memuse, proc[i].CreationTime)
		}
	} else if srt == "cpu" {
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
		for i := 0; i < viewLength; i++ {
			cpuuse := fmt.Sprintf("%.2f", proc[i].CPUUsage)
			tbl.AddRow(proc[i].PID, proc[i].Name, proc[i].User, cpuuse, proc[i].CreationTime)
		}
	} else if srt == "longrun" {
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
		for i := 0; i < viewLength; i++ {
			cpuuse := fmt.Sprintf("%.2f", proc[i].CPUUsage)
			memuse := fmt.Sprintf("%.2f", proc[i].MemoryUsage)
			cm := fmt.Sprintf("%s | %s", cpuuse, memuse)
			tbl.AddRow(proc[i].PID, proc[i].Name, proc[i].User, cm, proc[i].CreationTime)
		}
	}
	tbl.Print()
}


var shells = [] string {"bash", "csh", "ksh", "mksh", "oksh", "sh", "tcsh", "yash", "zsh"}
var rootUserGroups = [] string {"sudo", "sudoers", "admin", "wheel", "staff"}
type userInfo struct {
	userName       string
	userID         string
	groupID        string
	description    string
	homeDirectory  string
	shell          string
}


func users_info(cmd *cobra.Command, args []string) {
	currYear :=  time.Now().Year()
	currentUser, _ := getCurrentUser()
	// Get all root users by group
	su := ""
	for _, g := range rootUserGroups{
        h, _ := getUsersByGroup(g)
		if h != "" {
			su += su + "," + h
		}
	}
	superusers := strings.Split(su, ",")
	usrs , e, size := list_users()
	if e != nil {
		fmt.Printf("%v\n", e)
		os.Exit(0)
	}
	// Find real users. Must have a terminal
	realUsers := make([]userInfo, size)
	counter := 0
	realUsersCounter := 0
	for _, v := range usrs {
		if v.shell != "" {
		  shl := strings.Split(v.shell, "/")
		  if len(shl) >1 && arrContains(shells, shl[len(shl)-1]) {
			realUsers[realUsersCounter].userName = v.userName
			realUsers[realUsersCounter].userID = v.userID
			realUsers[realUsersCounter].groupID = v.groupID
			realUsers[realUsersCounter].description = v.description
			realUsers[realUsersCounter].shell = v.shell
			realUsersCounter += 1
		  }
		}
		counter += 1
		if counter == size {
			break
		}
	}
    
	if realUsersCounter > 1 {
		headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
		columnFmt := color.New(color.FgYellow).SprintfFunc()
		tbl := table.New("Username", "userID", "Description", "Shell", "Last Login", "Root Privilege")
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
		for i := 0; i < realUsersCounter; i++ {
			rootPrivilege := "No"
			if arrContains(superusers, realUsers[i].userName) || realUsers[i].userName == "root" {
				rootPrivilege = "Yes"
			}
			if currentUser != "" && currentUser == realUsers[i].userName {
				tbl.AddRow(realUsers[i].userName, realUsers[i].userID, realUsers[i].description, realUsers[i].shell, "Logged in now", rootPrivilege)
				continue
			}
			lTime, e := getLastLogin(realUsers[i].userName)
			if e!= nil && (currYear - lTime.Year() <= 1){
			  tbl.AddRow(realUsers[i].userName, realUsers[i].userID, realUsers[i].description, realUsers[i].shell, fmt.Sprintf("%v", lTime), rootPrivilege)
			} else {
			  tbl.AddRow(realUsers[i].userName, realUsers[i].userID, realUsers[i].description, realUsers[i].shell, fmt.Sprintf("%v", e), rootPrivilege)
			}
		}
		tbl.Print()
    } else {
		fmt.Printf("%sNo users found%s\n", colors["red"], colors["reset"])
	}
}
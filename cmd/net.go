package cmd

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
	"strings"

	"github.com/fatih/color"
	"github.com/jackpal/gateway"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"github.com/cakturk/go-netstat/netstat"
)

func init() {
	netCmd.AddCommand(netDetailsCmd)
	netCmd.AddCommand(netConnCmd)
}

var netCmd = &cobra.Command{
	Use:   "net",
	Short: "Information about the network interfaces, ports, and incoming/outgoing requests.",
	Long:  `Information about the network interfaces, ports, and incoming/outgoing requests. Run linate net --help for more options.`,
}

var netDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Information about the network interfaces and gateway.",
	Long:  `Information about the network interfaces and gateway.`,
	Run:   net_details_info,
}

var netConnCmd = &cobra.Command{
	Use:   "conn",
	Short: "Internet connection is available or not.",
	Long:  `Internet connection is available or not.`,
	Run:   net_conn_info,
}

var netGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Information about the listening connections.",
	Long:  `Information about the listening connections.`,
	Run:   net_get_info,
}

type NetDetails struct {
	InterfaceName string
	MacAddress    string
	IpAddress     string
}

type SockAddr struct {
	IP   net.IP
	Port uint16
}

type SockTabEntry struct {
	ino        string
	LocalAddr  *SockAddr
	RemoteAddr *SockAddr
	State      SkState
	UID        uint32
	Process    *Process
}

type Process struct {
	Pid  int
	Name string
}

type SkState uint8

type Socket struct {
	ino             string
	localIP         string
	localPort       string
	remoteIP        string
	remotePort      string
	state     		uint8
	UID             uint32 
	PID             int
	processName     string
	IP         		net.IP
	port            uint16
	conType         string
}

var portMap = map[string]string{
	"es": "9200",
	"http": "80",
	"https": "443",
	"mysql": "3306",
	"nfs": "2049",
	"smtp": "25",
	"ssh": "22",
	"sunrpc": "11",

}
func net_details_info(cmd *cobra.Command, args []string) {
	interfaces, e := net.Interfaces()
	if e != nil {
		exitWithError("Failed to get network interface information.")
	}
	var intfc = make([]NetDetails, len(interfaces))
	counter := 0

	// loop through the network interfaces
	for _, inter := range interfaces {
		intfc[counter].InterfaceName = inter.Name
		address := ""
		// Get a list of IP addresses for this network interface
		addrs, e := inter.Addrs()
		if e != nil {
			lexitWithError("Failed to obtain IP address list.")
		}
		for _, addr := range addrs {
			address = fmt.Sprintf("%s %s \n", address, addr)
		}
		intfc[counter].IpAddress = address

		// Get the MAC address of the network interface
		mac := inter.HardwareAddr
		intfc[counter].MacAddress = fmt.Sprintf("%s", mac)
		counter += 1
	}
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("Interface Name", "MAC address", "IP Address(s)")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	for i := 0; i < len(intfc); i++ {
		tbl.AddRow(intfc[i].InterfaceName, intfc[i].MacAddress, intfc[i].IpAddress)
	}
	tbl.Print()

	gatewayIP, e := gateway.DiscoverGateway()
	if e != nil {
		exitWithError("Cannot read gateway information.")
	}
	fmt.Println("")
	title := [1]string{"Gateway"}
	text_color := colors["yellow"]
	reset_color := colors["reset"]
	fmt.Printf("%-15s %s%s%s\n", title[0], text_color, gatewayIP, reset_color)
}

func net_conn_info(cmd *cobra.Command, args []string) {
	_, e := net.DialTimeout("tcp", "8.8.8.8:53", 15*time.Second)
	if e != nil {
		fmt.Printf("%sInternet connection is not available%s\n", colors["red"], colors["reset"])
	} else {
		fmt.Printf("%sInternet connection is available%s\n", colors["green"], colors["reset"])
	}
}

func net_get_info(cmd *cobra.Command, args []string) {
	// UDP sockets
	socks, e := netstat.UDPSocks(netstat.NoopFilter)
	if e != nil {
		exitWithError("Failed to obtain UDP connections.")
		os.Exit(1)
	}
	var socket = make([]Socket, len(socks))
	counter := 0
	for _, u := range socks {
		fmt.Printf("%v\n", u)
		la := strings.Split(fmt.Sprintf("%s", u.LocalAddr), ":")
		//ra := strings.Split(fmt.Sprintf("%s", u.RemoteAddr), ":")
		socket[counter].localIP = la[0]
		socket[counter].localPort = la[1]

		// socket[counter].remoteAddr = fmt.Sprintf("%s", u.RemoteAddr)
		// socket[counter].
		fmt.Printf("%s", la[0])
		fmt.Printf("\n")
		counter += 1
	}
    
	for _, v := range socket {
		fmt.Printf("%s", v.localIP)
		fmt.Printf("\n")
	}
	//TCP sockets
	// socks, e = netstat.TCPSocks(netstat.NoopFilter)
	// if e != nil {
	// 	log.Fatal("Failed to obtain TCP connections: ", e)
	// 	os.Exit(1)
	// }
	// for _, t := range socks {
	// 	fmt.Printf("%T\n", t)
	// }

	// // get only listening TCP sockets
	// tabs, err := netstat.TCPSocks(func(s *netstat.SockTabEntry) bool {
	// 	return s.State == netstat.Listen
	// })
	// if err != nil {
	// 	return err
	// }
	// for _, e := range tabs {
	// 	fmt.Printf("%v\n", e)
	// }

	// // list all the TCP sockets in state FIN_WAIT_1 for your HTTP server
	// tabs, err = netstat.TCPSocks(func(s *netstat.SockTabEntry) bool {
	// 	return s.State == netstat.FinWait1 && s.LocalAddr.Port == 80
	// })
	// // error handling, etc.

	// return nil
}

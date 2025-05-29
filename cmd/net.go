package cmd

import (
	"fmt"
	_ "log"
	"net"
	_"os"
	"strings"
	"time"
	"strconv"

	"github.com/fatih/color"
	"github.com/jackpal/gateway"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"github.com/bastjan/netstat"
)

func init() {
	netCmd.AddCommand(netDetailsCmd)
	netCmd.AddCommand(netConnCmd)
	netCmd.AddCommand(netInfoCmd)
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
	Short: "Check the internet connection.",
	Long:  `Check the internet connection.`,
	Run:   net_conn_info,
}

var netInfoCmd = &cobra.Command{
	Use:   "socket",
	Short: "Information about the network sockets. Listening, established and others.",
	Long:  `Information about the network sockets. Listening, established and others.`,
	Run:   net_socket,
}

type NetDetails struct {
	InterfaceName string
	MacAddress    string
	IpAddress     string
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
			exitWithError("Failed to obtain IP address list.")
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
	ino         string
	localIP     string
	localPort   string
	remoteIP    string
	remotePort  string
	state       uint8
	UID         uint32
	PID         int
	processName string
	IP          net.IP
	port        uint16
	conType     string
}

var portMap = map[string]string{
    "11": "SunRPC",
	"22": "SSH",
	"25": "SMTP",
	"53": "DNS",
	"80": "HTTP Service",
	"111": "RPC (Remote Procedure Call)",
	"443": "HTTPS Service",
	"631": "Internet Printing Protocol",
	"1433": "MSSQL",
	"1521": "Oracle DB",
	"2049": "NFS",
	"2377": "Docker Swarm",
	"2379": "Etcd Client",
	"2380": "Etcd Peer Communication",
	"2382": "SQL Server Analysis Services (SSAS)",
	"2525": "SMTP",
	"3000": "Grafana",
	"3306": "MySQL",
	"3389": "RDP",
	"5000": "Python Flask",
	"5355": "LLMNR (Link-Local Multicast Name Resolution)",
	"5380": "Technitium DNS",
	"5432": "PostgreSQL",
	"5666": "NRPE",
	"5672": "RabbitMQ",
	"6379": "Redis",
	"6443": "Kube Master",
	"7070": "Real Time Streaming",
	"8000": "Django",
	"8006": "Proxmox",
	"8080": "Jenkins / Other HTTP service",
	"8086": "InfluxDB",
	"9000": "Graylog",
	"9042": "Cassandra",
	"9090": "Prometheus",
	"9091": "Prometheus Pushgateway",
	"9092": "Kafka",
	"9200": "Elasticsearch",
	"9300": "Elasticsearch Cluster Communication",
	"9411": "Zipkin",
	"9876": "OpenStack Loadbalancer service (Octavia)",
	"10050": "Zabbix Agent",
    "27017": "MongoDB",
}
var showSocketArray = []string {"protocol", "state", "localIp", "lServiceOrPort", "remoteIp", "rServiceOrPort", "observation"}
var header = []string{"Proto", "Local Address", "Foreign Address", "State", "User"}

// Helper functions
func formatConnections(loc *netstat.Protocol) [][]string {
	connections, _ := loc.Connections()
	results := make([][]string, 0, len(connections))
	for _, conn := range connections {
		if !isReq(conn) {
			continue
		}
		results = append(results, []string{
			conn.Protocol.Name,
			fmt.Sprintf("%s:%s", conn.IP, formatPort(conn.Port)),
			fmt.Sprintf("%s:%s", conn.RemoteIP, formatPort(conn.RemotePort)),
			conn.State.String(),
			conn.UserID,
		})
	}

	return results
}

func isReq(conn *netstat.Connection) bool {
	tcpConn := strings.HasPrefix(conn.Protocol.Name, "tcp") && (conn.State == netstat.TCPListen || conn.State == netstat.TCPEstablished)
	udpListen := strings.HasPrefix(conn.Protocol.Name, "udp") && (conn.State == netstat.TCPListen || conn.State == netstat.TCPEstablished)
	return tcpConn || udpListen
}

func formatPort(port int) string {
	if port == 0 {
		return "*"
	}
	return strconv.Itoa(port)
}

func net_socket(cmd *cobra.Command, args []string) {
	sshCon := make(map[string][]string)
	listenCon := make(map[string][]string)
	establishedCon := make(map[string][]string)
    // Get all sockets
	out := [][]string{header}
	out = append(out, formatConnections(netstat.TCP)...)
	out = append(out, formatConnections(netstat.TCP6)...)
	out = append(out, formatConnections(netstat.UDP)...)
	out = append(out, formatConnections(netstat.UDP6)...)

	for _, v := range out {
		state, _ := v[3], v[4]
        localIp, localPort := splitOnLast(v[1], ":")
		remoteIp, remotePort := splitOnLast(v[2], ":")
        
		// Socket is open for an expected service(see portMap). We will not show all ports that are open
		if !portMapContains(portMap, localPort) && !portMapContains(portMap, remotePort){
			continue
		}
        
		lService := localPort
	    rService := remotePort
		observation := ""
		// Map port to service name
        if portMapContains(portMap, localPort){ lService = portMap[localPort]}
		if portMapContains(portMap, remotePort){ rService = portMap[remotePort]}

		// Add ssh connections (listening and established) to sshCon
		if state == "LISTEN" {
			if localPort == "22"{
				if ! conMapContains(sshCon, localIp){
                    sshCon[localIp] = []string{v[0], v[3], localIp, lService, remoteIp, rService, fmt.Sprintf("Listening on %s", lService)}
				}
			} else {
				conj := localIp + localPort
				if ! conMapContains(listenCon, conj){
                    listenCon[conj] = []string{v[0], v[3], localIp, lService, remoteIp, rService, fmt.Sprintf("Listening on %s | Forwards to %s", lService, rService)}
				}
			}
		} else if state == "ESTABLISHED"{
			if localPort == "22" || remotePort == "22" {
				conj := localIp + localPort + remoteIp + remotePort
				if !conMapContains(sshCon, conj){
					user  := v[4]
					userId, e := strconv.Atoi(v[4])
					if e == nil {
					    u, e := getUsernameFromUID(userId)
					    if e == nil {
                            user = u
					    }
				    }
					if localPort == "22" {observation = fmt.Sprintf("SSH connection from %s", remoteIp)}
					if remotePort == "22" {observation = fmt.Sprintf("SSH connection to %s | User: %s", remoteIp, user)}
					sshCon[conj] = []string{v[0], v[3], localIp, lService, remoteIp, rService, observation}
				}
			} else {
				conj := localIp + localPort + remoteIp + remotePort
				if !conMapContains(establishedCon, conj){
                    establishedCon[conj] = []string{v[0], v[3], localIp, lService, remoteIp, rService, ""}
				}
			}
		}
	}

    // Create the table view
	viewLength := 40
	totalLength := len(sshCon) + len(listenCon) + len(establishedCon)
	if totalLength < viewLength {
		viewLength = totalLength
	}

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("Protocol", "State", "LocalIP", "Local Service|Port", "RemoteIP", "Remote Service|Port", "Observation")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	counter := 0

	// Show ssh connections first
	for _, v := range sshCon {
		if counter == viewLength {
			break
		}
		tbl.AddRow(v[0], v[1], v[2], v[3], v[4], v[5], v[6])
		counter += 1
	}
    
	for _, v := range listenCon {
		if counter == viewLength {
			break
		}
		tbl.AddRow(v[0], v[1], v[2], v[3], v[4], v[5], v[6])
		counter += 1
	}

	for _, v := range establishedCon {
		if counter == viewLength {
			break
		}
		tbl.AddRow(v[0], v[1], v[2], v[3], v[4], v[5], v[6])
		counter += 1
	}
	tbl.Print()
}

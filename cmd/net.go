package cmd

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/fatih/color"
	"github.com/jackpal/gateway"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

func init() {
	netCmd.AddCommand(netDetailsCmd)
}

var netCmd = &cobra.Command{
	Use:   "net",
	Short: "Information about the network interfaces, ports, and incoming/outgoing requests.",
	Long:  `Information about the network interfaces, ports, and incoming/outgoing requests.. Run linate net --help for more options.`,
}

var netDetailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Information about the network interfaces",
	Long:  `Information about the network interfaces`,
	Run:   net_details_info,
}

func net_details_info(cmd *cobra.Command, args []string) {
	interfaces, e := net.Interfaces()
	if e != nil {
		log.Fatal("Failed to get network interface information: ", e)
		os.Exit(1)
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
			log.Fatal("Failed to obtain IP address list: ", e)
			os.Exit(1)
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
		log.Fatal("Cannot read gateway information.", e)
		os.Exit(1)
	}
	fmt.Println("")
	title := [1]string{"Gateway"}
	text_color := colors["yellow"]
	reset_color := colors["reset"]
	fmt.Printf("%-15s %s%s%s\n", title[0], text_color, gatewayIP, reset_color)
}

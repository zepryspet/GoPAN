package main

import (
	"fmt"
	"strings"
	"github.com/spf13/cobra"
    "./cps"
)

func main() {
	var firewallIP string
    var community string
    var minutes int
    
	var cmdScript = &cobra.Command{
		Use:	 "run",
		Short: "Pre-built scripts to collect and process firewall data",
		Long: `Set of scripts that will collect and process firewall 
information`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Print: " + strings.Join(args, " "))
		},
	}

	var cmdCPS = &cobra.Command{
        Use:	 "cps",
		Short: "Generate a CPS report from a Palo alto firewall",
		Long: `Send an SNMP query to get the CPS per zone information
and generate an excel table with the historical data. 
Recommended to run for a week. Only PAN-OS 8.0+ is supported`,
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("firewall IP: " + firewallIP)
                cps.Snmpgen(firewallIP, community, minutes)
		},
	}

	cmdCPS.Flags().StringVarP(&firewallIP, "ip-address", "i", "", "firewall IP address or FQDN")
	cmdCPS.MarkFlagRequired("ip-address")
    cmdCPS.Flags().StringVarP(&community, "community", "c", "", "SNMPv2 community")
	cmdCPS.MarkFlagRequired("community")
    cmdCPS.Flags().IntVarP(&minutes, "minutes", "m", 5, "SNMP interval in minutes to collect LPS")

	var rootCmd = &cobra.Command{Use: "pan"}
	rootCmd.AddCommand(cmdScript)
	cmdScript.AddCommand(cmdCPS)
	rootCmd.Execute()
}
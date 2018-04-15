package main

import (
	"fmt"
	"strings"
	"github.com/spf13/cobra"
    "GoPAN/run/ssh"
    "GoPAN/run/cps"
)

func main() {
    //Vriables used for flags on scripts
	var firewallIP string
    var community string
    var pass string
    var user string
    var minutes int
    var flag1 bool
    
    //Main commands under root
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

    //Run command subsection
    //Run CPS
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

    //Run SSH
	cmdCPS.Flags().StringVarP(&firewallIP, "ip-address", "i", "", "firewall IP address or FQDN")
	cmdCPS.MarkFlagRequired("ip-address")
    cmdCPS.Flags().StringVarP(&community, "community", "c", "", "SNMPv2 community")
	cmdCPS.MarkFlagRequired("community")
    cmdCPS.Flags().IntVarP(&minutes, "minutes", "m", 5, "SNMP interval in minutes to collect LPS")
    
    	var cmdSSH = &cobra.Command{
        Use:	 "ssh",
		Short: "send ssh commands to a PAN firewall",
		Long: `Send SSH commands to a PAN firewall
        Either send a single command or send multiple commands from a text file.
        Either one time or in an loop`,
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("firewall IP: " + firewallIP)
                fmt.Printf("%t\n", flag1)
                panssh.Send(firewallIP, user, pass)
		},
	}
    //missing timeout, filename or command
    cmdSSH.Flags().StringVarP(&firewallIP, "ip-address", "i", "", "firewall IP address or FQDN")
	cmdSSH.MarkFlagRequired("ip-address")
    cmdSSH.Flags().StringVarP(&user, "user", "u", "", "firewall username")
	cmdSSH.MarkFlagRequired("user")
    cmdSSH.Flags().StringVarP(&pass, "password", "p", "", "firewall password")
	cmdSSH.MarkFlagRequired("password")
    cmdSSH.Flags().BoolVarP(&flag1, "isFile", "f", false , "set this flag to true if you're sending commands in a text file")

    //Setting command and subcommand structure
    
	var rootCmd = &cobra.Command{Use: "pan"}
	rootCmd.AddCommand(cmdScript)
    //Run sub-commands
	cmdScript.AddCommand(cmdCPS)
    cmdScript.AddCommand(cmdSSH)
	rootCmd.Execute()
}
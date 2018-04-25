package main

import (
	"fmt"
	"strings"
	"github.com/spf13/cobra"
    "GoPAN/run/ssh"
    "GoPAN/run/cps"
    "GoPAN/api/urlcat"
    "GoPAN/utils"
    "time"
)

func main() {
    //Variables used for flags
	var firewallIP string
    var community string
    var pass string
    var user string
    var command string
    var minutes int
    var seconds int
    var flag1 bool
    var flag2 bool
    
    //"Run" top command section
    //
	var cmdScript = &cobra.Command{
		Use:	 "run",
		Short: "Pre-built scripts to collect and process firewall data using non-api methods like SNMP or SSH",
		Long: `Set of scripts that will collect and process firewall 
information`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Print: " + strings.Join(args, " "))
		},
	}

    //"run cps"
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
    //Seting requiered flags
	cmdCPS.Flags().StringVarP(&firewallIP, "ip-address", "i", "", "firewall IP address or FQDN")
	cmdCPS.MarkFlagRequired("ip-address")
    cmdCPS.Flags().StringVarP(&community, "community", "c", "", "SNMPv2 community")
	cmdCPS.MarkFlagRequired("community")
    cmdCPS.Flags().IntVarP(&minutes, "minutes", "m", 5, "SNMP interval in minutes to collect LPS")

    //"run ssh"
    var cmdSSH = &cobra.Command{
        Use:	 "ssh",
		Short: "send ssh commands to a PAN firewall",
		Long: `Send SSH commands to a PAN firewall
        Either send a single command or send multiple commands from a text file.
        Either one time or in an loop`,
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
            if seconds == 0 {
                panssh.Send(firewallIP, user, pass, command, flag1, flag2)
            }   else{
                for true{
                    panssh.Send(firewallIP, user, pass, command, flag1, flag2)
                    time.Sleep(time.Duration(seconds)*time.Second)
                }
            } 
            
		},
	}
    //missing timeout, filename or command
    cmdSSH.Flags().StringVarP(&firewallIP, "ip-address", "i", "", "firewall IP address or FQDN")
	cmdSSH.MarkFlagRequired("ip-address")
    cmdSSH.Flags().StringVarP(&user, "user", "u", "", "firewall username")
	cmdSSH.MarkFlagRequired("user")
    cmdSSH.Flags().StringVarP(&pass, "password", "p", "", "firewall password")
	cmdSSH.MarkFlagRequired("password")
    cmdSSH.Flags().StringVarP(&command, "command", "r", "", "ssh command or file with commands to run on the firewall")
	cmdSSH.MarkFlagRequired("password")
    cmdSSH.Flags().BoolVarP(&flag1, "isFile", "f", false , "set this flag if you're sending commands in a text file")
    cmdSSH.Flags().BoolVarP(&flag2, "isConfig", "c", false , "set this flag if you're sending configuration commands")
    cmdSSH.Flags().IntVarP(&seconds, "times", "t", 0 , "time in seconds to repeat the commands indefinetily")


    //"api" subtop command section
    //
	var cmdApi = &cobra.Command{
		Use:	 "api",
        TraverseChildren: true,
		Short: "scripts using the api calls",
		Long: `Set of scripts that will use the api calls -http- and process firewall or panorama data using one of the available commands`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Print: " + strings.Join(args, " "))
		},
	}
   //Adding fqdn, username and password as required flags for all api commands
    cmdApi.PersistentFlags().StringVarP(&firewallIP, "ip-address", "i", "", "firewall IP address or FQDN")
    cmdApi.PersistentFlags().StringVarP(&user, "user", "u", "", "firewall username")
    cmdApi.PersistentFlags().StringVarP(&pass, "password", "p", "", "firewall password")

    var cmdKey = &cobra.Command{
        Use:	 "keygen",
		Short: "generate an API key",
		Long: "Generate an API key for the username provided",
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
            fmt.Println(pan.Keygen(firewallIP, user, pass))
		},
	}
    
    var cmdURL = &cobra.Command{
        Use:	 "urlcat",
		Short: "Get a url category for a website or from a file with website",
		Long: "Generate an API key for the username provided",
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
            urlcat.Request(firewallIP, pan.Keygen(firewallIP, user, pass), command, flag1)
		},
	}
    cmdURL.Flags().StringVarP(&command, "website", "w", "", "url or filename with urls")
	cmdURL.MarkFlagRequired("website")
    cmdURL.Flags().BoolVarP(&flag1, "isFile", "f", false , "set this flag to true if you're have the urls in a text file")
    
    //Setting command and subcommand structure
	var rootCmd = &cobra.Command{Use: "pan"}
	rootCmd.AddCommand(cmdScript)
    rootCmd.AddCommand(cmdApi)
    //Run sub-commands
	cmdScript.AddCommand(cmdCPS)
    cmdScript.AddCommand(cmdSSH)
    //Api sub-commands
    cmdApi.AddCommand(cmdKey)
    cmdApi.AddCommand(cmdURL)
	rootCmd.Execute()
}
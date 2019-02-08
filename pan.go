package main

import (
	"fmt"
	"strings"
	"github.com/spf13/cobra"
  "github.com/zepryspet/GoPAN/run/ssh"
  "github.com/zepryspet/GoPAN/run/cps"
  "github.com/zepryspet/GoPAN/api/urlcat"
  "github.com/zepryspet/GoPAN/api/cutover"
  "github.com/zepryspet/GoPAN/api/threat"
  "github.com/zepryspet/GoPAN/utils"
  "time"
)

func main() {
  //Variables used for flags
	var firewallIP string
  var community string
  var pass string
  var user string
  var command string
  var polltime int
  var seconds int
	var version int
	var privacyv3 string
	var authv3 string
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
                cps.Snmpgen(firewallIP, community, polltime, version, authv3, privacyv3)
		},
	}
    //Seting requiered flags
	cmdCPS.Flags().StringVarP(&firewallIP, "ip-address", "i", "", "firewall IP address or FQDN")
	cmdCPS.MarkFlagRequired("ip-address")
  cmdCPS.Flags().StringVarP(&community, "community", "c", "", "SNMPv2 community or SNMPv3 user")
	cmdCPS.MarkFlagRequired("community")
  cmdCPS.Flags().IntVarP(&polltime, "polltime", "s", 10, "SNMP interval in polltime to collect LPS")
	cmdCPS.Flags().IntVarP(&version, "version", "v", 2, "use 2 for SNMPv2 or 3 for SNMPv3")
	cmdCPS.Flags().StringVarP(&authv3, "authv3", "a", "", "Auth password in case snmpv3 is used")
	cmdCPS.Flags().StringVarP(&privacyv3, "privacyv3", "x", "", "password in case snmpv3 is used")


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


    //"api" top command section
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
   //Adding fqdn, username and password as persistend flags for all api sub-commands
    cmdApi.PersistentFlags().StringVarP(&firewallIP, "ip-address", "i", "", "firewall IP address or FQDN")
    cmdApi.PersistentFlags().StringVarP(&user, "user", "u", "", "firewall username")
    cmdApi.PersistentFlags().StringVarP(&pass, "password", "p", "", "firewall password")

    //"api keygen" command to generate an API key
    var cmdKey = &cobra.Command{
        Use:	 "keygen",
		Short: "generate an API key",
		Long: "Generate an API key for the username provided",
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
            fmt.Println(pan.Keygen(firewallIP, user, pass))
		},
	}

    //"api urlcat" command to check url categories
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

    //"api cutover" basic checks during maintenance windows.
    var cmdCut = &cobra.Command{
        Use:	 "cutover",
        Short: "Checklist for common issues during cutovers",
        Long: "Check incomplete arps, duplex mismatch, send GARPs on all interfaces, global counters, CRCs and enabled interfaces in down state",
        Args: cobra.MinimumNArgs(0),
        Run: func(cmd *cobra.Command, args []string) {
            cutover.Check(firewallIP,user, pass)
        },
	}

    //"api threat" things related to the threat db
    var cmdThreat = &cobra.Command{
        Use:	 "threat",
        Short: "Export the threat Database",
        Long: "Subcommand that export the current threat DB to csv for easier parsing",
        Args: cobra.MinimumNArgs(0),
        Run: func(cmd *cobra.Command, args []string) {
            threat.Export(firewallIP,user, pass)
        },
	}

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
    cmdApi.AddCommand(cmdCut)
    cmdApi.AddCommand(cmdThreat)
	rootCmd.Execute()
}

package cps

import (
	"fmt"
	g "github.com/soniah/gosnmp"
	"log"
	"os"
	"strconv"
	"GoPAN/utils"
	"time"
)

var zones []string
var cps []int

//function to SNMP walk on the  CPS OIDs and save the output into text files based on the zone name.
func Snmpgen(fqdn string, community string, minutes int) {
	min := time.Duration(60 * minutes)
	loop := true
	for loop {
		// Do verbose logging of packets.
		port, _ := strconv.ParseUint("161", 10, 16)
		params := &g.GoSNMP{
			Target:    fqdn,
			Port:      uint16(port),
			Community: community,
			Version:   g.Version2c,
			Timeout:   time.Duration(2) * time.Second,
			Logger:    log.New(os.Stdout, "", 0),
		}
		err := params.Connect()
		if err != nil {
			pan.Logerror(err)
		}
		defer params.Conn.Close()
		//oids := []string{".1.3.6.1.4.1.25461.2.1.2.3.10.1.1", ".1.3.6.1.4.1.25461.2.1.2.3.10.1.1.5.84.114.117.115.116"}
		var oid [3]string
		//zone names OID
		oidZone := ".1.3.6.1.4.1.25461.2.1.2.3.10.1.1"
		//TCP CPS per zone oID, returns unsigned 32bit integer
		oid[0] = ".1.3.6.1.4.1.25461.2.1.2.3.10.1.2"
		//UDP CPS per zone oID, returns unsigned 32bit integer
		oid[1] = ".1.3.6.1.4.1.25461.2.1.2.3.10.1.3"
		//other IP CPS per zone oID, returns unsigned 32bit integer
		oid[2] = ".1.3.6.1.4.1.25461.2.1.2.3.10.1.3"

		err = params.BulkWalk(oidZone, saveZone)
		if err != nil {
			pan.Logerror(err)
		}
		fmt.Printf("%v", zones)

		for _, element := range oid {
			err = params.BulkWalk(element, saveCPS)
			if err != nil {
				pan.Logerror(err)
			}
			fmt.Printf("%v", cps)
			//Saving CPS info to the respective zone name using CSV
			for i := range zones {
				pan.Wlog(zones[i]+".csv", strconv.Itoa(cps[i])+",", false)
			}
			//Clearing the slice to reuse it for the next SNMP  walk
			cps = nil
		}
		//Ading a breakline at the end of all zone statistic CPS zone files
		for i := range zones {
			pan.Wlog(zones[i]+".csv", "\n", false)
		}
		//clearing slice for re-use
		zones = nil
		time.Sleep(min * time.Second)
	}
}

func saveZone(pdu g.SnmpPDU) error {
	fmt.Printf("%s = ", pdu.Name)
	switch pdu.Type {
	case g.OctetString:
		b := pdu.Value.([]byte)
		zones = append(zones, string(b))
	default:
		pan.Wlog("error.txt", "received zone name NOT as string, please check the OIDs for CPS zone names", true)
		os.Exit(1)
	}
	return nil
}

func saveCPS(pdu g.SnmpPDU) error {
	fmt.Printf("%s = ", pdu.Name)
	switch pdu.Type {
	case g.OctetString:
		pan.Wlog("error.txt", "received CPS as string instead as int, please check the OIDs for CPS", true)
		os.Exit(1)
	default:
		cps = append(cps, pdu.Value.(int))
	}
	return nil
}

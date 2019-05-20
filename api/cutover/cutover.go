package cutover

import (
	"github.com/beevik/etree"
	"github.com/zepryspet/GoPAN/http"
	"github.com/zepryspet/GoPAN/run/ssh"
	"github.com/zepryspet/GoPAN/utils"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

func Check(fqdn string, user string, pass string) {
	var message string
	//Generating api Key
	apikey := pan.Keygen(fqdn, user, pass)
	req, err := url.Parse("https://" + fqdn + "/api/?")
	pan.Logerror(err, true)
	/////////////////////////
	//Check incomplete arps.
	//Generating HTTP req.
	cmd := (pan.CmdGen("show arp n_all"))
	q := url.Values{}
	q.Add("type", "op")
	q.Add("key", apikey)
	q.Add("cmd", cmd)
	req.RawQuery = q.Encode()
	//Sending the query to the firewall
	resp, err := pan.HttpValidate(req.String(), false)
	pan.Logerror(err, true)
	doc := etree.NewDocument()
	doc.ReadFromBytes(resp)
	for _, e := range doc.FindElements("./response/result/entries/*") {
		status := strings.TrimSpace(e.SelectElement("status").Text())
		if status == "i" {
			intf := e.SelectElement("interface").Text()
			println("incomplete arp found on: " + intf)
		}
	}
	///////////////////
	//Check CRCs
	//
	//Generating HTTP req.
	cmd = pan.CmdGen("show system state filter-pretty t_sys.*")
	q.Set("cmd", cmd)
	req.RawQuery = q.Encode()
	//Sending the query to the firewall
	resp, err = pan.HttpValidate(req.String(), false)
	pan.Logerror(err, true)
	//Regex will check the first capturing group that will define the physical por as 'sX.pY'
	r := regexp.MustCompile(`sys.(s\d.p\d).detail:\s{\s[^}]*crc`)
	ports := r.FindAllStringSubmatch(string(resp), -1)
	for i := 0; i < len(ports); i++ {
		index := strings.Split(ports[i][1], ".")
		slot := strings.TrimPrefix(index[0], "s")
		port := strings.TrimPrefix(index[1], "p")
		message = "CRCs detected in ethernet" + slot + "/" + port
		println(message)
		pan.Wlog("output.txt", message, true)
	}
	/////////////////////////////////
	///Check global Counters
	gcCheck(fqdn, apikey)

	/////////////////////////////////
	//Verify they interfaces aren't running either 100/10 Mbps or half duplex or are down.
	//Generating HTTP req.
	cmd = pan.CmdGen("show interface t_all")
	q.Set("cmd", cmd)
	req.RawQuery = q.Encode()
	//Sending the query to the firewall
	resp, err = pan.HttpValidate(req.String(), false)
	pan.Logerror(err, true)
	doc = etree.NewDocument()
	doc.ReadFromBytes(resp)
	//parsing response for hardware status
	for _, e := range doc.FindElements("./response/result/hw/*") {
		name := strings.TrimSpace(e.SelectElement("name").Text())
		//Checking status, speed and state for all ethernet interfaces
		if strings.HasPrefix(name, "ethernet") {
			speed := strings.TrimSpace(e.SelectElement("speed").Text())
			state := strings.TrimSpace(e.SelectElement("state").Text())
			duplex := strings.TrimSpace(e.SelectElement("duplex").Text())
			//intf := e.SelectElement("interface").Text()
			if state != "up" {
				println("interface " + name + " is DOWN!")
			} else {
				if duplex != "full" {
					message = "interface " + name + " is half duplex!"
					println(message)
					pan.Wlog("output.txt", message, true)
				}
				if speed == "10" || speed == "100" {
					message = "interface " + name + " speed is less than 1 Gbps. Current speed: " + speed
					println(message)
					pan.Wlog("output.txt", message, true)
				}
			}
		}
	}
	///////////////////////////////
	//Generate commands for GARPs from the previous "show interface all"
	//Removing GARP text file if it exist
	if _, err := os.Stat("garp.txt"); err == nil {
		err = os.Remove("garp.txt")
		pan.Logerror(err, true)
	}
	//parsing response for interface ips
	for _, e := range doc.FindElements("./response/result/ifnet/*") {
		name := strings.TrimSpace(e.SelectElement("name").Text())
		for _, ip := range e.SelectElements("ip") {
			if ip.Text() != "N/A" {
				ip_netmask := strings.Split(ip.Text(), "/")
				//if netmask is less than 24 we'll use /24 otherwise the firewall will complain
				if strings.Compare(ip_netmask[1], "24") < 0 {
					message = "Warning: Netmask lower than 24, we'll use /24 for GARPs make sure there isn't NATs outside this net. Interface: " + name
					println(message)
					pan.Wlog("output.txt", message, true)
					pan.Wlog("garp.txt", "test arp gratuitous interface "+name+" ip "+ip_netmask[0]+"/24", true)
				} else {
					//Generating file to send GARP commands using SSH
					pan.Wlog("garp.txt", "test arp gratuitous interface "+name+" ip "+ip.Text(), true)
				}
			}
		}
	}
	//Sending SSH commands
	println("sending GARPs using SSH please wait...")
	panssh.Send(fqdn, user, pass, "garp.txt", true, false)

}

//Function to check Global counters for know issues
func gcCheck(fqdn string, api string) {
	//Global counters and their possible causes
	counters := map[string]string{"flow_fwd_zonechange": "Assymetric routing",
		"flow_rcv_dot1q_tag_err": "Received traffic on a non configured vlan"}
	//Map to store the base global counter value during the first call.
	value := make(map[string]int)
	//Delay in seconds to calculate the rate
	delay := 5
	//Max rate to avoid false positives. Global counter of 1 or 5 lps could be meaningless
	maxrate := 20
	for counter := range counters {
		value[counter] = show.GlobalCounter(fqdn, api, counter)
	}
	time.Sleep(time.Duration(delay) * time.Second)
	for counter := range counters {
		tmp := show.GlobalCounter(fqdn, api, counter)
		//Calculating rate and comparing it with the max rate expected
		if ((tmp - value[counter]) / delay) > maxrate {
			message := "packet rate exceeding threshold for global counter " + counter + " this could be because of the following reason: " + counters[counter]
			println(message)
			pan.Wlog("output.txt", message, true)
		}
	}
}

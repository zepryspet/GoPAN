package cutover
import(
    "GoPAN/utils"
    "GoPAN/run/ssh"
    "github.com/beevik/etree"
    "net/url"
    "strings"
    "os"
    "regexp"
)

func Check (fqdn string, user string, pass string) {
    var message string
    //Generating api Key
    apikey := pan.Keygen(fqdn, user, pass)
    req, err := url.Parse("https://" + fqdn + "/api/?" )
	pan.Logerror(err, true)
    /////////////////////////
    //Check incomplete arps.
    //Generating HTTP req.
    cmd := (pan.CmdGen("show arp n_all"))
    q := url.Values{}
    q.Add("type", "op")
    q.Add("key", apikey)
    q.Add("cmd", cmd )
    req.RawQuery = q.Encode()
    //Sending the query to the firewall
    resp, err := pan.HttpValidate(req.String(), false)
    pan.Logerror(err, true)
    doc := etree.NewDocument()
    doc.ReadFromBytes(resp)
    for _, e := range doc.FindElements("./response/result/entries/*") {
        status := strings.TrimSpace(e.SelectElement("status").Text())
        if status == "i"{
            intf := e.SelectElement("interface").Text()
            println ("incomplete arp found on: " + intf)
        }  
    }
    /////////////////////////////////
    //Verify they interfaces aren't running either 100/10 Mbps or half duplex or are down.
    //Generating HTTP req.
    cmd = pan.CmdGen("show interface t_all")
    q.Set("cmd", cmd )
    req.RawQuery = q.Encode()
    //Sending the query to the firewall
    resp, err = pan.HttpValidate(req.String(), false)
    pan.Logerror(err, true)
    doc = etree.NewDocument()
    doc.ReadFromBytes(resp)
    //parsing response for hardware status
    for _ , e := range doc.FindElements("./response/result/hw/*") {
        name := strings.TrimSpace(e.SelectElement("name").Text())
        //Checking status, speed and state for all ethernet interfaces
        if strings.HasPrefix(name , "ethernet"){
            speed := strings.TrimSpace(e.SelectElement("speed").Text())
            state := strings.TrimSpace(e.SelectElement("state").Text())
            duplex := strings.TrimSpace(e.SelectElement("duplex").Text())
            //intf := e.SelectElement("interface").Text()
            if state != "up"{
                println ("interface " + name + " is DOWN!")
            }else{
                if duplex != "full"{
                    message = "interface " + name + " is half duplex!"
                    println (message) 
                    pan.Wlog ( "output.txt",message, true)
                }
                if (speed == "10" || speed == "100"){
                    message ="interface " + name + " speed is less than 1 Gbps. Current speed: " +speed
                    println (message) 
                    pan.Wlog ( "output.txt",message, true)
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
    for _ , e := range doc.FindElements("./response/result/ifnet/*") {
        name := strings.TrimSpace(e.SelectElement("name").Text())
        for _, ip := range e.SelectElements ("ip"){
            if ip.Text() != "N/A"{
                ip_netmask := strings.Split(ip.Text() , "/")
                //if netmask is less than 24 we'll use /24 otherwise the firewall will complain
                if strings.Compare (ip_netmask[1] , "24") < 0 {
                    message ="Warning: Netmask lower than 24, we'll use /24 for GARPs make sure there isn't NATs outside this net. Interface: " + name
                    println (message) 
                    pan.Wlog ( "output.txt",message, true)
                    pan.Wlog ( "garp.txt","test arp gratuitous interface "+ name + " ip " + ip_netmask[0] +"/24", true)
                }else{
                //Generating file to send GARP commands using SSH
                pan.Wlog ( "garp.txt","test arp gratuitous interface "+ name + " ip " + ip.Text(), true)
                }
            }
        }
    }
    //Sending SSH commands
    println ("sending GARPs using SSH please wait...")
    panssh.Send(fqdn, user, pass, "garp.txt" , true, false)
    
    ///////////////////
    //Check CRCs 
    //
    //Generating HTTP req.
    cmd = pan.CmdGen("show system state filter-pretty t_sys.*")
    q.Set("cmd", cmd )
    req.RawQuery = q.Encode()
    //Sending the query to the firewall
    resp, err = pan.HttpValidate(req.String(), false)
    pan.Logerror(err, true)
    r := regexp.MustCompile(`sys.(s\d.p\d).detail:\s{\s[^}]*crc`)
    ports := r.FindAllStringSubmatch(string(resp), -1)
    for i := 0; i < len(ports); i++ {
        index := strings.Split(ports[i][1], ".")
        slot := strings.TrimPrefix(index[0], "s")
        port := strings.TrimPrefix(index[1], "p")
        message = "CRCs detected in ethernet" + slot + "/" + port
        println(message)
        pan.Wlog ( "output.txt",message, true)
    }
    //parsing response for hardware status
    //Regex check the first capturing group ''
    
}
    
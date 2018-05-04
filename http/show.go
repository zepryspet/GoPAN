package show

//package with common show commands

import(
    "github.com/zepryspet/GoPAN/utils"
    "net/url"
    "github.com/beevik/etree"
    "strings"
    "strconv"
)



//Show global counter and return the current value, to calculate rate send 2 of them in a second and calculate diference
func GlobalCounter(fqdn string, api string, counter string) int {
    //Global variables
    q := url.Values{}
    q.Add("type", "op")
    q.Add("key", api)
    //Building HTTP request
    req, err := url.Parse("https://" + fqdn + "/api/?" )
	pan.Logerror(err, true)
    counter = "t_" + counter
    cmd := pan.CmdGen("show counter global name " + counter)
    q.Add("cmd", cmd )
    req.RawQuery = q.Encode()
    //Sending the query to the firewall
    resp, err := pan.HttpValidate(req.String(), false)
    pan.Logerror(err, true)
    doc := etree.NewDocument()
    doc.ReadFromBytes(resp)
    number, err := strconv.Atoi(strings.TrimSpace(doc.FindElement("./response/result/global/counters/entry/value").Text()))
    pan.Logerror(err, true)
    return  (number)
}

package loadconfig

import (
	"fmt"
	"github.com/beevik/etree"
	pan "github.com/zepryspet/GoPAN/utils"
	"net/url"
	"os"
)

func Load(fqdn string, user string, pass string, fn string) {
	//Generating api Key
	apikey := pan.Keygen(fqdn, user, pass)
	req, err := url.Parse("https://" + fqdn + "/api/?")
	pan.Logerror(err, true)

	q := url.Values{}
	q.Add("type", "import")
	q.Add("category", "configuration")
	q.Add("key", apikey)
	req.RawQuery = q.Encode()
	//Sending the query to the firewall
	resp, err := pan.HttpPostFile(req.String(), fn)
	os.Exit(0)
	pan.Logerror(err, true)
	doc := etree.NewDocument()
	doc.ReadFromBytes(resp)
	e := doc.FindElement("./response/result/msg")
	fmt.Print("%v\n", e.Text())
}

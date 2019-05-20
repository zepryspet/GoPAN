package loadconfig

import (
	"fmt"
	"github.com/beevik/etree"
	pan "github.com/zepryspet/GoPAN/utils"
	"net/url"
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
	resp, fileUploadName, err := pan.HttpPostFile(req.String(), fn)
	pan.Logerror(err, true)
	doc := etree.NewDocument()
	doc.ReadFromBytes(resp)
	e := doc.FindElement("./response/msg/line")
	fmt.Printf("%v\n", e.Text())

	cmd := fmt.Sprintf("<load><config><from>%v</from></config></load>", fileUploadName)
	q = url.Values{}
	q.Add("type", "op")
	q.Add("key", apikey)
	q.Add("cmd", cmd)
	req.RawQuery = q.Encode()
	//Sending the query to the firewall
	resp, err = pan.HttpValidate(req.String(), false)
	pan.Logerror(err, true)
	doc = etree.NewDocument()
	doc.ReadFromBytes(resp)
	e = doc.FindElement("./response/result/msg/line")
	fmt.Printf("%v\n", e.Text())

}

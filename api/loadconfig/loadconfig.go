package loadconfig

import (
	"fmt"
	"github.com/beevik/etree"
	pan "github.com/zepryspet/GoPAN/utils"
	"net/url"
	"time"
)

/*
Load a configuration into a firewall, without committing.
fqdn: Firewall fqdn
user: username
pass: password
fn: Path to local file to upload.
*/
func Load(fqdn string, user string, pass string, fn string) {
	//Generating api Key
	apikey := pan.Keygen(fqdn, user, pass)
	req, err := url.Parse("https://" + fqdn + "/api/?")
	pan.Logerror(err, true)

	fmt.Printf("Importing configuration from local file %v...\n", fn)
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

	fmt.Print("Loading configuration into candidate...\n")
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

	fmt.Print("Committing configuration...\n")
	Commit(fqdn, apikey)
}

/*
Commit the current candidate configuration
*/
func Commit(fqdn string, apikey string) {
	req, err := url.Parse("https://" + fqdn + "/api/?")
	cmd := (pan.CmdGen("commit"))
	q := url.Values{}
	q.Add("type", "commit")
	q.Add("key", apikey)
	q.Add("cmd", cmd)
	req.RawQuery = q.Encode()
	//Sending the query to the firewall
	resp, err := pan.HttpValidate(req.String(), false)
	pan.Logerror(err, true)
	doc := etree.NewDocument()
	doc.ReadFromBytes(resp)

	e := doc.FindElement("./response/result/msg/line")
	fmt.Printf("%v\n", e.Text())
	jobid := doc.FindElement("./response/result/job").Text()

	cmd = fmt.Sprintf("<show><jobs><id>%v</id></jobs></show>", jobid)
	q = url.Values{}
	q.Add("type", "op")
	q.Add("key", apikey)
	q.Add("cmd", cmd)
	req.RawQuery = q.Encode()

	resp, err = pan.HttpValidate(req.String(), false)
	pan.Logerror(err, true)
	doc = etree.NewDocument()
	doc.ReadFromBytes(resp)

	status := doc.FindElement("./response/result/job/status").Text()
	for status == "ACT" {
		resp, err = pan.HttpValidate(req.String(), false)
		pan.Logerror(err, true)
		doc = etree.NewDocument()
		doc.ReadFromBytes(resp)

		status = doc.FindElement("./response/result/job/status").Text()
		progress := doc.FindElement("./response/result/job/progress").Text()
		fmt.Printf("Commit status: %v (%v)\n", status, progress)
		time.Sleep(2 * time.Second)
	}
}

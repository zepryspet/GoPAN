package urlcat

import (
    "GoPAN/utils"
    "net/url"
    "github.com/beevik/etree"
)

func Request (fqdn string, apiKey string, urls string, isFile bool) {
    //Creating API url
	req, err := url.Parse("https://" + fqdn + "/api/?")
	if err != nil {
		pan.Logerror(err, true)
	}
    //Creating cmd section
    root := etree.NewDocument()
    //Creating command "test url-info-cloud <url>"
    test := root.CreateElement("test")
    urlinfo := test.CreateElement("url-info-cloud")
    urlinfo.CreateCharData(urls)
    //cmd := root.WriteToString()
    
    command, err := root.WriteToString()
    if err != nil {
		pan.Logerror(err, true)
	} 
    println ( "can you see it?" + command)
    
	q := url.Values{}
	q.Add("type", "op")
    q.Add("key", apiKey)
    q.Add("cmd", command )
	req.RawQuery = q.Encode()

    resp, err := pan.HttpValidate(req.String(), true)
    if err != nil {
		pan.Logerror(err, true)
	} 
    doc := etree.NewDocument()
    doc.ReadFromBytes(resp)
    for _, e := range doc.FindElements("./response/result/*") {
        apiKey = e.Text()
    }
    
}
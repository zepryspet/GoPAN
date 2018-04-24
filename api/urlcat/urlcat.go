package urlcat

import (
    "GoPAN/utils"
    "net/url"
    "github.com/beevik/etree"
    "strings"
    "bufio"
    "os"
)

func Request (fqdn string, apiKey string, urls string, isFile bool) {
    var category string
    //Creating API url
	fqdn = "https://" + fqdn + "/api/?"
    //encoding url options
    q := url.Values{}
    q.Add("type", "op")
    q.Add("key", apiKey)
    
    //if it's a file read the file and send the commands line by line
    if isFile{
        file, err := os.Open(urls)
        if err != nil {
            pan.Logerror(err, true)
        }
        defer file.Close()

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            category = rcategory (fqdn, q, scanner.Text())
            //Writing to a csv file. creates one if it doesn't exist, appends if it does.
            pan.Wlog("categories.csv", strings.TrimSpace(scanner.Text()) + "," + category, true)
        }
        if err := scanner.Err(); err != nil {
            pan.Logerror(err, true)
        }
        
    } else{
        category = rcategory (fqdn, q, urls)
        //Writing to a csv file. creates one if it doesn't exist, appends if it does.
        pan.Wlog("categories.csv", urls + "," + category, true)
        //Printing the category in case there's only one object
        println(category)
    }
}

func rcategory (fqdn string, queries url.Values, website string) string{
    req, err := url.Parse(fqdn )
	if err != nil {
		pan.Logerror(err, true)
	}
    //Creating cmd section
    root := etree.NewDocument()
    //Creating command "test url-info-cloud <url>"
    test := root.CreateElement("test")
    urlinfo := test.CreateElement("url-info-cloud")
    urlinfo.CreateCharData(website)
    command, err := root.WriteToString()
    if err != nil {
        pan.Logerror(err, true)
    } 
    queries.Add("cmd", command )
    req.RawQuery = queries.Encode()
        //Sending the query to the firewall and validates the response.
    resp, err := pan.HttpValidate(req.String(), false)
    if err != nil {
        pan.Logerror(err, true)
    } 
    doc := etree.NewDocument()
    doc.ReadFromBytes(resp)
    //splitting response using spaces
    categories := strings.Fields(doc.FindElement("./response/result").Text()) 
    //the best match is located at 2nd postion just after "BM:" and the category in the 4th place after splitting by commas
    category := strings.Split(categories[1], ",")[3]
    //removing previous cmd from the queries
    queries.Del("cmd")
    return category
}
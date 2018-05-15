package threat

import (
    "net/url"
    "fmt"
    "github.com/zepryspet/GoPAN/utils"
    "github.com/beevik/etree"
    "github.com/360EntSecGroup-Skylar/excelize"
    "strings"
    "strconv"
)

func Export (fqdn string, user string, pass string){
    
    //Generating api Key
    apikey := pan.Keygen(fqdn, user, pass)
    req, err := url.Parse("https://" + fqdn + "/api/?" )
	pan.Logerror(err, true)
    spyware, vulnerability := threat(apikey, req)
    
    //Generating new excel file
    xlsx := excelize.NewFile()
    // Create a new sheet and deleting old
    index := xlsx.NewSheet("threats")
    xlsx.DeleteSheet("Sheet1")
    // Setting table headers.
    xlsx.SetCellValue("threats", "A1", "Threat ID")
    xlsx.SetCellValue("threats", "B1", "Threat Type")
    xlsx.SetColWidth("threats", "A", "B", 12)
    xlsx.SetCellValue("threats", "C1", "Threat Name")
    xlsx.SetCellValue("threats", "D1", "Threat Description")
    xlsx.SetColWidth("threats", "C", "D", 40)
    xlsx.SetCellValue("threats", "E1", "Threat Severity")
    xlsx.SetCellValue("threats", "F1", "CVE")
    xlsx.SetColWidth("threats", "E", "F", 20)
    
    //Checking spyware IDs
    i := 2  //Integer to keep track of the rows
    for k, v := range spyware {
        fmt.Printf("%d -> %s\n", k, v)
        values := description(apikey, req, k)
        // Writing cell values
        c := strconv.Itoa(i) 
        xlsx.SetCellValue("threats", "A"+c , strconv.Itoa(k))
        xlsx.SetCellValue("threats", "B"+c , "spyware")
        xlsx.SetCellValue("threats", "C"+c , v)
        xlsx.SetCellValue("threats", "D"+c , values[0])
        xlsx.SetCellValue("threats", "E"+c , values[1])
        xlsx.SetCellValue("threats", "F"+c , values[2])
        i += 1
    }
    //Checking vulnerability IDs
    for k, v := range vulnerability {
        fmt.Printf("%d -> %s\n", k, v)
        values := description(apikey, req, k)
        // Writing cell values
        c := strconv.Itoa(i) 
        xlsx.SetCellValue("threats", "A"+c , strconv.Itoa(k))
        xlsx.SetCellValue("threats", "B"+c , "vulnerability")
        xlsx.SetCellValue("threats", "C"+c , v)
        xlsx.SetCellValue("threats", "D"+c , values[0])
        xlsx.SetCellValue("threats", "E"+c , values[1])
        xlsx.SetCellValue("threats", "F"+c , values[2])
        i += 1
    }
    //Adding a table
    xlsx.AddTable("threats", "A1", "F" +strconv.Itoa(i-1), `{"table_name":"DB","table_style":"TableStyleMedium2", "show_first_column":true,"show_last_column":true,"show_row_stripes":false,"show_column_stripes":true}`)
    // Set active sheet of the workbook.
    xlsx.SetActiveSheet(index)
    // Saving xlsx file.
    err = xlsx.SaveAs("./Threats.xlsx")
    pan.Logerror(err, true)
}

//function that returns all predefined ids with its human readable name. It returns both spyware and vulnrability maps
func threat(apikey string, req *url.URL) (map[int]string, map[int]string ){
    spyware := make(map[int]string)
    vulnerability := make(map[int]string)
    //Generating HTTP req.
    q := url.Values{}
    q.Add("type", "config")
    q.Add("action", "get")
    q.Add("key", apikey)
    q.Add("xpath", "/config/predefined/threats" )
    req.RawQuery = q.Encode()
    //Sending the query to the firewall
    resp, err := pan.HttpValidate(req.String(), false)
    pan.Logerror(err, true)
    //fmt.Printf("%s", resp)
    doc := etree.NewDocument()
    doc.ReadFromBytes(resp)
    //Spyware
    for _, e := range doc.FindElements("./response/result/threats/phone-home/*") {
        //Converting threat id to number to create a dictionary
        id, err := strconv.Atoi(strings.TrimSpace(e.SelectAttr("name").Value))
        pan.Logerror(err, true)
        spyware[id] = e.SelectElement("threatname").Text()
    }
    //vulnerability
    for _, v := range doc.FindElements("./response/result/threats/vulnerability/*") {
        //Converting threat id to number to create a dictionary
        id, err := strconv.Atoi(strings.TrimSpace(v.SelectAttr("name").Value))
        pan.Logerror(err, true)
        vulnerability[id] = v.SelectElement("threatname").Text()
    }
    
    return spyware, vulnerability
}


func description(apikey string, req *url.URL, paloid int) []string{
    //Changiong int id to text and prepending the "t_" to indicate that's text inside the op command
    var val = []string{}
    id:= strconv.Itoa(paloid)
    id = "t_" + id
    //Generating HTTP req.
    cmd := pan.CmdGen("show threat id " + id)
    q := url.Values{}
    q.Add("type", "op")
    q.Add("key", apikey)
    q.Add("cmd", cmd )
    req.RawQuery = q.Encode()
    //Sending the query to the firewall
    resp, err := pan.HttpValidate(req.String(), false)
    pan.Logerror(err, true)
    //fmt.Printf("%s", resp)
    doc := etree.NewDocument()
    doc.ReadFromBytes(resp)
    entry := doc.FindElement("./response/result/entry")
    //Looking for threat description
    if description := entry.SelectElement("description"); description != nil {
        val = append(val, strings.TrimSpace(description.Text()))
    }else{
        val = append(val, "")
    }
    //Looking for threat severity
    if severity := entry.SelectElement("severity"); severity != nil {
        val = append(val,severity.Text())
    }else{
        
    }
    //Looking for vulnerability
    if vul := entry.SelectElement("vulnerability"); vul != nil {
        if e := vul.SelectElement("cve"); e != nil{
            val = append(val,e.SelectElement("member").Text())
        }else{ val = append(val, "")}
    }else{val = append(val, "")}
    
    return val
}
    
//func cve(apikey string, req *url.URL){
    
//}
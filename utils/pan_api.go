package pan

import (
	"crypto/tls"
	"fmt"
	"github.com/beevik/etree"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

//Function to generate
func Keygen(user string, pass string, fqdn string) {
	//Defining main variables
	//user := "admin"
	//pass := "admin"
	//fqdn := "192.168.55.135"

	//Defyning secondary variables
	req, err := url.Parse("https://" + fqdn + "/api/?")
	if err != nil {
		log.Fatal(err)
	}
	q := url.Values{}
	q.Add("password", pass)
	q.Add("user", user)
	q.Add("type", "keygen")
	req.RawQuery = q.Encode()
	fmt.Println(req.String())

	//Ignoring TLS certificate checking
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	//Setting HTTP timeout as 10 seconds.
	netClient := &http.Client{
		Timeout:   time.Second * 10,
		Transport: tr,
	}

	resp, err := netClient.Get(req.String())
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//making sure the API responds with a 200 code
	if resp.StatusCode == 200 {
		fmt.Println(string(body))
		//creating new XML object and parsing the response
		doc := etree.NewDocument()
		doc.ReadFromBytes(body)
		for _, e := range doc.FindElements("./response/result/*") {
			fmt.Println(e.Text())
		}
	} else {
		fmt.Println("error with HTTP request:\t" + req.String())
	}
}

func Wlog(fileName string, text string, newline bool) {
	// If the file doesn't exist, create it, or append to the file
	if newline {
		text += "\n"
	}
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte(text)); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

//Prints the error and exist
func Logerror(err error) {
	if err != nil {
		log.Fatal(err)
		Wlog("error.txt", err.Error(), true)
		os.Exit(1)
	}
}

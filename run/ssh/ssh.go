package panssh

import (
	//"swisspan/cps"
	//"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"strings"
	"GoPAN/utils"
	"time"
)

func Send(fqdn string, user string, pass string) {
	//Defining main variables
	//user := "admin"
	//pass := "admin"
	//pan.Keygen(user, pass, fqdn)
	//fqdn := "172.16.1.1"
	//community := "public"
	//minutes := 1
	//cps.Snmpgen(fqdn, community, minutes)
	//user := "admin"
	//pass := "admin"
	// Start up ssh process
	sshClt, err := ssh.Dial("tcp", fqdn+":22", &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(pass)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	pan.Logerror(err, true)
	session, err := sshClt.NewSession()
	pan.Logerror(err, true)
	sshOut, err := session.StdoutPipe()
	pan.Logerror(err, true)
	sshIn, err := session.StdinPipe()
	pan.Logerror(err, true)

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	// Request pseudo terminal
	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
		log.Fatal("request for pseudo terminal failed: ", err)
	}
	// Start remote shell
	if err := session.Shell(); err != nil {
		log.Fatal("failed to start shell: ", err)
	}
	//wait for banner
	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.

	cmdSend(sshOut, sshIn, "show system info", false, false, 20)
	session.Close()
}

func cmdSend(sshOut io.Reader, sshIn io.WriteCloser, cmd string, isfile bool, isconfig bool, timeout int) {
	readBuff(">", sshOut, timeout)
	if _, err := writeBuff("set cli pager off", sshIn); err != nil {
		log.Fatal("Failed to run: %s", err)
	}
	readBuff(">", sshOut, timeout)
	if _, err := writeBuff(cmd, sshIn); err != nil {
		log.Fatal("Failed to run: %s", err)
	}
	readBuff(">", sshOut, timeout)
	if _, err := writeBuff("exit", sshIn); err != nil {
		log.Fatal("Failed to run: %s", err)
	}
}

func readBuffForString(whattoexpect string, sshOut io.Reader, buffRead chan<- string) {
	buf := make([]byte, 2000)
	n, err := sshOut.Read(buf) //this reads the ssh terminal
	waitingString := ""
	if err == nil {
		waitingString = string(buf[:n])
	}
	for (err == nil) && (!strings.Contains(waitingString, whattoexpect)) {
		n, err = sshOut.Read(buf)
		waitingString += string(buf[:n])
		//fmt.Println(waitingString) //uncommenting this might help you debug if you are coming into errors with timeouts when correct details entered
	}
	fmt.Println(waitingString)
	pan.Wlog("output.txt", waitingString, true)
	buffRead <- waitingString
}

func readBuff(whattoexpect string, sshOut io.Reader, timeoutSeconds int) string {
	ch := make(chan string)
	go func(whattoexpect string, sshOut io.Reader) {
		buffRead := make(chan string)
		go readBuffForString(whattoexpect, sshOut, buffRead)
		select {
		case ret := <-buffRead:
			ch <- ret
		case <-time.After(time.Duration(timeoutSeconds) * time.Second):
			pan.Wlog("error.txt", "timeout waiting for command", true)
			break
		}
	}(whattoexpect, sshOut)
	return <-ch
}

func writeBuff(command string, sshIn io.WriteCloser) (int, error) {
	returnCode, err := sshIn.Write([]byte(command + "\r"))
	return returnCode, err
}
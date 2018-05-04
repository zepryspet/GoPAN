package panssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"strings"
	"github.com/zepryspet/GoPAN/utils"
	"time"
    "bufio"
    "os"
)

func Send(fqdn string, user string, pass string, command string, isFile bool, isConfig bool) {

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
        pan.Wlog ("request for pseudo terminal failed: ", "error.txt", false)
		pan.Logerror(err, true)
	}
	// Start remote shell
	if err := session.Shell(); err != nil {
        pan.Wlog ("request for remote shell failed: ", "error.txt", false)
		pan.Logerror(err, true)
	}
	//wait for banner
	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.

	cmdSend(sshOut, sshIn, command, isFile, isConfig, 20)
	session.Close()
}

func cmdSend(sshOut io.Reader, sshIn io.WriteCloser, cmd string, isFile bool, isConfig bool, timeout int) {
    //setting up the initial prompt
    prompt := ">"
	readBuff(prompt, sshOut, timeout)
    //disabling the CLI pager to avoid having to tab on large outputs
	if _, err := writeBuff("set cli pager off", sshIn); err != nil {
		pan.Logerror(err, true)
	}
	readBuff(prompt, sshOut, timeout)
    
    //verifying if the commands need to be run in config mode
    if isConfig{
        //changing the prompt to bash as due config mode
        prompt = "#"
        //Sending configuration
        if _, err := writeBuff("configure", sshIn); err != nil {
            pan.Logerror(err, true)
        }
        readBuff(prompt, sshOut, timeout)
    }
    
    //Sending the command (s) to the endpoints
    
    //checking if it's a file or a single command
    if isFile{
        file, err := os.Open(cmd)
        if err != nil {
            pan.Logerror(err, true)
        }
        defer file.Close()

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            //removing empty spaces
            newline := strings.TrimSpace(scanner.Text())
            if _, err := writeBuff(newline, sshIn); err != nil {
                pan.Logerror(err, true)
            }
            readBuff(prompt, sshOut, timeout)
        }
        if err := scanner.Err(); err != nil {
            pan.Logerror(err, true)
        }
    } else{
        if _, err := writeBuff(cmd, sshIn); err != nil {
            pan.Logerror(err, true)
        }
        readBuff(prompt, sshOut, timeout)
    }
    //Exiting from config mode
    if isConfig{
        if _, err := writeBuff("exit", sshIn); err != nil {
		pan.Logerror(err, true)
	   }
    }
    //Exiting from operational mode
	if _, err := writeBuff("exit", sshIn); err != nil {
		pan.Logerror(err, true)
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
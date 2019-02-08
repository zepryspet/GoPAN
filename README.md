Non-oficcial script set to deal with Palo Alto networks products.

___
I(a). Compilying the binary. skip this step if you want to download the precompiled binary files.

1. Install go for your system.
>https://golang.org/doc/install

2. Download the script source code
> go get github.com/zepryspet/GoPAN

3. Change directory to the dowloaded folder
> cd $HOME/go/src/github.com/zepryspet/GoPAN

4. Generate the binary file (this option is not needed since you can also run the script without compiling it). This will generate the binary in the current directory
> go build

5. run the script
>./GoPAN
___

I(b). Downloading the precompiled binary scripts:

For MAC

>curl -L https://github.com/zepryspet/GoPAN/raw/master/binaries/mac/GoPAN -o GoPAN

>chmod 755 GoPAN

>./GoPAN -h


For Windows x32

> dowload https://github.com/zepryspet/GoPAN/raw/master/binaries/windows32bit/GoPAN.exe

> open cmd in the dowload directory

>GoPAN.exe -h


For Windows x64

> dowload https://github.com/zepryspet/GoPAN/raw/master/binaries/windows64bit/GoPAN.exe

> open cmd in the dowload directory

>GoPAN.exe -h


For Linux x64 (tested on ubuntu 16)

>curl -L https://github.com/zepryspet/GoPAN/raw/master/binaries/linux64bit/GoPAN -o GoPAN

>chmod 755 GoPAN

>./GoPAN -h


___
II. Running the scripts

1. Change the directory where the binary file is stored.
> cd $HOME/go/src/github.com/zepryspet/GoPAN

2. Run the script to get help

> ./GoPAN


    >Usage:
      pan [command]

    >Available Commands:
      api         scripts using the api calls
      help        Help about any command
      run         Pre-built scripts to collect and process firewall data using non-api methods like SNMP or SSH


    > Flags:
      -h, --help   help for pan

___

III. Current script structure:

run:

    ssh

    cps

api:

    urlcat

    cutover

    keygen

___
IV. Script details:

1. cps.

    Getting the CPS per zone in palo alto firewall using SNMP (requieres PAN-OS 8.0+) the csv files will be saved in the working directory. You usually run it for a week or 2 and then open the csv to get the max CPS per zone and set the alert to 1.2x and the block rate to 1.5x

    Example using SNMPv2
    > ./GoPAN run cps -c <snmp-community> -i <firewall-ip>

    Example using SNMPv3
    >./GoPAN run cps -i <firewall-ip> -c <username> -v 3 -a <auth-password> -x <privacy-password>

    Note:Default polling time is 10 seconds (SNMP statistics are updated every 10 sec). use "-s" flag to change it if needed. Keep in mind that file rotation is not implemented so keep an eye on disk space.

2. ssh.

    Sent a batch of commands to a firewall or panorama using SSH. The input accepts either a filename in the current working directory containing a set of commands or a single command. It can send either configuration commands or operational commands (default) but no both. The commands can be run once or every X minutes.

    Examples:

    running "show clock" once on the endpoint every 10 seconds (-t 10)
    >./GoPAN run ssh -i firewall-ip -p password -u username -r 'show clock' -t 10

    running "run show system info" on the endpoint in config mode(used run because the firewall will be in config mode but could be any configuration command)
    >./GoPAN run ssh -i firewall-ip -p password -u username -r 'run show system info' -c

    Running commands from "commands.txt" (-f) in configuration mode (-c)
    >./GoPAN run ssh -i firewall-ip -p password -u username -r commands.txt -fc


3. urlcat.

    Script that request the url category for a single website or for multiple websites stored in a clear text file, the output will be saved in a csv file within the same folder named "categories.csv"

    Examples:

    requesting a category for a single website

    >./GoPAN api urlcat -i firewall-ip -p password -u username -w www.facebook.com

    requesting categories for urls inside a text file

    >./GoPAN api urlcat -i firewall-ip -p password -u username -w websites.txt -f

4. cutover

    Useful scripts for maintenance windows. It checks: incomplete mac addresses, interface status, speed and duplex information and send gratuitous ARPs on all interfaces.

    >./GoPAN api cutover -i firewall-ip -p password -u username

5. Keygen

    Generates and prints an API key from a palo alto firewall

    >./GoPAN api keygen -i firewall-ip -p password -u username

6. Threat

    Exports the firewall threat database from a firewall into excel. Including threat ID, name, description, type, severity and CVE.

    >./GoPAN api threat -i firewall-ip -p password -u username

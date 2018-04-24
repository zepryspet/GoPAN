Non-oficcial script set to deal with Palo Alto networks products.

___
Installation steps (in case you don't want to dowload the precompiled binaries or your system isn't MAC or windowsx64):

1. Install go for your system.
>https://golang.org/doc/install

2. Download the script 
> go get github.com/zepryspet/GoPAN

3. Change directory to the dowloaded folder
> cd $HOME/go/src/github.com/zepryspet/GoPAN

4. Generate the binary files (this option is not needed since you can also run the script without compiling it). This will generate the binary in the current directory
> go build pan.go


Running the scripts

1. Change the directory where the binary file is stored.
> cd $HOME/go/src/github.com/zepryspet/GoPAN

2. Run the script to get help

> ./pan
    
    
    >Usage:
      pan [command]

    >Available Commands:
      help        Help about any command
      run         Pre-built scripts to collect and process firewall data

    > Flags:
      -h, --help   help for pan

___

current script structure:

run:

    ssh

    cps

api:

    urlcat
    
___
Available scripts:

1. cps. 

    Getting the CPS per zone in palo alto firewall using SNMP (requieres PAN-OS 8.0+) the csv files will be saved in the working directory. You usually run it for a week or 2 and then open the csv to get the max CPS per zone and set the alert to 1.2x and the block rate to 1.5x

    Example
    > ./pan run cps -c snmp-community -i firewall-ip

2. ssh.

    Sent a batch of commands to a firewall or panorama using SSH. The input accepts either a filename in the current working directory containing a set of commands or a single command. It can send either configuration commands or operational command but no both. The commands can be once or every X minutes.

    Example:
    
3. urlcat.

    Script that request the url category for a single website or for multiple websites stored in a clear text file, the output will be saved in a csv file within the same folder named "categories.csv"

    Examples:
    >requesting a category for a single website
    
    >./pan api urlcat -i firewall-ip -p password -u username -w www.facebook.com 

    >requesting categories for urls inside a text file
    
    >./pan api urlcat -i firewall-ip -p password -u username -w websites.txt -f

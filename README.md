# nonstd
Go Language utility to encode/decode ip addresses for the Twister Network

This is a toool I created for the dns seeder for the [Twister P2P network](https://github.com/gombadi/dnsseeder/)

This tool allows you to test a DNS server and check it is returning the correct information.



> **NOTE:** This repository is under ongoing development and
is likely to break over time. Use at your own risk.


## Installing

Simply use go get to download the code:

    $ go get github.com/gombadi/nonstd


## Usage

    $ nonstd -h <dns server to test>

    $ nonstd -i ipaddress -p port // to see an encoded ip for a real ip & port

    $ nonstd -t1 ipaddress -t2 ipaddress // check if these two ip addresses are associated and if so display the port


### Command line Options:

```


Usage of nonstd:
  -h string
        DNS host to query for addresses
  -i string
        Real ip to encode
  -p string
        Port encode
  -t1 string
        Single ip to test. Also needs -t2
  -t2 string
        Single ip to test. Also needs -t1
  -v    Display additional information

```


/*

 */
package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var host, port, ip, t1, t2 string
var verbose bool

func main() {

	flag.StringVar(&t1, "t1", "", "Single ip to test. Also needs -t2")
	flag.StringVar(&t2, "t2", "", "Single ip to test. Also needs -t1")
	flag.StringVar(&host, "h", "", "DNS host to query for addresses")
	flag.StringVar(&port, "p", "", "Port encode")
	flag.StringVar(&ip, "i", "", "Real ip to encode")
	flag.BoolVar(&verbose, "v", false, "Display additional information")
	flag.Parse()

	if host != "" && ip != "" {
		fmt.Printf("Please only select a host with -h OR an ip address with -i\n")
	}

	if t1 != "" && t2 == "" {
		fmt.Printf("Please enter both -t1 and -t2 at the same time\n")
	}

	if t2 != "" && t1 == "" {
		fmt.Printf("Please enter both -t1 and -t2 at the same time\n")
	}

	// convert an ip & port to real & encoded addresses
	if ip != "" {
		encodeip(ip, port)
	}

	// pull ip addresses from dns server and decode them
	if host != "" {
		var lookuphost string
		if x := strings.HasPrefix(host, "nonstd."); x == false {
			lookuphost = "nonstd." + host
		} else {
			lookuphost = host
		}
		ips, err := net.LookupHost(lookuphost)
		if err != nil {
			fmt.Printf("Error - unable to lookup addresses from %s %v\n", host, err)
		}
		decodeips(ips)
	}

	if t1 != "" && t2 != "" {
		ips := []string{t1, t2}
		decodeips(ips)
	}

	fmt.Printf("\nProgram exiting. Bye\n")
}

// decideips receives a slice of ip addresses and displays information
// about the ip addresses and ports encoded in them
func decodeips(ips []string) {

	var match bool

	if verbose && host != "" {
		fmt.Printf("Received %v addresses from %s\n", len(ips), host)
	}

	for _, ip := range ips {

		netip := net.ParseIP(ip)
		netip = netip.To4()
		if netip == nil {
			continue
		}
		crcAddr := crc16(netip)
		match = false

		for _, ipport := range ips {

			netipport := net.ParseIP(ipport)
			netipport = netipport.To4()
			if netipport == nil {
				continue
			}
			if (netipport[0] == byte(crcAddr>>8)) && (netipport[1] == byte(crcAddr&0xff)) {
				theport := (uint16(netipport[2]) << 8) + uint16(netipport[3])
				fmt.Printf("realip: %s\tport: %v\tencodedip: %s\n", ip, theport, ipport)
				if verbose {
					fmt.Printf("crcAddr: 0x%x encoded bytes: 0x%x 0x%x 0x%x 0x%x\n", crcAddr, netipport[0], netipport[1], netipport[2], netipport[3])
				}
				match = true
			}
		}
		if match == false && verbose {
			fmt.Printf("No match found for ip: %s\n", ip)
		}
	}

}

// encodeip takes an ip address and port strings and displays the encoded
// ip address they would use
func encodeip(ip, port string) {

	var netport int
	if port != "" {
		netport, _ = strconv.Atoi(port)
	}

	netip := net.ParseIP(ip)
	netip = netip.To4()
	if netip == nil {
		fmt.Printf("error changing ip: %s to 4 bytes storage\n", ip)
		os.Exit(0)
	}

	crcAddr := crc16(netip)

	bs := make([]byte, 4)
	bs[0] = byte(crcAddr >> 8)
	bs[1] = byte(crcAddr & 0xff)
	bs[2] = byte(netport >> 8)
	bs[3] = byte(netport & 0xff)

	encodedip := net.IPv4(bs[0], bs[1], bs[2], bs[3])
	if x := encodedip.To4(); x == nil {
		fmt.Printf("Error checking encoded ip to real ip\n")
		fmt.Printf("Real ip:   \t%s\n", ip)
		fmt.Printf("Encoded ip:\t%v.%v.%v.%v\n", bs[0], bs[1], bs[2], bs[3])
	} else {

		fmt.Printf("crcAddr:\t%#v\n", crcAddr)
		fmt.Printf("Real ip:   \t%s\n", ip)
		fmt.Printf("Encoded ip:\t%s\n", encodedip.String())
	}
}

func crc16(bs []byte) uint16 {
	var x, crc uint16
	crc = 0xffff

	for _, v := range bs {

		x = crc>>8 ^ uint16(v)
		x ^= x >> 4
		crc = (crc << 8) ^ (x << 12) ^ (x << 5) ^ x
	}
	return crc
}

/*

 */

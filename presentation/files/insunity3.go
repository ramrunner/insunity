package main

import (
	"fmt"
	"log"
	"log/syslog"
	"net"
	"strings"
)

// START OMIT
func getIPv4ForInterfaceName(ifname string) net.IP {
	interfaces, _ := net.Interfaces()
	for _, inter := range interfaces {
		if inter.Name == ifname {
			if addrs, err := inter.Addrs(); err == nil {
				for _, addr := range addrs {
					switch ip := addr.(type) {
					case *net.IPNet:
						if ip.IP.DefaultMask() != nil {
							return ip.IP
						}
					}
				}
			}
		}
	}
	return nil
}

// END OMIT

func genIPs(ipstr string, ch *chan string, port int) {
	lastdotind := strings.LastIndex(ipstr, ".")
	if lastdotind == -1 {
		log.Printf("malformed IP")
	} else {
		prestr := ipstr[:lastdotind]
		var ipport string
		for i := 1; i < 255; i++ {
			ipport = fmt.Sprintf("%s.%d:%d", prestr, i, port)
			*ch <- ipport
		}
	}
	close(*ch)
}

func main() {
	//set log to syslog
	logwriter, e := syslog.New(syslog.LOG_NOTICE, "insunity")
	if e == nil {
		log.SetOutput(logwriter)
	}
	//figure out the IPs for all the interfaces
	ip := getIPv4ForInterfaceName("eth0")
	if ip == nil {
		log.Fatal("could not get IP address for eth0")
	}
	targetch := make(chan string)
	go genIPs(ip.String(), &targetch, 1337)
	for t := range targetch {
		log.Printf("will scan:%s\n", t)
	}
}

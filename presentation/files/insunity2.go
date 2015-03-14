package main

import (
	"log"
	"log/syslog"
	"net"
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

func main() {
	//set log to syslog
	logwriter, e := syslog.New(syslog.LOG_NOTICE, "insunity")
	if e == nil {
		log.SetOutput(logwriter)
	}
	//figure out the IPs for all the interfaces
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	for _, iface := range ifaces {
		log.Printf("[%v] - %v\n", iface.Name, getIPv4ForInterfaceName(iface.Name).String())
	}
}

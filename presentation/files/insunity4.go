package main

import (
	"fmt"
	"log"
	"log/syslog"
	"net"
	"strings"
	//"time"
	service "github.com/chai2010/protorpc/internal/service.pb"
	"github.com/chai2010/protorpc/proto"
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

type Echo int

func (t *Echo) Echo(args *service.EchoRequest, reply *service.EchoResponse) error {
	reply.Msg = proto.String(args.GetMsg())
	return nil
}

func (t *Echo) EchoTwice(args *service.EchoRequest, reply *service.EchoResponse) error {
	reply.Msg = proto.String(args.GetMsg() + args.GetMsg())
	return nil
}

func init() {
	go service.ListenAndServeEchoService("tcp", `127.0.0.1:9527`, new(Echo))
}

func main() {
	//set log to syslog
	logwriter, e := syslog.New(syslog.LOG_NOTICE, "insunity")
	if e == nil {
		log.SetOutput(logwriter)
	}
	/*
		//figure out the IPs for all the interfaces
		ip := getIPv4ForInterfaceName("wlan0")
		if ip == nil {
			log.Fatal("could not get IP address for eth0")
		}
		targetch := make(chan string)
		go genIPs(ip.String(), &targetch, 1337)
		for t := range targetch {
			con, err := net.DialTimeout("tcp", t, time.Second*2)
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("connected to:%s\n", t)
				con.Close()
				break
			}
		}
	*/
	echoClient, _, err := service.DialEchoService("tcp", `127.0.0.1:9527`)
	if err != nil {
		log.Fatalf("service.DialEchoService: %v", err)
	}
	defer echoClient.Close()

	args := &service.EchoRequest{Msg: proto.String("你好, 世界!")}
	reply := &service.EchoResponse{}
	err = echoClient.EchoTwice(args, reply)
	if err != nil {
		log.Fatalf("echoClient.EchoTwice: %v", err)
	}
	fmt.Println(reply.GetMsg())
}

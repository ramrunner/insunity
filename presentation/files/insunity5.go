package main

import (
	service "./entropy.pb"
	"bytes"
	"flag"
	"fmt"
	"github.com/chai2010/protorpc/proto"
	"log"
	"log/syslog"
	"math"
	"math/rand"
	"net"
	"strings"
	"time"
)

type logentry struct {
	time time.Time
	host string
	val  float32
}

type shannonsrv struct {
	log     []logentry
	id      string
	version uint32
}

func H(a string) float32 {
	freqs := make(map[byte]float64)
	abytes := []byte(a)
	entropy := float64(0)
	for i := 0; i < 256; i++ {
		freqs[byte(i)] = 0
	}
	ln := len(a)
	for k, _ := range freqs {
		px := float64(bytes.Count(abytes, []byte{k})) / float64(ln)
		freqs[k] = px
		if px > 0 {
			entropy += -float64(px) * math.Log2(px)
		}
	}
	return float32(entropy)
}

func (t *shannonsrv) Entropy(args *service.EntropyRequest, reply *service.EntropyResponse) error {
	reply.Entropy = proto.Float32(H(args.GetData()))
	reply.Id = proto.String(t.id)
	reply.Version = proto.Uint32(t.version)
	return nil
}

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
	return
}

func randStr(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	logwriter, e := syslog.New(syslog.LOG_NOTICE, "insunity")
	if e == nil {
		log.SetOutput(logwriter)
	}
	srv_flag := flag.Bool("srv", false, "should insunity start in server mode?")
	flag.Parse()
	ip := getIPv4ForInterfaceName("wlan0")
	if ip == nil {
		log.Fatal("could not get IP address for eth0")
	}
	if *srv_flag == true {
		ticker := time.NewTicker(time.Minute)
		ss := &shannonsrv{log: []logentry{}, id: ip.String(), version: 1}
		go func() {
			for t := range ticker.C {
				log.Printf("server:%s log:%s", ss.id, ss.log)
			}
		}()
		err := service.ListenAndServeEntropyService("tcp", fmt.Sprintf("%s:%d", ip, 9527), ss)
	} else {
		targetch := make(chan string)
		go genIPs(ip.String(), &targetch, 9527)
		for t := range targetch {
			log.Printf("will scan:%s\n", t)
			entropyClient, _, err := service.DialEntropyServiceTimeout("tcp", t, time.Second)
			if err != nil {
				log.Fatalf("service.DialEntropyService: %v", err)
			}
		}
		defer entropyClient.Close()
		for i := 0; i < 10; i++ {
			dat := proto.String(randStr(32))
			args := &service.EntropyRequest{Data: dat, Id: proto.String(ip.String()), Version: proto.Uint32(1)}
			reply := &service.EntropyResponse{}
			err = entropyClient.Entropy(args, reply)
			if err != nil {
				log.Fatalf("entropyClient.Entropy: %v", err)
			}
			fmt.Println(reply.GetEntropy(), reply.GetId())
		}
	}
}

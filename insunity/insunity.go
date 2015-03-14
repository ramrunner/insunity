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
	time time.Time //time when request was serviced
	host string    //caller id (IP)
	val  float32   // value returned to caller
}

type logslice []logentry

//the context struct of our server
type shannonsrv struct {
	log     logslice //all the requests served
	id      string   //Id string. usually IP
	version uint32   //service version number
}

//implementing the String interface for printing
func (l logslice) String() string {
	var ret []string
	for _, lg := range l {
		ret = append(ret, fmt.Sprintf("%s@%s -> %v", lg.host, lg.time, lg.val))
	}
	return strings.Join(ret, ", ")
}

//Shannon Entropy of a string
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

//RPC service call. arguments are described in the service pkg generated from the .proto file
func (t *shannonsrv) Entropy(args *service.EntropyRequest, reply *service.EntropyResponse) error {
	entr := H(args.GetData())
	reply.Entropy = proto.Float32(entr)
	reply.Id = proto.String(t.id)
	reply.Version = proto.Uint32(t.version)
	t.log = append(t.log, logentry{time: time.Now(), host: args.GetId(), val: entr})
	return nil
}

//get the IPv4 of an interface
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

//generate ip:port strings for the subnet ipstr
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

//generate a random string of length n
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
	srv_flag := flag.Bool("srv", false, "start in server mode. other options are ignored. eth0 IP will be used")
	port_flag := flag.Int("port", 9527, "port for the server to bind on")
	msg_amount_flag := flag.Int("amount", 10, "amount of random messages to generate")
	msg_length_flag := flag.Int("len", 32, "character length of each random message")
	srv_ip_flag := flag.String("dial", "", "IP:Port of server to dial")
	flag.Parse()
	ip := getIPv4ForInterfaceName("eth0")
	if ip == nil {
		log.Fatal("could not get IP address for eth0")
	}
	if *srv_flag == true {
		ticker := time.NewTicker(time.Minute)
		ss := &shannonsrv{log: []logentry{}, id: ip.String(), version: 1}
		go func() {
			for _ = range ticker.C {
				log.Printf("server:%s log:%s", ss.id, ss.log)
			}
		}()
		err := service.ListenAndServeEntropyService("tcp", fmt.Sprintf("%s:%d", ip, *port_flag), ss)
		ticker.Stop()
		log.Fatal(err)
	} else {
		var (
			entropyClient *service.EntropyServiceClient
			err           error
		)
		targetch := make(chan string)
		if *srv_ip_flag == "" {
			go genIPs(ip.String(), &targetch, 9527)
		} else {
			go func() {
				targetch <- *srv_ip_flag
				close(targetch)
			}()
		}
		for t := range targetch {
			log.Printf("connecting to:%s with 1s timeout\n", t)
			entropyClient, _, err = service.DialEntropyServiceTimeout("tcp", t, time.Second)
			if err != nil {
				log.Printf("service.DialEntropyService: %v", err)
			}
			if entropyClient != nil {
				defer entropyClient.Close()
				break
			}
		}
		if entropyClient == nil { //we never found a server
			log.Fatal("no servers could be contacted.")
		}
		for i := 0; i < *msg_amount_flag; i++ {
			dat := proto.String(randStr(*msg_length_flag))
			args := &service.EntropyRequest{Data: dat, Id: proto.String(ip.String()), Version: proto.Uint32(1)}
			reply := &service.EntropyResponse{}
			err = entropyClient.Entropy(args, reply)
			if err != nil {
				log.Fatalf("entropyClient.Entropy: %v", err)
			}
			fmt.Println("server:", reply.GetId(), " str", args.GetData(), " entropy ", reply.GetEntropy())
		}
	}
}

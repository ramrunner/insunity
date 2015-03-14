package main

import (
	"log"
	"log/syslog"
)

func main() {
	//set log to syslog
	logwriter, e := syslog.New(syslog.LOG_NOTICE, "insunity")
	if e == nil {
		log.SetOutput(logwriter)
	}
	log.Printf("χελόου γουορλντ\n")
}

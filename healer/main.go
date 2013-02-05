package main

import (
	"fmt"
	"log/syslog"
	"os"
	"time"
)

func main() {
	var err error
	log, err = syslog.New(syslog.LOG_INFO, "tsuru-healer")
	if err != nil {
		panic(err)
	}
	if len(os.Args) < 1 {
		fmt.Println("Healer expects the endpoint to connect with tsuru.")
		return
	}
	endpoint := os.Args[1]
	registerTicker(time.Tick(time.Minute*15), endpoint)
	healTicker(time.Tick(time.Minute))
}

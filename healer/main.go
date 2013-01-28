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
	if len(os.Args) < 3 {
		fmt.Println("Healer expects email, password and endpoint to connect with tsuru.")
		return
	}
	// email := os.Args[1]
	// password := os.Args[2]
	endpoint := os.Args[3]
	// healer := newInstanceHealer(email, password, endpoint)
	// register("instance-healer", healer)
	registerTicker(time.Tick(time.Minute * 15), endpoint)
	healTicker(time.Tick(time.Minute))
}

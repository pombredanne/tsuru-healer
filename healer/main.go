// +build ignore

package main

import (
	"fmt"
	"github.com/globocom/tsuru-healer/healer"
	"log/syslog"
	"os"
	"time"
)

func main() {
	log, err := syslog.New(syslog.LOG_INFO, "tsuru-healer")
	if err != nil {
		panic(err)
	}
	if len(os.Args) < 3 {
		fmt.Println("Healer expects email, password and endpoint to connect with tsuru.")
		return
	}
	email := os.Args[1]
	password := os.Args[2]
	endpoint := os.Args[3]
	healer := healer.NewInstanceHealer(email, password, endpoint)
	for _ = range time.Tick(time.Minute) {
		err := healer.Heal()
		if err != nil {
			log.Err("Got error while healing: " + err.Error())
		}
	}
}

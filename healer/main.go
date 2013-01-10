// +build ignore

package main

import (
    "fmt"
    "github.com/globocom/tsuru-healer/healer"
	"os"
    "time"
    "log/syslog"
)

func main() {
    var log *syslog.Writer
    if len(os.Args) < 3 {
        fmt.Println("Healer expects email, password and endpoint to connect with tsuru.")
        return
    }
	email := os.Args[0]
	password := os.Args[1]
	endpoint := os.Args[2]
	healer := healer.NewTsuruHealer(email, password, endpoint)
    for _ = range time.Tick(time.Minute) {
        err := healer.Heal()
        if err != nil {
            log.Err("Got error while healing: " + err.Error())
        }
    }
}

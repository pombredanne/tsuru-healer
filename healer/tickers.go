package main

import (
	"fmt"
	"sync"
	"time"
)

// healTicker execute the healers.heal.
func healTicker(ticker <-chan time.Time) {
	log.Info("running heal ticker")
	var wg sync.WaitGroup
	for _ = range ticker {
		healers := getHealers()
		wg.Add(len(healers))
		for name, h := range healers {
			log.Info(fmt.Sprintf("running verification/heal for %s", name))
			go func(healer *healer) {
				err := healer.heal()
				if err != nil {
					log.Info(err.Error())
				}
				wg.Done()
			}(h)
		}
		wg.Wait()
	}
}

// registerTicker register healers from resource.
func registerTicker(ticker <-chan time.Time, endpoint string) {
	var registerHealer = func() {
		log.Info("running register ticker")
		if healers, err := healersFromResource(endpoint); err == nil {
			setHealers(healers)
		}
	}
	registerHealer()
	go func() {
		for _ = range ticker {
			registerHealer()
		}
	}()
}

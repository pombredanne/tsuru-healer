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
			go func(healer healer) {
				err := healer.heal()
				if err != nil {
					log.Info(err)
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
		healers, _ := healersFromResource(endpoint)
		for name, healer := range healers {
			register(name, &healer)
		}
	}
	registerHealer()
	go func() {
		for _ = range ticker {
			registerHealer()
		}
	}()
}

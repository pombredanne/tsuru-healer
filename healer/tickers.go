package main

import (
	"sync"
	"time"
)

// healTicker execute the healers.heal.
func healTicker(ticker <-chan time.Time) {
	var wg sync.WaitGroup
	for _ = range ticker {
		healers := getHealers()
		wg.Add(len(healers))
		for _, h := range healers {
			go func(healer healer) {
				healer.heal()
				wg.Done()
			}(h)
		}
		wg.Wait()
	}
}

// registerTicker register healers from resource.
func registerTicker(ticker <-chan time.Time, endpoint string) {
	var registerHealer = func() {
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

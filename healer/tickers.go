package main

import "time"

// healTicker execute the healers.heal.
func healTicker(ticker chan time.Time) {
	for _ = range ticker {
		healers := getHealers()
		for _, healer := range healers {
			go healer.heal()
		}
	}
}

func registerTicker(ticker chan time.Time, endpoint string) {
	go func() {
		for _ = range ticker {
			healers, _ := healersFromResource(endpoint)
			for _, healer := range healers {
				register(&healer)
			}
		}
	}()
}

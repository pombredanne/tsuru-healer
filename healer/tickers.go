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

package main

import (
	. "launchpad.net/gocheck"
	"net/http"
	"net/http/httptest"
	"time"
)

func (s *S) TestHealTicker(c *C) {
	var called bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()
	h := &tsuruHealer{url: ts.URL}
	register("ticker-healer", h)
	ch := make(chan time.Time)
	ok := make(chan bool)
	go func() {
		healTicker(ch)
		ok <- true
	}()
	ch <- time.Now()
	time.Sleep(1 * time.Second)
	close(ch)
	<-ok
	c.Assert(called, Equals, true)
}

func (s *S) TestRegisterTicker(c *C) {
	var called bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()
	ch := make(chan time.Time)
	ok := make(chan bool)
	go func() {
		registerTicker(ch, ts.URL)
		ok <- true
	}()
	ch <- time.Now()
	time.Sleep(1 * time.Second)
	close(ch)
	<-ok
	c.Assert(called, Equals, true)
}

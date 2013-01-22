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
	register(h)
	ch := make(chan time.Time)
	ok := make(chan bool)
	go func() {
		healTicker(ch)
		ok <- true
	}()
	ch <- time.Now()
	time.Sleep(1 * time.Millisecond)
	close(ch)
	<-ok
	c.Assert(called, Equals, true)
}

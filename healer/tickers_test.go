package main

import (
	. "launchpad.net/gocheck"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"time"
)

func (s *S) TestHealTicker(c *C) {
	var called int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.StoreInt32(&called, 1)
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
	c.Assert(atomic.LoadInt32(&called), Equals, int32(1))
}

func (s *S) TestRegisterTicker(c *C) {
	var called int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.StoreInt32(&called, 1)
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
	c.Assert(atomic.LoadInt32(&called), Equals, int32(1))
}

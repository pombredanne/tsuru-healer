package main

import (
	"fmt"
	. "launchpad.net/gocheck"
	"net/http"
	"net/http/httptest"
)

func (s *S) TestHealersFromResource(c *C) {
	reqs := []*http.Request{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqs = append(reqs, r)
		w.Write([]byte(`{"bootstrap":"/bootstrap"}`))
	}))
	defer ts.Close()
	expected := map[string]*healer{
		"bootstrap": {url: fmt.Sprintf("%s/bootstrap", ts.URL)},
	}
	healers, err := healersFromResource(ts.URL)
	c.Assert(err, IsNil)
	c.Assert(healers, DeepEquals, expected)
}

func (s *S) TestTsuruHealer(c *C) {
	var called bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()
	h := healer{url: ts.URL}
	err := h.heal()
	c.Assert(err, IsNil)
	c.Assert(called, Equals, true)
}

func (s *S) TestSetAndGetHealers(c *C) {
	h := &healer{url: ""}
	setHealers(map[string]*healer{"test-healer": h})
	healers := getHealers()
	healer, ok := healers["test-healer"]
	c.Assert(healer, DeepEquals, h)
	c.Assert(ok, Equals, true)
}

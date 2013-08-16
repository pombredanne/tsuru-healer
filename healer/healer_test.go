package main

import (
	"fmt"
	. "launchpad.net/gocheck"
	"net/http"
	"net/http/httptest"
	"os"
)

func (s *S) TestHealersFromResource(c *C) {
	os.Setenv("TSURU_TOKEN", "token123")
	defer os.Setenv("TSURU_TOKEN", "")
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
	c.Assert(reqs, HasLen, 1)
	c.Assert(reqs[0].Header.Get("Authorization"), Equals, "bearer token123")
}

func (s *S) TestTsuruHealer(c *C) {
	os.Setenv("TSURU_TOKEN", "token123")
	defer os.Setenv("TSURU_TOKEN", "")
	var reqs []*http.Request
	var called bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		reqs = append(reqs, r)
	}))
	defer ts.Close()
	h := healer{url: ts.URL}
	err := h.heal()
	c.Assert(err, IsNil)
	c.Assert(called, Equals, true)
	c.Assert(reqs, HasLen, 1)
	c.Assert(reqs[0].Header.Get("Authorization"), Equals, "bearer token123")
}

func (s *S) TestSetAndGetHealers(c *C) {
	h := &healer{url: ""}
	setHealers(map[string]*healer{"test-healer": h})
	healers := getHealers()
	healer, ok := healers["test-healer"]
	c.Assert(healer, DeepEquals, h)
	c.Assert(ok, Equals, true)
}

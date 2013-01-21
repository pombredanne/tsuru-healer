package healer

import (
	"github.com/flaviamissi/go-elb/aws"
	"github.com/flaviamissi/go-elb/elb"
	. "launchpad.net/gocheck"
	"net/http"
	"net/http/httptest"
)

func (s *S) TestGetToken(c *C) {
	var req *http.Request
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req = r
		w.Write([]byte(`{"token": "token-123"}`))
	}))
	defer ts.Close()
	token, err := getToken("test@test.com", "test123", ts.URL)
	c.Assert(err, IsNil)
	c.Assert(token, Equals, "token-123")
	c.Assert(req.URL.String(), Equals, "/users/test@test.com/tokens")
}

func (s *S) TestSpawn(c *C) {
	var req *http.Request
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req = r
		w.Write([]byte(""))
	}))
	defer ts.Close()
	s.healer.Endpoint = ts.URL
	err := s.healer.Spawn("testlb")
	c.Assert(err, IsNil)
	c.Assert(req.URL.String(), Equals, "/apps/testlb/units")
	c.Assert(req.Method, Equals, "PUT")
}

func (s *S) TestTerminate(c *C) {
	var req *http.Request
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req = r
		w.Write([]byte(""))
	}))
	defer ts.Close()
	s.healer.Endpoint = ts.URL
	err := s.healer.Terminate("testlb", "i-123")
	c.Assert(err, IsNil)
	c.Assert(req.URL.String(), Equals, "/apps/testlb/unit")
	c.Assert(req.Method, Equals, "DELETE")
}

func (s *S) TestHealer(c *C) {
	reqs := []*http.Request{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqs = append(reqs, r)
		w.Write([]byte("123"))
	}))
	defer ts.Close()
	s.healer.seeker = &AWSSeeker{
		ELB: elb.New(aws.Auth{AccessKey: "auth", SecretKey: "s3cr3t"}, aws.Region{ELBEndpoint: s.elbsrv.URL()}),
	}
	s.healer.Endpoint = ts.URL
	state := elb.InstanceState{
		Description: "Instance has failed at least the UnhealthyThreshold number of health checks consecutively.",
		State:       "OutOfService",
		ReasonCode:  "Instance",
		InstanceId:  s.instId,
	}
	s.elbsrv.ChangeInstanceState("testlb", state)
	err := s.healer.Heal()
	c.Assert(err, IsNil)
	c.Assert(len(reqs), Equals, 2)
	c.Assert(reqs[0].URL.String(), Equals, "/apps/testlb/unit")
	c.Assert(reqs[0].Method, Equals, "DELETE")
	c.Assert(reqs[0].Header.Get("Authorization"), Equals, s.token)
	c.Assert(reqs[1].URL.String(), Equals, "/apps/testlb/units")
	c.Assert(reqs[1].Method, Equals, "PUT")
	c.Assert(reqs[1].Header.Get("Authorization"), Equals, s.token)
}

func (s *S) TestHealersFromResource(c *C) {
	reqs := []*http.Request{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqs = append(reqs, r)
		w.Write([]byte(`{"bootstrap":"/bootstrap"}`))
	}))
	defer ts.Close()
	expected := map[string]string{
		"bootstrap": "/bootstrap",
	}
	healers, err := healersFromResource(ts.URL)
	c.Assert(err, IsNil)
	c.Assert(healers, DeepEquals, expected)
}

func (s *S) TestGenericHealer(c *C) {
	var called bool
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))
	defer ts.Close()
	h := GenericHealer{
		url: ts.URL,
	}
	err := h.Heal()
	c.Assert(err, IsNil)
	c.Assert(called, Equals, true)
}

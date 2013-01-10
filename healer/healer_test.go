package healer

import (
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
	c.Assert(req.URL.String(), Equals, "/apps/testlb/units")
	c.Assert(req.Method, Equals, "DELETE")
}

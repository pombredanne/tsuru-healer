package main

import (
	. "launchpad.net/gocheck"
	"log/syslog"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type S struct {
	instId string
	token  string
}

var _ = Suite(&S{})

func (s *S) SetUpSuite(c *C) {
	var err error
	log, err = syslog.New(syslog.LOG_INFO, "tsuru-healer")
	c.Assert(err, IsNil)
}

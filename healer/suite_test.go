package main

import (
	"github.com/flaviamissi/go-elb/aws"
	"github.com/flaviamissi/go-elb/elb"
	"github.com/flaviamissi/go-elb/elb/elbtest"
	. "launchpad.net/gocheck"
	"log/syslog"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type S struct {
	elbsrv *elbtest.Server
	seeker seeker
	healer *instanceHealer
	instId string
	token  string
}

var _ = Suite(&S{})

func (s *S) SetUpSuite(c *C) {
	var err error
	log, err = syslog.New(syslog.LOG_INFO, "tsuru-healer")
	c.Assert(err, IsNil)
	s.setUpELB(c)
	s.token = "123456"
	s.healer = &instanceHealer{token: s.token}
}

func (s *S) setUpELB(c *C) {
	var err error
	s.elbsrv, err = elbtest.NewServer()
	c.Assert(err, IsNil)
	s.elbsrv.NewLoadBalancer("testlb")
	s.instId = s.elbsrv.NewInstance()
	s.elbsrv.RegisterInstance(s.instId, "testlb")
	s.seeker = &awsSeeker{
		elb: elb.New(aws.Auth{AccessKey: "auth", SecretKey: "s3cr3t"}, aws.Region{ELBEndpoint: s.elbsrv.URL()}),
	}
}

func (s *S) TearDownSuite(c *C) {
	s.elbsrv.Quit()
}

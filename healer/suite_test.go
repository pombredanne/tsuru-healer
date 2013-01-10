package healer

import (
	"github.com/flaviamissi/go-elb/aws"
	"github.com/flaviamissi/go-elb/elb"
	"github.com/flaviamissi/go-elb/elb/elbtest"
	. "launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type S struct {
	elbsrv *elbtest.Server
	seeker Seeker
	healer *TsuruHealer
	instId string
}

var _ = Suite(&S{})

func (s *S) SetUpSuite(c *C) {
	s.setUpELB(c)
	s.healer = &TsuruHealer{}
}

func (s *S) setUpELB(c *C) {
	var err error
	s.elbsrv, err = elbtest.NewServer()
	c.Assert(err, IsNil)
	s.elbsrv.NewLoadBalancer("testlb")
	s.instId = s.elbsrv.NewInstance()
	s.elbsrv.RegisterInstance(s.instId, "testlb")
	s.seeker = &AWSSeeker{
		ELB: elb.New(aws.Auth{AccessKey: "auth", SecretKey: "s3cr3t"}, aws.Region{ELBEndpoint: s.elbsrv.URL()}),
	}
}

func (s *S) TearDownSuite(c *C) {
	s.elbsrv.Quit()
}

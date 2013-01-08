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
	srv    *elbtest.Server
	seeker Seeker
}

var _ = Suite(&S{})

func (s *S) SetUpSuite(c *C) {
	var err error
	s.srv, err = elbtest.NewServer()
	c.Assert(err, IsNil)
	s.srv.NewLoadBalancer("testlb")
	s.seeker = AWSSeeker{
		ELB: elb.New(aws.Auth{"auth", "s3cr3t"}, aws.Region{ELBEndpoint: s.srv.URL()}),
	}
}

func (s *S) TearDownSuite(c *C) {
	s.srv.Quit()
}

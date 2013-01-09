package healer

import (
	"github.com/flaviamissi/go-elb/elb"
	. "launchpad.net/gocheck"
)

func (s *S) TestDescribeLoadBalancers(c *C) {
	lbs, err := s.seeker.DescribeLoadBalancers()
	c.Assert(err, IsNil)
	c.Assert(len(lbs) > 0, Equals, true)
	c.Assert(lbs[0].Name, Equals, "testlb")
	c.Assert(lbs[0].DNSName, Matches, "^testlb.*")
}

func (s *S) TestDescribeInstancesHealth(c *C) {
	instances, err := s.seeker.DescribeInstancesHealth("testlb")
	c.Assert(err, IsNil)
	c.Assert(len(instances) > 0, Equals, true)
	c.Assert(instances[0].InstanceId, Equals, s.instId)
	c.Assert(instances[0].Description, Not(Equals), "")
	c.Assert(instances[0].ReasonCode, Not(Equals), "")
	c.Assert(instances[0].State, Not(Equals), "")
}

func (s *S) TestSeekUnhealthyInstances(c *C) {
	state := elb.InstanceState{
		Description: "Instance has failed at least the UnhealthyThreshold number of health checks consecutively",
		State:       "OutOfService",
		ReasonCode:  "Instance",
		InstanceId:  s.instId,
	}
	s.srv.ChangeInstanceState("testlb", state)
	instances, err := s.seeker.SeekUnhealthyInstances()
	c.Assert(err, IsNil)
	expected := []Instance{
		{
			Description: "Instance has failed at least the UnhealthyThreshold number of health checks consecutively",
			State:       "OutOfService",
			ReasonCode:  "Instance",
			InstanceId:  s.instId,
		},
	}
	c.Assert(instances, DeepEquals, expected)
}

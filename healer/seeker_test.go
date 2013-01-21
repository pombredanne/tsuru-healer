package main

import (
	"github.com/flaviamissi/go-elb/elb"
	. "launchpad.net/gocheck"
)

func (s *S) TestDescribeLoadBalancers(c *C) {
	lbs, err := s.seeker.describeLoadBalancers()
	c.Assert(err, IsNil)
	c.Assert(len(lbs) > 0, Equals, true)
	c.Assert(lbs[0].name, Equals, "testlb")
	c.Assert(lbs[0].dnsName, Matches, "^testlb.*")
}

func (s *S) TestDescribeInstancesHealth(c *C) {
	instances, err := s.seeker.describeInstancesHealth("testlb")
	c.Assert(err, IsNil)
	c.Assert(len(instances) > 0, Equals, true)
	c.Assert(instances[0].instanceId, Equals, s.instId)
	c.Assert(instances[0].description, Not(Equals), "")
	c.Assert(instances[0].reasonCode, Not(Equals), "")
	c.Assert(instances[0].state, Not(Equals), "")
	c.Assert(instances[0].loadBalancer, Equals, "testlb")
}

func (s *S) TestSeekUnhealthyInstances(c *C) {
	state := elb.InstanceState{
		Description: "Instance has failed at least the UnhealthyThreshold number of health checks consecutively.",
		State:       "OutOfService",
		ReasonCode:  "Instance",
		InstanceId:  s.instId,
	}
	s.elbsrv.ChangeInstanceState("testlb", state)
	instances, err := s.seeker.seekUnhealthyInstances()
	c.Assert(err, IsNil)
	expected := []instance{
		{
			description:  "Instance has failed at least the UnhealthyThreshold number of health checks consecutively.",
			state:        "OutOfService",
			reasonCode:   "Instance",
			instanceId:   s.instId,
			loadBalancer: "testlb",
		},
	}
	c.Assert(instances, DeepEquals, expected)
}

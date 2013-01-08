package healer

import (
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

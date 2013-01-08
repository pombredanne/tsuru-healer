package healer

import (
	. "launchpad.net/gocheck"
)

func (s *S) TestDescribeLoadBalancers(c *C) {
	lbs, err := s.seeker.DescribeLoadBalancers()
	c.Assert(err, IsNil)
	c.Assert(len(lbs) > 0, Equals, true)
	c.Assert(lbs[0].Name, Equals, "testlb")
}

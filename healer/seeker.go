package healer

import (
	"github.com/flaviamissi/go-elb/aws"
	"github.com/flaviamissi/go-elb/elb"
)

// Seeker is responsible to seek for unhealthy instances
// under a Load Balancer
type Seeker interface {
	DescribeInstancesHealth() ([]Instance, error)
	DescribeLoadBalancers() ([]LoadBalancer, error)
}

type Instance struct {
	InstanceId  string
	Description string
	ReasonCode  string
	State       string
}

type LoadBalancer struct {
	AvailZones []string
	DNSName    string
	Instances  []Instance
	Name       string
}

type AWSSeeker struct {
	ELB *elb.ELB
}

func NewAWSSeeker() AWSSeeker {
	auth, err := aws.EnvAuth()
	if err != nil {
		panic(err.Error())
	}
	return AWSSeeker{
		ELB: elb.New(auth, aws.USEast),
	}
}

func (aws AWSSeeker) DescribeInstancesHealth() ([]Instance, error) {
	return nil, nil
}

func (aws AWSSeeker) DescribeLoadBalancers() ([]LoadBalancer, error) {
	lbResp, err := aws.ELB.DescribeLoadBalancers()
	if err != nil {
		return nil, err
	}
	lbs := []LoadBalancer{}
	for _, lbDesc := range lbResp.LoadBalancerDescriptions {
		lb := LoadBalancer{
			AvailZones: lbDesc.AvailZones,
			Name:       lbDesc.LoadBalancerName,
		}
		instances := make([]Instance, len(lbDesc.Instances))
		for i, instance := range lbDesc.Instances {
			instances[i].InstanceId = instance.InstanceId
		}
		lbs = append(lbs, lb)
	}
	return lbs, nil
}

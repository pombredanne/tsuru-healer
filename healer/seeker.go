package healer

import (
	"github.com/flaviamissi/go-elb/aws"
	"github.com/flaviamissi/go-elb/elb"
	"log/syslog"
)

// Seeker is responsible to seek for unhealthy instances
// under a Load Balancer
type Seeker interface {
	DescribeInstancesHealth(lb string) ([]Instance, error)
	DescribeLoadBalancers() ([]LoadBalancer, error)
	SeekUnhealthyInstances() ([]Instance, error)
}

type Instance struct {
	InstanceId   string
	Description  string
	ReasonCode   string
	State        string
	LoadBalancer string
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

var log *syslog.Writer

func NewAWSSeeker() *AWSSeeker {
	auth, err := aws.EnvAuth()
	if err != nil {
		panic(err.Error())
	}
	// receive region?
	return &AWSSeeker{
		ELB: elb.New(auth, aws.USEast),
	}
}

func (s *AWSSeeker) matchCriteria(instances []Instance, model Instance) []Instance {
	matches := []Instance{}
	for _, instance := range instances {
		if instance.Description == model.Description && instance.State == model.State &&
			instance.ReasonCode == model.ReasonCode {
			matches = append(matches, instance)
		}
	}
	return matches
}

func (s *AWSSeeker) SeekUnhealthyInstances() ([]Instance, error) {
	log.Info("Seeking for unhealthy instances..")
	lbs, err := s.DescribeLoadBalancers()
	if err != nil {
		return nil, err
	}
	unhealthy := []Instance{}
	for _, lb := range lbs {
		instances, err := s.DescribeInstancesHealth(lb.Name)
		if err != nil {
			return nil, err
		}
		model := Instance{
			Description: "Instance has failed at least the UnhealthyThreshold number of health checks consecutively",
			State:       "OutOfService",
			ReasonCode:  "Instance",
		}
		unhealthy = append(unhealthy, s.matchCriteria(instances, model)...)
	}
	log.Info("Found " + len(unhealthy) + " unhealthy instances.")
	return unhealthy, nil
}

func (s *AWSSeeker) DescribeInstancesHealth(lb string) ([]Instance, error) {
	resp, err := s.ELB.DescribeInstanceHealth(lb)
	if err != nil {
		return nil, err
	}
	instances := make([]Instance, len(resp.InstanceStates))
	for i, state := range resp.InstanceStates {
		instances[i].InstanceId = state.InstanceId
		instances[i].Description = state.Description
		instances[i].ReasonCode = state.ReasonCode
		instances[i].State = state.State
		instances[i].LoadBalancer = lb
	}
	return instances, nil
}

func (s *AWSSeeker) DescribeLoadBalancers() ([]LoadBalancer, error) {
	resp, err := s.ELB.DescribeLoadBalancers()
	if err != nil {
		return nil, err
	}
	lbs := []LoadBalancer{}
	for _, lbDesc := range resp.LoadBalancerDescriptions {
		lb := LoadBalancer{
			AvailZones: lbDesc.AvailZones,
			Name:       lbDesc.LoadBalancerName,
			DNSName:    lbDesc.DNSName,
		}
		instances := make([]Instance, len(lbDesc.Instances))
		for i, instance := range lbDesc.Instances {
			instances[i].InstanceId = instance.InstanceId
		}
		lbs = append(lbs, lb)
	}
	return lbs, nil
}

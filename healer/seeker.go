package main

import (
	"fmt"
	"github.com/flaviamissi/go-elb/aws"
	"github.com/flaviamissi/go-elb/elb"
)

// Seeker is responsible to seek for unhealthy instances
// under a Load Balancer
type seeker interface {
	describeInstancesHealth(lb string) ([]instance, error)
	describeLoadBalancers() ([]loadBalancer, error)
	seekUnhealthyInstances() ([]instance, error)
}

type instance struct {
	instanceId   string
	description  string
	reasonCode   string
	state        string
	loadBalancer string
}

type loadBalancer struct {
	availZones []string
	dnsName    string
	instances  []instance
	name       string
}

type awsSeeker struct {
	elb *elb.ELB
}

func newAWSSeeker() *awsSeeker {
	auth, err := aws.EnvAuth()
	if err != nil {
		panic(err.Error())
	}
	// receive region?
	return &awsSeeker{
		elb: elb.New(auth, aws.USEast),
	}
}

func (s *awsSeeker) matchCriteria(instances []instance, model instance) []instance {
	matches := []instance{}
	for _, instance := range instances {
		if instance.description == model.description && instance.state == model.state &&
			instance.reasonCode == model.reasonCode {
			matches = append(matches, instance)
		}
	}
	return matches
}

func (s *awsSeeker) seekUnhealthyInstances() ([]instance, error) {
	log.Info("Seeking for unhealthy instances..")
	lbs, err := s.describeLoadBalancers()
	if err != nil {
		return nil, err
	}
	unhealthy := []instance{}
	for _, lb := range lbs {
		instances, err := s.describeInstancesHealth(lb.name)
		if err != nil {
			return nil, err
		}
		model := instance{
			description: "Instance has failed at least the UnhealthyThreshold number of health checks consecutively.",
			state:       "OutOfService",
			reasonCode:  "Instance",
		}
		unhealthy = append(unhealthy, s.matchCriteria(instances, model)...)
	}
	log.Info(fmt.Sprintf("Found %d unhealthy instances.", len(unhealthy)))
	return unhealthy, nil
}

func (s *awsSeeker) describeInstancesHealth(lb string) ([]instance, error) {
	resp, err := s.elb.DescribeInstanceHealth(lb)
	if err != nil {
		return nil, err
	}
	instances := make([]instance, len(resp.InstanceStates))
	for i, state := range resp.InstanceStates {
		instances[i].instanceId = state.InstanceId
		instances[i].description = state.Description
		instances[i].reasonCode = state.ReasonCode
		instances[i].state = state.State
		instances[i].loadBalancer = lb
	}
	return instances, nil
}

func (s *awsSeeker) describeLoadBalancers() ([]loadBalancer, error) {
	resp, err := s.elb.DescribeLoadBalancers()
	if err != nil {
		return nil, err
	}
	lbs := []loadBalancer{}
	for _, lbDesc := range resp.LoadBalancerDescriptions {
		lb := loadBalancer{
			availZones: lbDesc.AvailZones,
			name:       lbDesc.LoadBalancerName,
			dnsName:    lbDesc.DNSName,
		}
		instances := make([]instance, len(lbDesc.Instances))
		for i, instance := range lbDesc.Instances {
			instances[i].instanceId = instance.InstanceId
		}
		lbs = append(lbs, lb)
	}
	return lbs, nil
}

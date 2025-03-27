package ssmutil

import (
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/aws"
)

func GetInstanceIDByPrivateIP(ip, region, profile string) string {
	sess := GetAWSSession(region, profile)
	svc := ec2.New(sess)
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{Name: aws.String("private-ip-address"), Values: []*string{aws.String(ip)}},
			{Name: aws.String("instance-state-name"), Values: []*string{aws.String("running")}},
		},
	}
	result, err := svc.DescribeInstances(input)
	if err != nil || len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
		log.Printf("No instance found with private IP %s in region %s", ip, region)
		return ""
	}
	return *result.Reservations[0].Instances[0].InstanceId
}

func GetInstanceIDByTag(tag, region, profile string) string {
	sess := GetAWSSession(region, profile)
	svc := ec2.New(sess)
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{Name: aws.String("tag:Name"), Values: []*string{aws.String(tag)}},
			{Name: aws.String("instance-state-name"), Values: []*string{aws.String("running")}},
		},
	}
	result, err := svc.DescribeInstances(input)
	if err != nil || len(result.Reservations) == 0 {
		log.Fatalf("Error: No instance found for tag:Name = %s in region %s", tag, region)
	}
	for _, r := range result.Reservations {
		for _, inst := range r.Instances {
			if inst.Placement != nil && inst.Placement.AvailabilityZone != nil {
				az := *inst.Placement.AvailabilityZone
				if strings.HasPrefix(az, region) {
					return *inst.InstanceId
				}
			}
		}
	}
	log.Fatalf("Error: No running instance with tag:Name = %s found in region %s", tag, region)
	return ""
}

func GetInstanceID(hostname, region, profile string) string {
	ip := ResolveHostnameToIP(hostname)
	if ip != "" {
		log.Printf("Resolved private IP: %s", ip)
		if id := GetInstanceIDByPrivateIP(ip, region, profile); id != "" {
			return id
		}
	}
	log.Printf("Looking up instance by tag:Name = %s in %s...", hostname, region)
	return GetInstanceIDByTag(hostname, region, profile)
}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/singularet/ssmtools/internal/ssmutil"
)

func startSSMSession(instanceID, region, profile string) {
	log.Printf("Starting SSM session on instance: %s in region %s using profile %s", instanceID, region, profile)
	cmd := exec.Command("bash", "-c", fmt.Sprintf("exec aws ssm start-session --target %s --region %s --profile %s", instanceID, region, profile))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error: Failed to start SSM session for %s: %v", instanceID, err)
	}
}

func main() {
	hostname := flag.String("hostname", "", "Hostname or private IP of the instance")
	region := flag.String("region", "us-west-2", "AWS region")
	role := flag.String("role", "", "AWS IAM role to assume using OneLogin")
	profile := flag.String("profile", "", "AWS CLI profile to use")
	flag.Parse()

	if *hostname == "" {
		fmt.Println("Usage: ssmssh --hostname <hostname> [--region us-west-2] [--role <role-name>] [--profile <aws-profile>]")
		os.Exit(1)
	}

	if *role != "" {
		ssmutil.AssumeAWSRole(*role)
		*profile = *role
	} else if *profile != "" {
		ssmutil.AssumeAWSRole(*profile)
	}

	instanceID := ssmutil.GetInstanceID(*hostname, *region, *profile)
	startSSMSession(instanceID, *region, *profile)
}

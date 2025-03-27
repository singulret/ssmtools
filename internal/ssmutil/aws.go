package ssmutil

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func AssumeAWSRole(profile string) {
	log.Printf("Assuming AWS role using profile: %s", profile)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}
	scriptPath := filepath.Join(homeDir, "workspace/singular/profiles/aws-onelogin-access.sh")

	cmd := exec.Command("bash", "-i", scriptPath, "--profile", profile)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("ONELOGIN_ROLE=%s", profile),
		fmt.Sprintf("AWS_PROFILE=%s", profile),
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("AWS role assumption script failed: %v", err)
	}
}

func GetAWSSession(region, profile string) *session.Session {
	options := session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}
	if region != "" {
		options.Config.Region = aws.String(region)
	}
	if profile != "" {
		options.Profile = profile
	}
	return session.Must(session.NewSessionWithOptions(options))
}

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"strings"

	"github.com/singularet/ssmtools/internal/ssmutil"

)

func startSSMTunnel(instanceID string, remotePort, localPort int, profile, region string) *exec.Cmd {
	log.Printf("Starting port forward from localhost:%d to instance:%d via SSM", localPort, remotePort)

	args := []string{
		"ssm", "start-session",
		"--target", instanceID,
		"--document-name", "AWS-StartPortForwardingSession",
		"--parameters", fmt.Sprintf("{\"portNumber\":[\"%d\"], \"localPortNumber\":[\"%d\"]}", remotePort, localPort),
		"--region", region,
	}
	if profile != "" {
		args = append(args, "--profile", profile)
	}

	log.Printf("Running using command: aws %s", strings.Join(args, " "))

	cmd := exec.Command("aws", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = nil

	if err := cmd.Start(); err != nil {
		log.Fatalf("SSM tunnel failed to start: %v", err)
	}

	time.Sleep(3 * time.Second)
	return cmd
}

func uploadFile(localPath, remotePath string) error {
	file, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://localhost:48080%s", remotePath)
	req, err := http.NewRequest("PUT", url, file)
	if err != nil {
		return err
	}
	req.ContentLength = stat.Size()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("upload failed with status: %s", resp.Status)
	}
	fmt.Println("✅ Upload complete")
	return nil
}

func main() {
	var (
		hostname   = flag.String("hostname", "", "Target hostname or tag:Name")
		profile    = flag.String("profile", "", "AWS CLI profile to use")
		role       = flag.String("role", "", "AWS IAM role to assume using OneLogin")
		region     = flag.String("region", "us-west-2", "AWS region")
		from       = flag.String("from", "", "Local path to upload")
		to         = flag.String("to", "", "Target path on remote")
		localPort  = flag.Int("local-port", 48080, "Local port for tunnel")
		remotePort = flag.Int("remote-port", 48080, "Remote port on instance")
	)
	flag.Parse()

	if *hostname == "" || *from == "" || *to == "" {
		fmt.Println("Usage: ssmcp --hostname bastion01 --from ./file.txt --to /home/ubuntu/shared/file.txt [--role role-name] [--profile dev] [--region us-west-2]")
		os.Exit(1)
	}

	if *role != "" {
		ssmutil.AssumeAWSRole(*role)
		*profile = *role
	} else if *profile != "" {
		ssmutil.AssumeAWSRole(*profile)
	}

	instanceID := ssmutil.GetInstanceID(*hostname, *region, *profile)
	absPath, _ := filepath.Abs(*from)
	log.Printf("Uploading %s to %s on instance %s...", absPath, *to, instanceID)

	tunnelCmd := startSSMTunnel(instanceID, *remotePort, *localPort, *profile, *region)
	defer func() {
		if tunnelCmd.Process != nil {
			_ = tunnelCmd.Process.Kill()
			_ = tunnelCmd.Wait()
			log.Println("SSM tunnel closed")
		}
	}()

	err := uploadFile(absPath, *to)
	if err != nil {
		log.Fatalf("❌ Upload failed: %v", err)
	}
}

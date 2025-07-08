package internal

import (
	"github.com/spf13/cobra"

	"golang.org/x/crypto/ssh"

	"fmt"
	"os"
	"os/exec"
)

func copyToServer(localFile string, remotePath string) error {

	remoteUser := "root"
	remoteHost := "192.168.64.4"

	// scp command
	scpCmd := exec.Command("scp", "-i", os.Getenv("HOME")+"/.ssh/id_rsa", localFile, fmt.Sprintf("%s@%s:%s", remoteUser, remoteHost, remotePath))

	println(scpCmd.String())

	// Forward stdout/stderr
	scpCmd.Stdout = os.Stdout
	scpCmd.Stderr = os.Stderr

	fmt.Println("Transferring file...")
	err := scpCmd.Run()
	if err != nil {
		fmt.Printf("SCP failed: %v\n", err)
		return nil
	}
	fmt.Println("File transferred successfully.")
	return nil
}

func copyAgentToRemote() error {
	agentFile := "./temp/pbdeploy-agent"
	agentRemotePath := "/usr/local/bin/pbdeploy-agent"

	err := copyToServer(agentFile, agentRemotePath)
	if err != nil {
		fmt.Printf("Failed to copy agent to remote: %v\n", err)
		return err
	}

	return nil
}

func createSystemdService() error {
	serviceFile := "./scripts/pbdeploy-agent.service"
	serviceRemotePath := "/etc/systemd/system/"

	err := copyToServer(serviceFile, serviceRemotePath)
	if err != nil {
		fmt.Printf("Failed to copy service file to remote: %v\n", err)
		return err
	}

	return nil
}

func enableService() {
	// Config
	remoteHost := "192.168.64.4:22"
	username := "root"
	privateKeyPath := os.Getenv("HOME") + "/.ssh/id_rsa"
	localFile := "scripts/pbdeploy-agent.service"
	remotePath := "/etc/systemd/system/pbdeploy-agent.service"

	// Load private key
	key, err := os.ReadFile(privateKeyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key: %w", err))
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		panic(fmt.Errorf("failed to parse private key: %w", err))
	}

	config := &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect
	client, err := ssh.Dial("tcp", remoteHost, config)
	if err != nil {
		panic(fmt.Errorf("failed to dial: %w", err))
	}
	defer client.Close()

	// Read local service file
	data, err := os.ReadFile(localFile)
	if err != nil {
		panic(fmt.Errorf("failed to read local service file: %w", err))
	}

	// Start a session to upload
	session, err := client.NewSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		panic(err)
	}

	go func() {
		defer stdin.Close()
		stdin.Write(data)
	}()

	if err := session.Run(fmt.Sprintf("tee %s", remotePath)); err != nil {
		panic(fmt.Errorf("failed to upload file: %w", err))
	}

	fmt.Println("✅ Service file uploaded directly to /etc/systemd/system")

	// Move file with sudo and reload systemd
	session2, err := client.NewSession()
	if err != nil {
		panic(err)
	}
	defer session2.Close()

	cmd := "systemctl daemon-reload && systemctl enable --now pbdeploy-agent"
	if err := session2.Run(cmd); err != nil {
		panic(fmt.Errorf("failed to enable/start service: %w", err))
	}

	fmt.Println("🚀 Service deployed and started successfully.")

}

func Install(cmd *cobra.Command, args []string) {
	copyAgentToRemote()
	createSystemdService()
	enableService()

}

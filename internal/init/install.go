package internal

import (
	"log"

	"github.com/spf13/cobra"

	"golang.org/x/crypto/ssh"

	"fmt"
	"os"
	"os/exec"
)

func copyToServer(localFile string, remotePath string, server string) error {

	remoteUser := "root"
	remoteHost := server

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

func copyAgentToRemote(config DeployConfig) error {
	agentFile := "./temp/pbdeploy-agent-linux-arm64"
	agentRemotePath := "/usr/local/bin/pbdeploy-agent-linux-arm64"

	err := copyToServer(agentFile, agentRemotePath, config.Server)
	if err != nil {
		fmt.Printf("Failed to copy agent to remote: %v\n", err)
		return err
	}

	return nil
}

func createSystemdService(config DeployConfig) error {
	serviceFile := "./scripts/pbdeploy-agent.service"
	serviceRemotePath := "/etc/systemd/system/"

	err := copyToServer(serviceFile, serviceRemotePath, config.Server)
	if err != nil {
		fmt.Printf("Failed to copy service file to remote: %v\n", err)
		return err
	}

	return nil
}

func enableService(yamlConfig DeployConfig) {

	// Config
	remoteHost := yamlConfig.Server + ":22"
	username := "root"
	privateKeyPath := os.Getenv("HOME") + "/.ssh/id_rsa"
	localFile := "scripts/pbdeploy-agent.service"
	remotePath := "/etc/systemd/system/pbdeploy-agent.service"

	err := copyToServer(localFile, remotePath, yamlConfig.Server)

	if err != nil {
		log.Fatal("Error copying " + localFile + " to " + yamlConfig.Server + ".")
	}
	fmt.Println("✅ Service file uploaded directly to /etc/systemd/system")

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

	fmt.Println("🚀 Service enable/started successfully.")

}

func Install(cmd *cobra.Command, args []string) {

	config, err := GetConfig("pbdeploy.yaml")

	if err != nil {
		log.Fatal("Error loading configuration file " + "pbdeploy.yaml")
	}

	err = copyAgentToRemote(config)
	if err != nil {
		log.Print("Error copy agent to remote")
		return
	}

	err = createSystemdService(config)
	if err != nil {
		log.Print("Error create systemd service")
		return
	}

	enableService(config) // cant given the error handling

}

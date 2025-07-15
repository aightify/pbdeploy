package main

import (
	"log"
	"os"
	"testing"
)

// Sample valid YAML content for testing
const sampleYAML = `
repo: git@github.com:example/repo
app_location: "/home/user/app"
ssh_key: "/home/user/.ssh/id_rsa"
app_name: testapp
post_deploy: systemctl restart testapp
branch: main
env:
  DATABASE_URL: postgres://user:pass@localhost:5432/db
  PORT: "8080"
`

func TestGetConfig(t *testing.T) {
	// Create temporary YAML file
	tmpFile, err := os.CreateTemp("", "deploy-config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // clean up

	_, err = tmpFile.WriteString(sampleYAML)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Call the function
	cfg, err := GetConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("GetConfig failed: %v", err)
	}

	// Validate fields
	if cfg.Repo != "git@github.com:example/repo" {
		t.Errorf("Expected Repo, got %s", cfg.Repo)
	}
	if cfg.AppLocation != "/home/user/app" {
		t.Errorf("Expected AppLocation, got %s", cfg.AppLocation)
	}
	if cfg.SshKey != "/home/user/.ssh/id_rsa" {
		t.Errorf("Expected SshKey, got %s", cfg.SshKey)
	}
	if cfg.AppName != "testapp" {
		t.Errorf("Expected AppName, got %s", cfg.AppName)
	}
	if cfg.Branch != "main" {
		t.Errorf("Expected Branch, got %s", cfg.Branch)
	}
	if cfg.Env["DATABASE_URL"] != "postgres://user:pass@localhost:5432/db" {
		t.Errorf("Expected DATABASE_URL, got %s", cfg.Env["DATABASE_URL"])
	}
	if cfg.Env["PORT"] != "8080" {
		t.Errorf("Expected PORT, got %s", cfg.Env["PORT"])
	}
}

func TestCloneRepo_BuildApp(t *testing.T) {
	conf := DeployConfig{
		AppLocation: "./tmp/applocation",
		Repo:        "https://github.com/aightify/spyder-web.git",
		Branch:      "main",
	}

	cloneRepo(&conf)

	//Expect
	//check if the the repo is cloned in AppLocation
	//If it cloned then it is pass, or else it is fail

	_, err := os.Stat(conf.AppLocation)

	if err != nil {
		t.Error("Repo is not cloned")
	}
}

func TestBuildApp(t *testing.T) {
	conf := DeployConfig{
		AppLocation: "./tmp/applocation",
		Repo:        "https://github.com/aightify/spyder-web.git",
		Branch:      "main",
		AppName:     "./bin/mytestapp",
	}

	cloneRepo(&conf)

	//build

	buildApp(&conf)

	//Test the build
	binaryPath := conf.AppLocation + "/" + conf.AppName
	log.Print("looking for " + binaryPath)
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Error("build is not done")
	}

}

func TestGenerateServiceFileForApp(t *testing.T) {
	conf := DeployConfig{
		Branch:    "main",
		AppName:   "myapp",
		ExecStart: "usr/local/bin/myapp",
	}

	serviceFileLocation := "./tmp/"

	GenerateServiceFile(conf, serviceFileLocation)

	//assert to see if the file generated

	//Test the build
	serviceFileLocation = serviceFileLocation + conf.AppName + ".service"

	log.Print("looking for " + serviceFileLocation)

	if _, err := os.Stat(serviceFileLocation); os.IsNotExist(err) {
		t.Error("Service file is not generated", serviceFileLocation)
	}

	//check the service file and see if the data is correctly transformed.

	expected :=
		`[Unit]
Description=myapp

[Service]
ExecStart=usr/local/bin/myapp
Restart=always
User=root

[Install]
WantedBy=multi-user.target`

	transformed, err := os.ReadFile(serviceFileLocation)
	if err != nil {
		t.Error("Cannot read file " + serviceFileLocation)
	}

	if string(transformed) != expected {
		t.Errorf("template output mismatch.\nExpected:\n%s\nGot:\n%s", expected, transformed)
	}

}

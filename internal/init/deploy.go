package internal

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"text/template"

	"gopkg.in/yaml.v2"
)

type DeployConfig struct {
	Repo        string            `yaml:"repo"`
	AppLocation string            `yaml:"app_location"`
	SshKey      string            `yaml:"ssh_key"`
	AppName     string            `yaml:"app_name"`
	Branch      string            `yaml:"branch"`
	ExecStart   string            `yaml:"exec_start"`
	Env         map[string]string `yaml:"env"`
	Server      string            `yaml:"server"`
}

func DeployHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	config, err := GetConfig("pbdeploy.yaml")
	if err != nil {
		http.Error(w, "config error", http.StatusInternalServerError)
		return
	}

	err = cloneRepo(&config)
	if err != nil {
		log.Println("Error building app", err)
		// w.Write([]byte("deployment failed, error cloning the repo"))
		http.Error(w, "deployment failed,error cloning the repo: "+err.Error(), http.StatusInternalServerError)

		return
	}

	err = buildApp(&config)
	if err != nil {
		log.Print("Error building app")
		w.Write([]byte("deployment failed, error building the repo"))
		return
	}

	err = restartAppService(&config)
	if err != nil {
		log.Print("error restart app")
		w.Write([]byte("restart failed , error restart the app"))
		return
	}

	w.Write([]byte("deployment successfully"))
}

func GetConfig(yamlFile string) (DeployConfig, error) {

	data, err := os.ReadFile(yamlFile)

	if err != nil {
		log.Fatal("Error reading the config " + yamlFile + " file.")
		return DeployConfig{}, err
	}

	var cfg DeployConfig

	err = yaml.Unmarshal(data, &cfg)

	if err != nil {
		log.Fatal("Cannot unmarshal yaml file " + yamlFile)
	}

	return cfg, nil
}

func cloneRepo(config *DeployConfig) error {

	// remove config.AppLocation
	log.Print("Removing folder " + config.AppLocation)

	removeAppCmd := exec.Command("rm", "-rf", config.AppLocation)
	removeAppCmd.CombinedOutput()
	// clone into config.AppLocation

	//change working directory to config.AppLocation

	log.Print("Cloning repo into " + config.AppLocation)
	gitCmd := exec.Command("git", "clone", "-b", config.Branch, config.Repo, config.AppLocation)

	cloneOut, err := gitCmd.CombinedOutput()

	if err != nil {
		fmt.Println("failed to clone", err, string(cloneOut))
		return err
	}

	fmt.Println("successfully cloned", gitCmd)
	return nil

}

func buildApp(config *DeployConfig) error {

	log.Print("Building app " + config.AppName)

	buildCmd := exec.Command("go", "build", "-o", config.AppName)

	buildCmd.Dir = config.AppLocation

	log.Print(buildCmd.Args, buildCmd.Dir)

	buildOut, err := buildCmd.CombinedOutput()

	if err != nil {
		fmt.Println("failed to build", err, string(buildOut))
		return err
	}

	fmt.Println("go build successfully output", string(buildOut))
	return nil

}

func GenerateServiceFile(config DeployConfig, serviceFileLocation string) {

	type ServiceTemplateData struct {
		ServiceName string
		ExecCommand string
	}
	//read the template
	serviceTemplate :=
		`[Unit]
Description={{.ServiceName}}

[Service]
ExecStart={{.ExecCommand}}
Restart=always
User=root

[Install]
WantedBy=multi-user.target`

	//transform with config data
	data := ServiceTemplateData{
		ServiceName: config.AppName,
		ExecCommand: config.ExecStart,
	}

	// Create a new template and parse the string
	tmpl, err := template.New("service").Parse(serviceTemplate)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	//write the service file with name <appname>.service

	serviceFile := serviceFileLocation

	file, err := os.Create(serviceFile)
	if err != nil {
		log.Fatal("Error creating service file "+serviceFile+"\n", err)
	}
	defer file.Close()

	// Execute the template, writing the output to standard output
	err = tmpl.Execute(file, data)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

}

func restartAppService(config *DeployConfig) error {

	serviceFile := config.AppName + ".service"

	GenerateServiceFile(*config, "/etc/systemd/system/"+serviceFile)

	restartCmd := exec.Command("systemctl", "restart", serviceFile)

	restartOut, err := restartCmd.CombinedOutput()

	if err != nil {
		fmt.Println("failed to restart", err, string(restartOut))
		return err
	}

	fmt.Println("successfully restarted", restartOut)
	return nil

}

package internal

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func GenerateYml(cmd *cobra.Command, args []string) {
	if _, err := os.Stat("pbdeploy.yaml"); err == nil {
		fmt.Println("already initialized, skipping generation")
		return
	}

	config := map[string]interface{}{
		"server":   "ubuntu@1.2.3.4",
		"ssh_key":  "~/.ssh/id_rsa",
		"repo":     "git@github.com/user/myapp.git",
		"branch":   "main",
		"app_name": "myapp",
		"env": map[string]string{
			"PORT":         "8080",
			"DATABASE_URL": "sqlite://...",
		},
		"post_deploy": "systemctl restart myapp",
		"webhook": map[string]string{
			"enabled": "true",
			"secret":  "abcdef123456",
		},
	}

	data, err := yaml.Marshal(&config)
	if err != nil {
		fmt.Printf("Error generating YAML: %v\n", err)
		return
	}

	err = os.WriteFile("pbdeploy.yaml", data, 0644)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}

	fmt.Println("Data successfully saved to pbdeploy.yaml")

}

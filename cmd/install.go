/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	internal "github.com/aightify/pbdeploy/internal/init"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs pbdeploy-agent remotely via SSH using sudo",
	Long: `This command connects to the target host using the provided SSH credentials,
	 transfers the compiled pbdeploy-agent binary, 
	 and installs it as a systemd service (or init system depending on the target OS).
It ensures the agent runs in the background and starts on system boot.`,

	Run: internal.Install,
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

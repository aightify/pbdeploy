/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	internal "github.com/aightify/pbdeploy/internal/init"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates pbdeploy.yml with deployment descriptor",
	Long: `The init command sets up a new pbdeploy configuration file (pbdeploy.yml)
		defines how the app should be built, deployed,and managed using pbdeploy.

		It typically includes settings such as the application name, remote SSH host,
		deployment directory, build command, and other environment-specific configurations.

		Running 'pbdeploy init' will prompt for key deployment values or generate
		a default pbdeploy.yml file.so that can customize later.`,

	Run: internal.GenerateYml,
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

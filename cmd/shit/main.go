package main

import (
	"fmt"

	"github.com/ShreyeshArangath/shit/internal/git"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var log = logrus.New()

var rootCmd = &cobra.Command{
	Use:   "shit",
	Short: "shreyesh's git tool",
	Long:  "shit is a version control tool inspired by Git.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to shit! Use 'shit --help' for available commands.")
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new repository",
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("path")
		git.Init(path)
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
	}
}

func init() {
	// Add the 'path' flag to the 'init' command
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("path", "p", ".", "Where to create the repository.")
}

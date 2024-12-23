package main

import (
	"fmt"

	"github.com/ShreyeshArangath/shit/internal/git"
	"github.com/spf13/cobra"
)

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
		fmt.Println("Initialized an empty repository.")
	},
}

func main() {
	rootCmd.AddCommand(initCmd)
	rootCmd.Execute()
	git.Init()
}

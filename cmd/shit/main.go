package main

import (
	"fmt"

	"github.com/ShreyeshArangath/shit/internal/git/plumbing"
	"github.com/ShreyeshArangath/shit/internal/git/porcelain"
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

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.AddCommand(porcelain.GetInitCmd())
	rootCmd.AddCommand(plumbing.GetCatFileCmd())
	rootCmd.AddCommand(plumbing.GetHashObjectCmd())
	rootCmd.AddCommand(porcelain.GetLogCmd())
	rootCmd.AddCommand(plumbing.GetLsTreeCmd())
}

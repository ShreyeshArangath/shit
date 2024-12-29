package main

import (
	"fmt"
	"os"

	"github.com/ShreyeshArangath/shit/internal/git/plumbing"
	"github.com/ShreyeshArangath/shit/internal/git/porcelain"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
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

func GetRootCmd() *cobra.Command {
	return rootCmd
}

func init() {
	rootCmd.AddCommand(porcelain.GetInitCmd())
	rootCmd.AddCommand(plumbing.GetCatFileCmd())
	rootCmd.AddCommand(plumbing.GetHashObjectCmd())
	rootCmd.AddCommand(porcelain.GetLogCmd())
	rootCmd.AddCommand(plumbing.GetLsTreeCmd())
	rootCmd.AddCommand(porcelain.GetCheckoutCmd())
	rootCmd.AddCommand(plumbing.GetShowRefCmd())
	// Create the directory to store the documentation
	err := os.MkdirAll("docs", 0755)
	if err != nil {
		log.Fatal(err)
	}

	// Generate the Markdown documentation
	err = doc.GenMarkdownTree(GetRootCmd(), "docs")
	if err != nil {
		log.Fatal(err)
	}
}

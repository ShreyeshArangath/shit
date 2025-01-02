package plumbing

import (
	"fmt"
	"os"

	"github.com/ShreyeshArangath/shit/pkg/models"
	"github.com/spf13/cobra"
)

var revparsecmd = &cobra.Command{
	Use:   "rev-parse",
	Short: "Parse revision identifiers or other object identifiers",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve the flag value
		objType, _ := cmd.Flags().GetString("wyag-type")

		// Validate the wyag-type flag
		validTypes := map[string]bool{
			"blob":   true,
			"commit": true,
			"tag":    true,
			"tree":   true,
		}

		if objType != "" && !validTypes[objType] {
			fmt.Fprintf(cmd.OutOrStderr(), "Invalid type: %s. Must be one of: blob, commit, tag, tree.\n", objType)
			os.Exit(1)
		}

		// Positional argument
		name := args[0]
		repo, err := models.RepoFind(".", true)
		if err != nil {
			log.Fatal(err)
		}
		obj, err := models.ObjectFind(repo, name, objType, true)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(obj)
	},
}

func GetRevParseCmd() *cobra.Command {
	return revparsecmd
}

func init() {
	revparsecmd.Flags().StringP("type", "t", "", "Specify the object type")

}

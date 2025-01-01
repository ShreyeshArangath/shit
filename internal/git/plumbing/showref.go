package plumbing

import (
	"fmt"
	"reflect"

	"github.com/ShreyeshArangath/shit/pkg/models"
	"github.com/spf13/cobra"
)

var ShowRefCmd = &cobra.Command{
	Use:   "show-ref",
	Short: "List references in a shit repository",
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := models.RepoFind(".", true)
		if err != nil {
			log.Fatal(err)
		}
		refList, err := models.ListRef(repo, "")
		if err != nil {
			log.Fatal(err)
		}
		ShowRef(refList, true, "")
	},
}

func GetShowRefCmd() *cobra.Command {
	return ShowRefCmd
}

func ShowRef(refs models.RefMap, withHash bool, prefix string) {
	for k, v := range refs {
		switch value := v.(type) {
		case string: // If the value is a string, it's a leaf
			if withHash {
				fmt.Printf("%s %s%s\n", value, prefixWithSlash(prefix), k)
			} else {
				fmt.Printf("%s %s\n", prefixWithSlash(prefix), k)
			}
		case models.RefMap: // If the value is another map, recurse
			newPrefix := fmt.Sprintf("%s%s%s", prefix, addSlashIfNeeded(prefix), k)
			// cast v to models.RefMap
			v := v.(models.RefMap)
			ShowRef(v, withHash, newPrefix)
		default:
			fmt.Printf("Unexpected type: %v\n", reflect.TypeOf(v))
		}
	}
}

// Helper to append a slash if the prefix is non-empty
func prefixWithSlash(prefix string) string {
	if prefix != "" {
		return prefix + "/"
	}
	return ""
}

// Helper to add a slash if the prefix exists
func addSlashIfNeeded(prefix string) string {
	if prefix != "" {
		return "/"
	}
	return ""
}

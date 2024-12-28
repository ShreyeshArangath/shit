package plumbing

import (
	"fmt"
	"path/filepath"

	"github.com/ShreyeshArangath/shit/pkg/models"
	"github.com/spf13/cobra"
)

var lsTreeCmd = &cobra.Command{
	Use:   "ls-tree",
	Short: "List the contents of a tree object",
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := models.RepoFind(".", true)
		if err != nil {
			fmt.Println(err)
			return
		}
		tree, _ := cmd.Flags().GetString("tree")
		recursive, _ := cmd.Flags().GetBool("recursive")
		err = lsTreeHelper(repo, tree, recursive, "")
		if err != nil {
			fmt.Println(err)
		}
	},
}

// lsTreeHelper recursively lists the contents of a tree object in a shit repository.
//
// Parameters:
//   - repo: A pointer to the Repository object representing the shit repository.
//   - ref: A string representing the reference (branch, tag, or commit SHA) to list the tree from.
//   - recursive: A boolean indicating whether to list tree contents recursively.
//   - prefix: A string representing the prefix to prepend to the file paths in the output.
//
// Returns:
//   - error: An error if any occurs during the operation, otherwise nil.
//
// The function retrieves the tree object corresponding to the given reference and iterates over its items.
// For each item, it determines the object type (tree, blob, commit) and prints its details.
// If the recursive flag is set and the item is a tree, the function calls itself recursively to list the contents of the subtree.
func lsTreeHelper(repo *models.Repository, ref string, recursive bool, prefix string) error {
	sha, err := models.ObjectFind(repo, ref, "tree", false)
	if err != nil {
		return err
	}
	object, err := models.ObjectRead(repo, sha)
	if err != nil {
		return err
	}
	tree, ok := object.(*models.ShitTree)
	if !ok {
		return fmt.Errorf("object is not a tree")
	}
	var objecttype string
	for _, item := range tree.Items {
		if len(item.Mode) == 5 {
			objecttype = item.Mode[0:1]
		} else {
			objecttype = item.Mode[0:2]
		}
		switch objecttype {
		case "04":
			objecttype = "tree"
		case "10":
			objecttype = "blob" // Regular file
		case "12":
			objecttype = "blob" // Symlink file
		case "16":
			objecttype = "commit"
		default:
			return fmt.Errorf("unknown object type: %s", objecttype)
		}
		if !(recursive && objecttype == "tree") {
			mode := fmt.Sprintf("%06s", item.Mode)
			fmt.Printf("%s %s %s\t%s\n",
				mode,
				objecttype,
				item.Sha,
				filepath.Join(prefix, item.Path),
			)
		} else {
			lsTreeHelper(repo, item.Sha, recursive, filepath.Join(prefix, item.Path))
		}
	}
	return nil
}

func GetLsTreeCmd() *cobra.Command {
	return lsTreeCmd
}

func init() {
	lsTreeCmd.Flags().StringP("tree", "t", "", "Specify the tree object")
	lsTreeCmd.Flags().BoolP("recursive", "r", false, "Recurse into subtrees")
	lsTreeCmd.MarkFlagRequired("tree")
}

package porcelain

import (
	"os"
	"path/filepath"

	"github.com/ShreyeshArangath/shit/pkg/models"
	"github.com/ShreyeshArangath/shit/pkg/utils"
	"github.com/spf13/cobra"
)

var checkoutCmd = &cobra.Command{
	Use:   "checkout",
	Short: "Checkout a branch or paths to the working tree",
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := models.RepoFind(".", true)
		if err != nil {
			log.Fatal(err)
		}
		commit, _ := cmd.Flags().GetString("commit")
		path, _ := cmd.Flags().GetString("path")

		name, err := models.ObjectFind(repo, commit, "commit", false)
		if err != nil {
			log.Fatal(err)
		}
		object, err := models.ObjectRead(repo, name)
		if err != nil {
			log.Fatal(err)
		}
		commitobj, ok := object.(*models.ShitCommit)
		if !ok {
			log.Fatalf("Object %s is not of the type commit", name)
		}
		object, err = models.ObjectRead(repo, commitobj.CommitMetadata.GetTree())
		if err != nil {
			log.Fatal(err)
		}
		tree, ok := object.(*models.ShitTree)
		if !ok {
			log.Fatalf("Object %s is not of the type tree", commitobj.CommitMetadata.GetTree())
		}
		exists, err := utils.PathExists(path)
		if err != nil {
			log.Fatal(err)
		}
		if !exists {
			err = os.MkdirAll(path, 0755)
		} else {
			isDir, err := utils.IsDirEmpty(path)
			if err != nil {
				log.Fatal(err)
			}
			if !isDir {
				log.Fatalf("Path %s is not an empty directory", path)
			}
		}
		checkouthelper(repo, tree, path)
	},
}

func checkouthelper(repo *models.Repository, tree *models.ShitTree, path string) error {
	for _, item := range tree.Items {
		obj, err := models.ObjectRead(repo, item.Sha)
		if err != nil {
			return err
		}
		dest := filepath.Join(path, item.Path)
		if obj.GetType() == "tree" {
			err = os.MkdirAll(dest, 0755)
			if err != nil {
				return err
			}
			err = checkouthelper(repo, obj.(*models.ShitTree), dest)
			if err != nil {
				return err
			}
		} else if obj.GetType() == "blob" {
			// TODO: Support symlinks (identified by mode 12****)
			err = os.WriteFile(dest, obj.(*models.ShitBlob).Data, 0644)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func init() {
	checkoutCmd.Flags().StringP("commit", "c", "", "Commit to checkout")
	checkoutCmd.Flags().StringP("path", "p", "", "The empty directory to checkout on")
}

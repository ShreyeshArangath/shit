package plumbing

import (
	"github.com/ShreyeshArangath/shit/pkg/models"
	"github.com/spf13/cobra"
)

var hashObjectCmd = &cobra.Command{
	Use:   "hash-object",
	Short: "Compute object ID and optionally creates a blob from a file",
	Run: func(cmd *cobra.Command, args []string) {
		write, _ := cmd.Flags().GetBool("write")
		path, _ := cmd.Flags().GetString("path")
		objectType, _ := cmd.Flags().GetString("type")
		err := hashObject(objectType, path, write)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func GetHashObjectCmd() *cobra.Command {
	return hashObjectCmd
}

func hashObject(objectype string, path string, write bool) error {
	var repo *models.Repository
	var err error
	if write {
		repo, err = models.RepoFind(path, true)
		if err != nil {
			return err
		}
	} else {
		repo = &models.Repository{}
	}
	sha, err := models.ObjectHash(repo, objectype, path)
	if err != nil {
		return err
	}
	log.Println(sha)
	return nil
}

func init() {
	hashObjectCmd.Flags().StringP("path", "p", ".", "Specify the path")
	hashObjectCmd.Flags().StringP("type", "t", "blob", "Specify the type")
	hashObjectCmd.Flags().BoolP("write", "w", false, "Actually write the object into the database")
}

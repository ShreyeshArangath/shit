package porcelain

import (
	"fmt"

	"github.com/ShreyeshArangath/shit/internal/git/plumbing"
	"github.com/ShreyeshArangath/shit/pkg/models"
	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "List and create shit tag objects",
	Run: func(cmd *cobra.Command, args []string) {
		repo, err := models.RepoFind(".", true)
		tag, err := cmd.Flags().GetString("tag")
		object, err := cmd.Flags().GetString("object")
		create, err := cmd.Flags().GetBool("create")
		if err != nil {
			log.Fatal(err)
		}
		if tag != "" {
			err = tagCreateHelper(repo, tag, object, create)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			refs, err := models.ListRef(repo, "")
			if err != nil {
				log.Fatal(err)
			}
			if tags, ok := refs["tags"].(models.RefMap); ok {
				plumbing.ShowRef(tags, true, "")
			} else {
				log.Fatal("Invalid type assertion for refs[\"tags\"]")
			}
		}
	},
}

func GetTagCmd() *cobra.Command {
	return tagCmd
}

func init() {
	tagCmd.Flags().StringP("tag", "t", "", "Tag name")
	tagCmd.Flags().BoolP("create", "c", false, "Create a new tag")
	tagCmd.Flags().StringP("object", "o", "HEAD", "Object that the tag points to")
}

func tagCreateHelper(repo *models.Repository, tag string, ref string, create bool) error {
	sha, err := models.ObjectFind(repo, ref, "commit", false)
	if err != nil {
		return err
	}
	if create {
		tagmetadata := models.CreateShitTagMetadataFromAttr(
			sha,
			"commit",
			tag,
			"Shreyesh Arangath <>",
			"Tagged by shit. Cannot be customized.",
		)
		tagObj := models.ShitTag{TagMetadata: tagmetadata}
		tagSha, err := models.ObjectWrite(&tagObj, repo)
		if err != nil {
			return err
		}
		err = models.CreateRef(repo, fmt.Sprintf("tags/%s", tag), tagSha)
		if err != nil {
			return err
		}
	} else {
		models.CreateRef(repo, fmt.Sprintf("tags/%s", tag), sha)
		if err != nil {
			return err
		}
	}
	return nil
}

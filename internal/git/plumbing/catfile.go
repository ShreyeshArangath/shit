package plumbing

import (
	"github.com/ShreyeshArangath/shit/pkg/models"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var log = logrus.New()

var catFileCmd = &cobra.Command{
	Use:   "cat-file",
	Short: "Provide content of repository objects",
	Run: func(cmd *cobra.Command, args []string) {
		objectType, _ := cmd.Flags().GetString("type")
		object, _ := cmd.Flags().GetString("object")
		catfileout, err := CatFile(object, objectType)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(catfileout)
	},
}

func GetCatFileCmd() *cobra.Command {
	return catFileCmd
}

func CatFile(sha string, objectType string) (string, error) {
	repo, err := models.RepoFind(".", true)
	if err != nil {
		return "", err
	}
	sha, err = models.ObjectFind(repo, sha, objectType, true)
	if err != nil {
		return "", err
	}
	object, err := models.ObjectRead(repo, sha)
	if err != nil {
		return "", err
	}
	bytes, err := object.Serialize(repo)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func init() {
	catFileCmd.Flags().StringP("type", "t", "", "Specify the type")
	catFileCmd.Flags().StringP("object", "o", "", "Specify the object")
	catFileCmd.MarkFlagRequired("type")
	catFileCmd.MarkFlagRequired("object")
}

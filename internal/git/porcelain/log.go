package porcelain

import (
	"strings"

	"github.com/ShreyeshArangath/shit/pkg/models"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Display history of a given commit",
	Run: func(cmd *cobra.Command, args []string) {
		commit, _ := cmd.Flags().GetString("commit")
		repo, err := models.RepoFind(".", true)
		if err != nil {
			log.Fatal(err)
		}
		object, err := models.ObjectFind(repo, commit, "commit", true)
		if err != nil {
			log.Fatal(err)
		}
		seen := mapset.NewSet[string]()
		err = loghelper(repo, object, seen)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func loghelper(repo *models.Repository, sha string, seen mapset.Set[string]) error {
	// Example usage of the set
	if seen.Contains(sha) {
		return nil
	}
	seen.Add(sha)
	var commit *models.ShitCommit
	obj, err := models.ObjectRead(repo, sha)
	if err != nil {
		return err
	}
	commit, ok := obj.(*models.ShitCommit)
	if !ok {
		log.Fatalf("Object %s is not a commit", sha)
		return nil
	}
	if err != nil {
		return err
	}
	short_hash := sha[:8]
	message := strings.TrimSpace(commit.CommitMetadata.GetMessage())
	if strings.Contains(message, "\n") {
		message = strings.Split(message, "\n")[0]
	}
	log.Printf(" c_%s [label=\"%s: %s\"]", sha, short_hash, message)
	parents := commit.CommitMetadata.GetParent()
	if len(parents) > 0 {
		for _, parent := range parents {
			log.Printf(" c_%s -> c_%s", sha, parent)
			err = loghelper(repo, parent, seen)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func init() {
	logCmd.Flags().StringP("commit", "c", "HEAD", "Commit to start at")
}

func GetLogCmd() *cobra.Command {
	return logCmd
}

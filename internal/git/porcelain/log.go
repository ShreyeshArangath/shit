package porcelain

import (
	"fmt"
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

// loghelper traverses the commit history of a repository starting from a given commit SHA,
// printing the commit details and their parent relationships in a specific format.
// It uses a set to keep track of seen commits to avoid processing the same commit multiple times.
//
// Parameters:
// - repo: A pointer to the Repository object representing the Git repository.
// - sha: The SHA-1 hash of the commit to start traversing from.
// - seen: A set of strings to keep track of already processed commits.
//
// Returns:
// - An error if any occurs during the processing of commits, otherwise nil.
func loghelper(repo *models.Repository, sha string, seen mapset.Set[string]) error {
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
		log.Fatalf("Object %s is not of the type commit", sha)
		return nil
	}
	if err != nil {
		return err
	}
	shortHash := sha[:8]
	message := strings.TrimSpace(commit.CommitMetadata.GetMessage())
	if strings.Contains(message, "\n") {
		message = strings.Split(message, "\n")[0]
	}
	fmt.Printf(" c_%s [label=\"%s: %s\"]\n", sha, shortHash, message)
	parents := commit.CommitMetadata.GetParent()
	if len(parents) > 0 {
		for _, parent := range parents {
			fmt.Printf(" c_%s -> c_%s\n", sha, parent)
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

package porcelain

import "github.com/spf13/cobra"

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Display history of a given commit",
	Run: func(cmd *cobra.Command, args []string) {
		commit, _ := cmd.Flags().GetString("commit")
		err := loghelper(commit)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func loghelper(commit string) error {

	return nil
}

func init() {
	logCmd.Flags().StringP("commit", "c", "HEAD", "Commit to start at")
}

func GetLogCmd() *cobra.Command {
	return logCmd
}

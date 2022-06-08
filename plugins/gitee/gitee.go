package main

import (
	"github.com/apache/incubator-devlake/plugins/gitee/impl"
	"github.com/apache/incubator-devlake/runner"
	"github.com/spf13/cobra"
)

var PluginEntry impl.Gitee //nolint

func main() {
	githubCmd := &cobra.Command{Use: "gitee"}
	owner := githubCmd.Flags().StringP("owner", "o", "", "github owner")
	repo := githubCmd.Flags().StringP("repo", "r", "", "github repo")
	_ = githubCmd.MarkFlagRequired("owner")
	_ = githubCmd.MarkFlagRequired("repo")

	githubCmd.Run = func(cmd *cobra.Command, args []string) {
		runner.DirectRun(cmd, args, PluginEntry, map[string]interface{}{
			"owner": *owner,
			"repo":  *repo,
		})
	}
	runner.RunCmd(githubCmd)
}

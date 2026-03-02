// Package cmd provides command-line interface commands for the gh-godo extension.
package cmd

import (
	"fmt"

	"github.com/lakruzz/gh-godo/cmd/mkissue"
	"github.com/spf13/cobra"
)

var (
	issueFile  string
	branchName string
	gistID     string
	repoName   string
)

var mkissueCmd = &cobra.Command{
	Use:   "mkissue",
	Short: "Create a GitHub issue from a markdown file",
	Long: `Create a GitHub issue from a markdown file with frontmatter support.
The markdown file should contain YAML frontmatter with metadata and a markdown body.

Usage variants:
  godo mkissue --file <file> [--branch <branch>] [--repo <owner/repo>]
  godo mkissue --file <file> [--gist <gist-id>]

Rules:
  --file is always required
  --branch is optional (defaults to the repo's default branch when used with --repo)
  --gist and --repo are mutually exclusive
  --branch is not valid with --gist`,
	RunE: func(_ *cobra.Command, _ []string) error {
		// Validate that branch and gist are not both specified
		if branchName != "" && gistID != "" {
			return fmt.Errorf("cannot use both --branch and --gist flags together")
		}
		// Validate that repo and gist are not both specified
		if repoName != "" && gistID != "" {
			return fmt.Errorf("cannot use both --repo and --gist flags together")
		}
		// Call the original mkissue logic with the file path, branch, gist, and repo
		return mkissue.RunWithFile(issueFile, branchName, gistID, repoName)
	},
}

func init() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.AddCommand(mkissueCmd)

	// Define flags for mkissue command
	mkissueCmd.Flags().StringVarP(&issueFile, "file", "f", "", "Path to the markdown file containing issue content (required)")
	mkissueCmd.Flags().StringVarP(&branchName, "branch", "b", "", "Branch name to get the file from (optional)")
	mkissueCmd.Flags().StringVarP(&gistID, "gist", "g", "", "Gist ID to get the file from (optional)")
	mkissueCmd.Flags().StringVarP(&repoName, "repo", "r", "", "Repository to get the file from, in owner/repo format (optional)")
	_ = mkissueCmd.MarkFlagRequired("file")
}

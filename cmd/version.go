package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	gitRevision = "HEAD"
	buildStamp  = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Git Revision: %s\n", gitRevision)
		fmt.Printf("UTC Build Time: %s\n", buildStamp)
		os.Exit(0)
	},
}

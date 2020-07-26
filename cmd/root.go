package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var (
	browserSource string
	browserChannel string
	driverVersion string
	noCache bool
	testsDir string
	skipTests bool
	tags []string

	rootCmd  = &cobra.Command{
		Use:   "images",
		Short: "images is a tool to build Docker images with browsers",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}
)

func initFlags() {
	rootCmd.PersistentFlags().StringSliceVarP(&tags, "tag", "t", []string{}, "image tag")
	rootCmd.PersistentFlags().StringVarP(&browserSource, "browser", "b", "", "browser APT package version, package file path, package file URL")
	rootCmd.PersistentFlags().StringVarP(&driverVersion, "driver-version", "d", "", "webdriver version")
	rootCmd.PersistentFlags().StringVarP(&browserChannel, "channel", "c", "default", "browser channel")
	rootCmd.PersistentFlags().BoolVarP(&noCache, "no-cache", "n", false, "do not use Docker cache")
	rootCmd.PersistentFlags().StringVar(&testsDir, "tests-dir", "", "directory with tests")
	rootCmd.PersistentFlags().BoolVar(&skipTests, "skip-tests", false, "skip test execution results")
}

func init() {
	initFlags()
	rootCmd.AddCommand(chromeCmd)
	rootCmd.AddCommand(versionCmd)
}

func Execute() {
	if _, err := rootCmd.ExecuteC(); err != nil {
		os.Exit(1)
	}
}

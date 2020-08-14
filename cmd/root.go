package cmd

import (
	"github.com/aerokube/selenoid-images/build"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

var (
	browserSource  string
	browserChannel string
	driverVersion  string
	noCache        bool
	testsDir       string
	test           bool
	ignoreTests    bool
	push           bool
	tags           []string

	rootCmd = &cobra.Command{
		Use:           "images",
		Short:         "images is a tool for building Docker images with browsers",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}
)

func initFlags() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("get current dir: %v", err)
	}
	defaultTestsDir := filepath.Join(cwd, "../selenoid-container-tests")

	rootCmd.PersistentFlags().StringSliceVarP(&tags, "tag", "t", []string{}, "image tag")
	rootCmd.PersistentFlags().StringVarP(&browserSource, "browser", "b", "", "browser APT package version, package file path, package file URL")
	rootCmd.PersistentFlags().StringVarP(&driverVersion, "driver-version", "d", build.LatestVersion, "webdriver version")
	rootCmd.PersistentFlags().StringVarP(&browserChannel, "channel", "c", "default", "browser channel")
	rootCmd.PersistentFlags().BoolVarP(&noCache, "no-cache", "n", false, "do not use Docker cache")
	rootCmd.PersistentFlags().StringVar(&testsDir, "tests-dir", defaultTestsDir, "directory with tests")
	rootCmd.PersistentFlags().BoolVar(&test, "test", false, "run tests")
	rootCmd.PersistentFlags().BoolVar(&ignoreTests, "ignore-tests", false, "continue to run even if tests failed")
	rootCmd.PersistentFlags().BoolVarP(&push, "push", "p", false, "push image to Docker registry")
}

func init() {
	initFlags()
	rootCmd.AddCommand(chromeCmd)
	rootCmd.AddCommand(firefoxCmd)
	rootCmd.AddCommand(operaCmd)
	rootCmd.AddCommand(yandexCmd)
	rootCmd.AddCommand(versionCmd)
}

func Execute() {
	if _, err := rootCmd.ExecuteC(); err != nil {
		log.Fatalf("command error: %v", err)
	}
}

package cmd

import (
	"github.com/aerokube/selenoid-images/build"
	"github.com/spf13/cobra"
)

var (
	chromeCmd  = &cobra.Command{
		Use:   "chrome",
		Short: "build Chrome image",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := build.Requirements{
				BrowserSource: build.BrowserSource(browserSource),
				BrowserChannel: browserChannel,
				DriverVersion: driverVersion,
				NoCache: noCache,
				TestsDir: testsDir,
				SkipTests: skipTests,
				Tags: tags,
			}
			chrome := &build.Chrome{req}
			return chrome.Build()
		},
	}
)

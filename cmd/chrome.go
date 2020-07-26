package cmd

import (
	"github.com/aerokube/selenoid-images/build"
	"github.com/spf13/cobra"
)

var (
	chromeCmd = &cobra.Command{
		Use:   "chrome",
		Short: "build Chrome image",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := build.Requirements{
				BrowserSource:  build.BrowserSource(browserSource),
				BrowserChannel: browserChannel,
				DriverVersion:  driverVersion,
				NoCache:        noCache,
				TestsDir:       testsDir,
				RunTests:       test,
				Tags:           tags,
				PushImage:      push,
			}
			chrome := &build.Chrome{req}
			return chrome.Build()
		},
	}
)

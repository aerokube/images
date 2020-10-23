package cmd

import (
	"github.com/aerokube/images/build"
	"github.com/spf13/cobra"
)

var (
	edgeCmd = &cobra.Command{
		Use:   "edge",
		Short: "build Microsoft Edge image",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := build.Requirements{
				BrowserSource:  build.BrowserSource(browserSource),
				BrowserChannel: browserChannel,
				DriverVersion:  driverVersion,
				NoCache:        noCache,
				TestsDir:       testsDir,
				RunTests:       test,
				IgnoreTests:    ignoreTests,
				Tags:           tags,
				PushImage:      push,
			}
			edge := &build.Edge{req}
			return edge.Build()
		},
	}
)

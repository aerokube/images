package cmd

import (
	"github.com/aerokube/images/build"
	"github.com/spf13/cobra"
)

var (
	chromiumCmd = &cobra.Command{
		Use:   "chromium",
		Short: "build Chromium image",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := build.Requirements{
				BrowserSource: build.BrowserSource(browserSource),
				NoCache:       noCache,
				TestsDir:      testsDir,
				RunTests:      test,
				IgnoreTests:   ignoreTests,
				Tags:          tags,
				PushImage:     push,
			}
			chromium := &build.Chromium{req}
			return chromium.Build()
		},
	}
)

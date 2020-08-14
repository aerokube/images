package cmd

import (
	"github.com/aerokube/selenoid-images/build"
	"github.com/spf13/cobra"
)

var (
	yandexCmd = &cobra.Command{
		Use:   "yandex",
		Short: "build Yandex.Browser image",
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
			yandexBrowser := &build.YandexBrowser{req}
			return yandexBrowser.Build()
		},
	}
)

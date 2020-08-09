package build

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	operaDriverBinary = "operadriver_linux64/operadriver"
)

type Opera struct {
	Requirements
}

func (o *Opera) Build() error {

	// Build dev image
	devDestDir, err := tmpDir()
	if err != nil {
		return fmt.Errorf("create dev temporary dir: %v", err)
	}

	srcDir := "opera/apt"
	pkgSrcPath, pkgVersion, err := o.BrowserSource.Prepare()
	if err != nil {
		return fmt.Errorf("invalid browser source: %v", err)
	}

	if pkgSrcPath != "" {
		srcDir = "opera/local"
		pkgDestPath := filepath.Join(devDestDir, "opera.deb")
		err = os.Rename(pkgSrcPath, pkgDestPath)
		if err != nil {
			return fmt.Errorf("move package: %v", err)
		}
	}

	pkgTagVersion := extractVersion(pkgVersion)
	devImageTag := fmt.Sprintf("selenoid/dev_opera:%s", pkgTagVersion)
	devImageRequirements := Requirements{NoCache: o.NoCache, Tags: []string{devImageTag}}
	devImage, err := NewImage(srcDir, devDestDir, devImageRequirements)
	if err != nil {
		return fmt.Errorf("init dev image: %v", err)
	}
	devBuildArgs := []string{fmt.Sprintf("VERSION=%s", pkgVersion)}
	devBuildArgs = append(devBuildArgs, o.channelToBuildArgs()...)
	devImage.BuildArgs = devBuildArgs
	if pkgSrcPath != "" {
		devImage.FileServer = true
	}

	err = devImage.Build()
	if err != nil {
		return fmt.Errorf("build dev image: %v", err)
	}

	// Build main image
	destDir, err := tmpDir()
	if err != nil {
		return fmt.Errorf("create temporary dir: %v", err)
	}

	image, err := NewImage("opera", destDir, o.Requirements)
	if err != nil {
		return fmt.Errorf("init image: %v", err)
	}
	image.BuildArgs = append(image.BuildArgs, fmt.Sprintf("VERSION=%s", pkgTagVersion))

	driverVersion, err := o.downloadOperaDriver(image.Dir)
	if err != nil {
		return fmt.Errorf("failed to download operadriver: %v", err)
	}
	image.Labels = []string{fmt.Sprintf("driver=operadriver:%s", driverVersion)}

	err = image.Build()
	if err != nil {
		return fmt.Errorf("build image: %v", err)
	}

	err = image.Test(o.TestsDir, "opera", pkgTagVersion)
	if err != nil {
		return fmt.Errorf("test image: %v", err)
	}

	err = image.Push()
	if err != nil {
		return fmt.Errorf("push image: %v", err)
	}

	return nil
}

func (o *Opera) channelToBuildArgs() []string {
	switch o.BrowserChannel {
	case "beta":
		return []string{"PACKAGE=opera-beta"}
	case "dev":
		return []string{"PACKAGE=opera-developer"}
	default:
		return []string{}
	}
}

func (o *Opera) downloadOperaDriver(dir string) (string, error) {
	version := o.DriverVersion
	if version == LatestVersion {
		v, err := latestGithubRelease("operasoftware/operachromiumdriver")
		if err != nil {
			return "", fmt.Errorf("latest Operadriver version: %v", err)
		}
		version = strings.TrimPrefix("v.", v)
	}

	u := fmt.Sprintf("https://github.com/operasoftware/operachromiumdriver/releases/download/v.%s/operadriver_linux64.zip", version)
	_, err := downloadDriver(u, operaDriverBinary, dir)
	if err != nil {
		return "", fmt.Errorf("download Operadriver: %v", err)
	}
	return version, nil
}

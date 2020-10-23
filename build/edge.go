package build

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	msedgeDriverBinary = "msedgedriver"
)

type Edge struct {
	Requirements
}

func (c *Edge) Build() error {

	pkgSrcPath, pkgVersion, err := c.BrowserSource.Prepare()
	if err != nil {
		return fmt.Errorf("invalid browser source: %v", err)
	}

	pkgTagVersion := extractVersion(pkgVersion)

	// Build dev image
	devDestDir, err := tmpDir()
	if err != nil {
		return fmt.Errorf("create dev temporary dir: %v", err)
	}

	srcDir := "edge/apt"

	if pkgSrcPath != "" {
		srcDir = "edge/local"
		pkgDestPath := filepath.Join(devDestDir, "microsoft-edge.deb")
		err = os.Rename(pkgSrcPath, pkgDestPath)
		if err != nil {
			return fmt.Errorf("move package: %v", err)
		}
	}

	devImageTag := fmt.Sprintf("selenoid/dev_edge:%s", pkgTagVersion)
	devImageRequirements := Requirements{NoCache: c.NoCache, Tags: []string{devImageTag}}
	devImage, err := NewImage(srcDir, devDestDir, devImageRequirements)
	if err != nil {
		return fmt.Errorf("init dev image: %v", err)
	}
	devBuildArgs := []string{fmt.Sprintf("VERSION=%s", pkgVersion)}
	devBuildArgs = append(devBuildArgs, c.channelToBuildArgs()...)
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

	image, err := NewImage("edge", destDir, c.Requirements)
	if err != nil {
		return fmt.Errorf("init image: %v", err)
	}
	image.BuildArgs = append(image.BuildArgs, fmt.Sprintf("VERSION=%s", pkgTagVersion))

	driverVersion, err := c.downloadMSEdgeDriver(image.Dir)
	if err != nil {
		return fmt.Errorf("failed to download msedgedriver: %v", err)
	}
	image.Labels = []string{fmt.Sprintf("driver=msedgedriver:%s", driverVersion)}

	err = image.Build()
	if err != nil {
		return fmt.Errorf("build image: %v", err)
	}

	err = image.Test(c.TestsDir, "MicrosoftEdge", pkgTagVersion)
	if err != nil {
		return fmt.Errorf("test image: %v", err)
	}

	err = image.Push()
	if err != nil {
		return fmt.Errorf("push image: %v", err)
	}

	return nil
}

func (c *Edge) channelToBuildArgs() []string {
	switch c.BrowserChannel {
	case "beta":
		return []string{"PACKAGE=microsoft-edge-beta", "INSTALL_DIR=msedge-beta"}
	case "dev":
		return []string{"PACKAGE=microsoft-edge-dev", "INSTALL_DIR=msedge-dev"}
	default:
		return []string{}
	}
}

func (c *Edge) downloadMSEdgeDriver(dir string) (string, error) {
	version := c.DriverVersion
	// Full driver versions list can be fetched as XML from https://msedgedriver.azureedge.net/
	u := fmt.Sprintf("https://msedgewebdriverstorage.blob.core.windows.net/edgewebdriver/%s/edgedriver_linux64.zip", version)
	_, err := downloadDriver(u, msedgeDriverBinary, dir)
	if err != nil {
		return "", fmt.Errorf("download msedgedriver: %v", err)
	}
	return version, nil
}

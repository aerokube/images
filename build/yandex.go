package build

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	yandexDriverBinary = "yandexdriver"
)

type YandexBrowser struct {
	Requirements
}

func (yb *YandexBrowser) Build() error {

	// Build dev image
	devDestDir, err := tmpDir()
	if err != nil {
		return fmt.Errorf("create dev temporary dir: %v", err)
	}

	srcDir := "yandex/apt"
	pkgSrcPath, pkgVersion, err := yb.BrowserSource.Prepare()
	if err != nil {
		return fmt.Errorf("invalid browser source: %v", err)
	}

	if pkgSrcPath != "" {
		srcDir = "yandex/local"
		pkgDestDir := filepath.Join(devDestDir, srcDir)
		err := os.MkdirAll(pkgDestDir, 0755)
		if err != nil {
			return fmt.Errorf("create %v temporary dir: %v", pkgDestDir, err)
		}
		pkgDestPath := filepath.Join(pkgDestDir, "yandex-browser.deb")
		err = os.Rename(pkgSrcPath, pkgDestPath)
		if err != nil {
			return fmt.Errorf("move package: %v", err)
		}
	}

	pkgTagVersion := extractVersion(pkgVersion)
	devImageTag := fmt.Sprintf("selenoid/dev_yandex:%s", pkgTagVersion)
	devImageRequirements := Requirements{NoCache: yb.NoCache, Tags: []string{devImageTag}}
	devImage, err := NewImage(srcDir, devDestDir, devImageRequirements)
	if err != nil {
		return fmt.Errorf("init dev image: %v", err)
	}
	devBuildArgs := []string{fmt.Sprintf("VERSION=%s", pkgVersion)}
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

	image, err := NewImage("yandex", destDir, yb.Requirements)
	if err != nil {
		return fmt.Errorf("init image: %v", err)
	}
	image.BuildArgs = append(image.BuildArgs, fmt.Sprintf("VERSION=%s", pkgTagVersion))

	driverVersion, err := yb.downloadYandexDriver(image.Dir)
	if err != nil {
		return fmt.Errorf("failed to download yandexdriver: %v", err)
	}
	image.Labels = []string{fmt.Sprintf("driver=yandexdriver:%s", driverVersion)}

	err = image.Build()
	if err != nil {
		return fmt.Errorf("build image: %v", err)
	}

	err = image.Test(yb.TestsDir, "yandex", pkgTagVersion)
	if err != nil {
		return fmt.Errorf("test image: %v", err)
	}

	err = image.Push()
	if err != nil {
		return fmt.Errorf("push image: %v", err)
	}

	return nil
}

func (yb *YandexBrowser) downloadYandexDriver(dir string) (string, error) {
	version := yb.DriverVersion
	if version == LatestVersion {
		a, err := latestGithubLinuxAsset("yandex/YandexDriver")
		if err != nil {
			return "", fmt.Errorf("latest yandexdriver linux asset: %v", err)
		}
		version = strings.Split(a, "-")[1]
	}

	bv := buildVersion(version)
	u := fmt.Sprintf("https://github.com/yandex/YandexDriver/releases/download/v%s-stable/yandexdriver-%s-linux.zip", bv, version)
	_, err := downloadDriver(u, yandexDriverBinary, dir)
	if err != nil {
		return "", fmt.Errorf("download Yandexdriver: %v", err)
	}
	return version, nil
}

package build

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

const (
	chromeDriverBinary = "chromedriver"
)

const (
	chromeDriverPath = "chromedriver-linux64/"
)

type Chrome struct {
	Requirements
}

func (c *Chrome) Build() error {

	pkgSrcPath, pkgVersion, err := c.BrowserSource.Prepare()
	if err != nil {
		return fmt.Errorf("invalid browser source: %v", err)
	}

	pkgTagVersion := extractVersion(pkgVersion)

	driverVersion, err := c.parseChromeDriverVersion(pkgTagVersion)
	if err != nil {
		return fmt.Errorf("parse chromedriver version: %v", err)
	}

	// Build dev image
	devDestDir, err := tmpDir()
	if err != nil {
		return fmt.Errorf("create dev temporary dir: %v", err)
	}

	srcDir := "chrome/apt"

	if pkgSrcPath != "" {
		srcDir = "chrome/local"
		pkgDestDir := filepath.Join(devDestDir, srcDir)
		err := os.MkdirAll(pkgDestDir, 0755)
		if err != nil {
			return fmt.Errorf("create %v temporary dir: %v", pkgDestDir, err)
		}
		pkgDestPath := filepath.Join(pkgDestDir, "google-chrome.deb")
		err = os.Rename(pkgSrcPath, pkgDestPath)
		if err != nil {
			return fmt.Errorf("move package: %v", err)
		}
	}

	devImageTag := fmt.Sprintf("selenoid/dev_chrome:%s", pkgTagVersion)
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

	image, err := NewImage("chrome", destDir, c.Requirements)
	if err != nil {
		return fmt.Errorf("init image: %v", err)
	}
	image.BuildArgs = append(image.BuildArgs, fmt.Sprintf("VERSION=%s", pkgTagVersion))

	err = c.downloadChromeDriver(image.Dir, driverVersion)
	if err != nil {
		return fmt.Errorf("failed to download chromedriver: %v", err)
	}
	image.Labels = []string{fmt.Sprintf("driver=chromedriver:%s", driverVersion)}

	err = image.Build()
	if err != nil {
		return fmt.Errorf("build image: %v", err)
	}

	err = image.Test(c.TestsDir, "chrome", pkgTagVersion)
	if err != nil {
		return fmt.Errorf("test image: %v", err)
	}

	err = image.Push()
	if err != nil {
		return fmt.Errorf("push image: %v", err)
	}

	return nil
}

func (c *Chrome) channelToBuildArgs() []string {
	switch c.BrowserChannel {
	case "beta":
		return []string{"PACKAGE=google-chrome-beta", "INSTALL_DIR=chrome-beta"}
	case "dev":
		return []string{"PACKAGE=google-chrome-unstable", "INSTALL_DIR=chrome-unstable"}
	default:
		return []string{}
	}
}

func (c *Chrome) parseChromeDriverVersion(pkgVersion string) (string, error) {
	version := c.DriverVersion
	if version == LatestVersion {
		const baseUrl = "https://chromedriver.storage.googleapis.com/"
		v, err := c.getLatestChromeDriver(baseUrl, pkgVersion)
		if err != nil {
			return "", err
		}
		version = v
	}
	return version, nil
}

func moveFile(srcPath, dstPath string) error {
    inputFile, err := os.Open(srcPath)
    if err != nil {
        return fmt.Errorf("Couldn't open source file: %s", err)
    }
    // Создаём нужный файл
    outputFile, err := os.Create(dstPath)
    if err != nil {
        inputFile.Close()
        return fmt.Errorf("Couldn't open dest file: %s", err)
    }
    defer outputFile.Close()
    // Копируем содержимое
    _, err = io.Copy(outputFile, inputFile)
    inputFile.Close()
    if err != nil {
        return fmt.Errorf("Writing to output file failed: %s", err)
    }
    // Удаляем исходный файл, если не было ошибок
    err = os.Remove(srcPath)
    if err != nil {
        return fmt.Errorf("Failed removing original file: %s", err)
    }
    return nil
}

func (c *Chrome) downloadChromeDriver(dir string, version string) error {
	u := fmt.Sprintf("https://edgedl.me.gvt1.com/edgedl/chrome/chrome-for-testing/%s/linux64/chromedriver-linux64.zip", version)
	_, err := downloadDriver(u, chromeDriverPath + chromeDriverBinary, dir)
	if err != nil {
		return fmt.Errorf("download chromedriver: %v", err)
	}
	return moveFile(dir + "/" + chromeDriverPath + chromeDriverBinary, dir + "/" + chromeDriverBinary)
}

func (c *Chrome) getLatestChromeDriver(baseUrl string, pkgVersion string) (string, error) {
	fetchVersion := func(url string) (string, error) {
		data, err := sendGet(url)
		if err != nil {
			return "", fmt.Errorf("read chromedriver version: %v", err)
		}
		return string(data), nil
	}

	if c.BrowserChannel != "dev" {
		chromeBuildVersion := buildVersion(pkgVersion)
		u := baseUrl + fmt.Sprintf("LATEST_RELEASE_%s", chromeBuildVersion)
		v, err := fetchVersion(u)
		if err == nil {
			return v, nil
		}
	}

	chromeMajorVersion, err := strconv.Atoi(majorVersion(pkgVersion))
	if err != nil {
		return "", fmt.Errorf("chrome major version: %v", err)
	}
	u := baseUrl + fmt.Sprintf("LATEST_RELEASE_%d", chromeMajorVersion)
	v, err := fetchVersion(u)
	if err == nil {
		return v, nil
	} else {
		previousChromeMajorVersion := chromeMajorVersion - 1
		u = baseUrl + fmt.Sprintf("LATEST_RELEASE_%d", previousChromeMajorVersion)
		v, err := fetchVersion(u)
		if err == nil {
			return v, nil
		} else {
			return "", errors.New("could not find compatible chromedriver")
		}
	}
}

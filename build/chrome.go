package build

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	chromeDriverBinary = "chromedriver"
)

type Chrome struct {
	Requirements
}

func (c *Chrome) Build() error {

	// Build dev image
	devDestDir, err := tmpDir()
	if err != nil {
		return fmt.Errorf("create dev temporary dir: %v", err)
	}

	srcDir := "chrome/apt"
	pkgSrcPath, pkgVersion, err := c.BrowserSource.Prepare()
	if err != nil {
		return fmt.Errorf("invalid browser source: %v", err)
	}

	if pkgSrcPath != "" {
		srcDir = "chrome/local"
		pkgDestPath := filepath.Join(devDestDir, "google-chrome.deb")
		err = os.Rename(pkgSrcPath, pkgDestPath)
		if err != nil {
			return fmt.Errorf("move package: %v", err)
		}
	}

	devImageTag := fmt.Sprintf("selenoid/dev_chrome:%s", pkgVersion)
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
		return fmt.Errorf("init dev image: %v", err)
	}
	image.BuildArgs = append(image.BuildArgs, fmt.Sprintf("VERSION=%s", pkgVersion))

	driverVersion, err := c.downloadChromedriver(image.Dir, pkgVersion)
	if err != nil {
		return fmt.Errorf("failed to download Chromedriver: %v", err)
	}
	image.Labels = []string{fmt.Sprintf("driver=chromedriver:%s", driverVersion)}

	err = image.Build()
	if err != nil {
		return fmt.Errorf("build image: %v", err)
	}

	err = image.Test(c.TestsDir, "chrome", pkgVersion)
	if err != nil {
		return fmt.Errorf("test image: %v", err)
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

func (c *Chrome) downloadChromedriver(dir string, pkgVersion string) (string, error) {
	version := c.DriverVersion
	if version == LatestVersion {
		const baseUrl = "https://chromedriver.storage.googleapis.com"
		v, err := c.getLatestChromeDriver(baseUrl, pkgVersion)
		if err != nil {
			return "", fmt.Errorf("latest Chromedriver version: %v", err)
		}
		version = v
	}

	u := fmt.Sprintf("http://chromedriver.storage.googleapis.com/%s/chromedriver_linux64.zip", version)
	_, err := downloadDriver(u, chromeDriverBinary, dir)
	if err != nil {
		return "", fmt.Errorf("download Chromedriver: %v", err)
	}
	return version, nil
}

func (c *Chrome) getLatestChromeDriver(baseUrl string, pkgVersion string) (string, error) {
	fetchVersion := func(url string) (string, error) {
		resp, err := http.Get(url)
		if err != nil {
			return "", fmt.Errorf("request error: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("unsuccessful response: %d %s", resp.StatusCode, resp.Status)
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("read driver version: %v", err)
		}
		return string(data), nil
	}

	switch c.BrowserChannel {
	case "dev":
		chromeMajorVersion, err := strconv.Atoi(strings.Split(pkgVersion, ".")[0])
		if err != nil {
			return "", fmt.Errorf("chrome major version: %v", err)
		}
		for i := chromeMajorVersion; i > 0; i-- {
			u := path.Join(baseUrl, fmt.Sprintf("LATEST_RELEASE_%d", chromeMajorVersion))
			v, err := fetchVersion(u)
			if err == nil {
				return v, nil
			}
		}
		u := path.Join(baseUrl, "LATEST_RELEASE")
		return fetchVersion(u)
	default:
		chromeBuildVersion := pkgVersion
		pieces := strings.Split(pkgVersion, ".")
		if len(pieces) >= 3 {
			chromeBuildVersion = strings.Join(pieces[0:3], ".")
		}
		u := path.Join(baseUrl, fmt.Sprintf("LATEST_RELEASE_%s", chromeBuildVersion))
		return fetchVersion(u)
	}
}

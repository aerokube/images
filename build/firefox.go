package build

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	geckoDriverBinary = "geckodriver"
)

type Firefox struct {
	SelenoidVersion string
	SeleniumVersion string
	Requirements
}

func (c *Firefox) Build() error {

	if c.SelenoidVersion == "" && c.SeleniumVersion == "" {
		return errors.New("missing Selenoid or Selenium JAR version")
	}

	// Build dev image
	devDestDir, err := tmpDir()
	if err != nil {
		return fmt.Errorf("create dev temporary dir: %v", err)
	}

	devSrcDir := "firefox/apt"
	pkgSrcPath, pkgVersion, err := c.BrowserSource.Prepare()
	if err != nil {
		return fmt.Errorf("invalid browser source: %v", err)
	}

	if pkgSrcPath != "" {
		devSrcDir = "firefox/local"
		pkgDestPath := filepath.Join(devDestDir, "firefox.deb")
		err = os.Rename(pkgSrcPath, pkgDestPath)
		if err != nil {
			return fmt.Errorf("move package: %v", err)
		}
	}

	pkgTagVersion := extractVersion(pkgVersion)
	devImageTag := fmt.Sprintf("selenoid/dev_firefox:%s", pkgTagVersion)
	devImageRequirements := Requirements{NoCache: c.NoCache, Tags: []string{devImageTag}}
	devImage, err := NewImage(devSrcDir, devDestDir, devImageRequirements)
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

	firefoxMajorVersion, err := strconv.Atoi(majorVersion(pkgTagVersion))
	geckoDriverCompatible := firefoxMajorVersion > 48
	srcDir := "firefox/selenoid"
	if !geckoDriverCompatible {
		srcDir = "firefox/selenium"
	}

	image, err := NewImage(srcDir, destDir, c.Requirements)
	if err != nil {
		return fmt.Errorf("init dev image: %v", err)
	}
	image.BuildArgs = append(image.BuildArgs, fmt.Sprintf("VERSION=%s", pkgTagVersion))

	if geckoDriverCompatible {
		driverVersion, err := c.downloadGeckoDriver(image.Dir)
		if err != nil {
			return fmt.Errorf("failed to download geckodriver: %v", err)
		}
		labels := []string{fmt.Sprintf("driver=geckodriver:%s", driverVersion)}

		selenoidVersion, err := c.downloadSelenoid(image.Dir)
		if err != nil {
			return fmt.Errorf("failed to download Selenoid: %v", err)
		}
		labels = append(labels, fmt.Sprintf("selenoid=%s", selenoidVersion))
		image.Labels = labels

		firefoxMajorMinorVersion := majorMinorVersion(pkgTagVersion)
		browsersJsonFile := filepath.Join(image.Dir, "browsers.json")
		data, err := ioutil.ReadFile(browsersJsonFile)
		if err != nil {
			return fmt.Errorf("failed to read browsers.json: %v", err)
		}
		newContents := strings.Replace(string(data), "@@VERSION@@", firefoxMajorMinorVersion, -1)
		err = ioutil.WriteFile(browsersJsonFile, []byte(newContents), 0)
		if err != nil {
			return fmt.Errorf("failed to update browsers.json: %v", err)
		}
	} else {
		driverVersion, err := c.downloadSeleniumJAR(image.Dir)
		if err != nil {
			return fmt.Errorf("failed to download Selenium JAR: %v", err)
		}
		image.Labels = []string{fmt.Sprintf("driver=selenium:%s", driverVersion)}
	}

	err = image.Build()
	if err != nil {
		return fmt.Errorf("build image: %v", err)
	}

	err = image.Test(c.TestsDir, "firefox", pkgTagVersion)
	if err != nil {
		return fmt.Errorf("test image: %v", err)
	}

	err = image.Push()
	if err != nil {
		return fmt.Errorf("push image: %v", err)
	}

	return nil
}

func (c *Firefox) channelToBuildArgs() []string {
	switch c.BrowserChannel {
	case "beta":
		return []string{"PPA=ppa:mozillateam/firefox-next"}
	case "dev":
		return []string{"PACKAGE=firefox-trunk", "PPA=ppa:ubuntu-mozilla-daily/ppa"}
	case "esr":
		return []string{"PACKAGE=firefox-esr", "PPA=ppa:mozillateam/ppa"}
	default:
		return []string{}
	}
}

func (c *Firefox) downloadGeckoDriver(dir string) (string, error) {
	version := c.DriverVersion
	if version == LatestVersion {
		v, err := latestGithubRelease("mozilla/geckodriver")
		if err != nil {
			return "", fmt.Errorf("latest geckodriver version: %v", err)
		}
		version = v
	}

	u := fmt.Sprintf("https://github.com/mozilla/geckodriver/releases/download/v%s/geckodriver-v%s-linux64.tar.gz", version, version)
	_, err := downloadDriver(u, geckoDriverBinary, dir)
	if err != nil {
		return "", fmt.Errorf("download geckodriver: %v", err)
	}
	return version, nil
}

func (c *Firefox) downloadSelenoid(dir string) (string, error) {
	version := c.SelenoidVersion
	if version == LatestVersion {
		v, err := latestGithubRelease("aerokube/selenoid")
		if err != nil {
			return "", fmt.Errorf("latest Selenoid version: %v", err)
		}
		version = v
	}

	u := fmt.Sprintf("https://github.com/aerokube/selenoid/releases/download/%s/selenoid_linux_amd64", version)
	data, err := downloadFile(u)
	if err != nil {
		return "", fmt.Errorf("download Selenoid: %v", err)
	}
	outputPath := filepath.Join(dir, "selenoid")
	err = ioutil.WriteFile(outputPath, data, 0755)
	if err != nil {
		return "", fmt.Errorf("save Selenoid: %v", err)
	}
	return version, nil
}

func (c *Firefox) downloadSeleniumJAR(dir string) (string, error) {
	version := c.SeleniumVersion
	var u string
	switch version {
	case "2.15.0", "2.19.0", "2.20.0", "2.21.0", "2.25.0", "2.32.0", "2.35.0", "2.37.0", "2.39.0", "2.40.0", "2.41.0", "2.43.1", "2.44.0", "2.45.0", "2.48.2":
		u = fmt.Sprintf("https://repo.jenkins-ci.org/releases/org/seleniumhq/selenium/selenium-server-standalone/%s/selenium-server-standalone-$selenium_version.jar", version)
	case "2.47.1":
		u = "http://selenium-release.storage.googleapis.com/2.47/selenium-server-standalone-2.47.1.jar"
	case "2.53.1":
		u = "http://selenium-release.storage.googleapis.com/2.53/selenium-server-standalone-2.53.1.jar"
	case "3.2.0":
		u = "http://selenium-release.storage.googleapis.com/3.2/selenium-server-standalone-3.2.0.jar"
	case "3.3.1":
		u = "http://selenium-release.storage.googleapis.com/3.3/selenium-server-standalone-3.3.1.jar"
	case "3.4.0":
		u = "https://selenium-release.storage.googleapis.com/3.4/selenium-server-standalone-3.4.0.jar"
	default:
		return "", fmt.Errorf("unsupported Selenium JAR version: %s", version)
	}
	data, err := downloadFile(u)
	if err != nil {
		return "", fmt.Errorf("download Selenium JAR: %v", err)
	}
	outputPath := filepath.Join(dir, "selenium-server-standalone.jar")
	err = ioutil.WriteFile(outputPath, data, 0644)
	if err != nil {
		return "", fmt.Errorf("save Selenium JAR: %v", err)
	}
	return version, nil

}

package build

import (
	"encoding/json"
	"errors"
	"fmt"
	hv "github.com/hashicorp/go-version"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	chromeDriverBinary    = "chromedriver"
	newChromeDriverBinary = "chromedriver-linux64/chromedriver"
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

	chromeDriverVersions, err := fetchChromeDriverVersions()
	if err != nil {
		return fmt.Errorf("fetch chromedriver versions: %v", err)
	}

	driverVersion, err := c.parseChromeDriverVersion(pkgTagVersion, chromeDriverVersions)
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

	err = c.downloadChromeDriver(image.Dir, driverVersion, chromeDriverVersions)
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

func (c *Chrome) parseChromeDriverVersion(pkgVersion string, chromeDriverVersions map[string]string) (string, error) {
	version := c.DriverVersion
	if version == LatestVersion {

		var matchingVersions []string
		for mv := range chromeDriverVersions {
			if strings.Contains(mv, pkgVersion) {
				matchingVersions = append(matchingVersions, mv)
			}
		}
		if len(matchingVersions) > 0 {
			sort.SliceStable(matchingVersions, func(i, j int) bool {
				l := matchingVersions[i]
				r := matchingVersions[j]
				lv, err := hv.NewVersion(l)
				if err != nil {
					return false
				}
				rv, err := hv.NewVersion(r)
				if err != nil {
					return false
				}
				return lv.LessThan(rv)
			})
			return matchingVersions[0], nil
		}

		const baseUrl = "https://chromedriver.storage.googleapis.com/"
		v, err := c.getLatestChromeDriver(baseUrl, pkgVersion)
		if err != nil {
			return "", err
		}
		return v, nil
	}
	return version, nil
}

func (c *Chrome) downloadChromeDriver(dir string, version string, chromeDriverVersions map[string]string) error {
	u := fmt.Sprintf("https://chromedriver.storage.googleapis.com/%s/chromedriver_linux64.zip", version)
	fn := chromeDriverBinary
	if cdu, ok := chromeDriverVersions[version]; ok {
		u = cdu
		fn = newChromeDriverBinary
	}
	outputPath, err := downloadDriver(u, fn, dir)
	if err != nil {
		return fmt.Errorf("download chromedriver: %v", err)
	}
	if fn == newChromeDriverBinary {
		err = os.Rename(outputPath, filepath.Join(dir, chromeDriverBinary))
		if err != nil {
			return fmt.Errorf("rename chromedriver: %v", err)
		}
	}
	return nil
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

func fetchChromeDriverVersions() (map[string]string, error) {
	const versionsURL = "https://googlechromelabs.github.io/chrome-for-testing/known-good-versions-with-downloads.json"
	resp, err := http.Get(versionsURL)
	if err != nil {
		return nil, fmt.Errorf("fetch chrome versions: %v", err)
	}
	defer resp.Body.Close()
	var cv ChromeVersions
	err = json.NewDecoder(resp.Body).Decode(&cv)
	if err != nil {
		return nil, fmt.Errorf("decode json: %v", err)
	}
	ret := make(map[string]string)
	const platformLinux64 = "linux64"
	const chromeDriver = "chromedriver"
	for _, v := range cv.Versions {
		version := v.Version
		if cd, ok := v.Downloads[chromeDriver]; ok {
			for _, d := range cd {
				u := d.URL
				if u != "" && d.Platform == platformLinux64 {
					ret[version] = u
				}
			}
		}
	}
	return ret, nil
}

type ChromeVersions struct {
	Versions []ChromeVersion `json:"versions"`
}

type ChromeVersion struct {
	Version   string                      `json:"version"`
	Downloads map[string][]ChromeDownload `json:"downloads"`
}

type ChromeDownload struct {
	Platform string `json:"platform"`
	URL      string `json:"url"`
}

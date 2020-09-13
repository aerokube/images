package build

import (
	"context"
	"errors"
	"fmt"
	"github.com/markbates/pkger"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

const (
	LatestVersion = "latest"
)

type Requirements struct {
	BrowserSource  BrowserSource
	BrowserChannel string // "beta", "esr", "dev" and so on
	DriverVersion  string
	Tags           []string
	NoCache        bool
	TestsDir       string
	RunTests       bool
	IgnoreTests    bool
	PushImage      bool
}

type BrowserSource string

// Return regular file corresponding to this source and optionally download this file
func (bs *BrowserSource) Prepare() (string, string, error) {
	src := string(*bs)
	if src == "" {
		return "", "", errors.New("empty browser source")
	}
	if _, err := os.Stat(src); err == nil {
		pkgName := filepath.Base(src)
		return src, extractVersion(pkgName), nil
	} else if u, err := url.Parse(src); strings.HasPrefix(src, "http") && err == nil {
		pkgName := path.Base(src)
		data, err := downloadFile(u.String())
		if err != nil {
			return "", "", fmt.Errorf("download file: %v", err)
		}
		f, err := ioutil.TempFile("", "selenoid-images")
		if err != nil {
			return "", "", fmt.Errorf("temporary file: %v", err)
		}
		outputFileName := f.Name()
		err = ioutil.WriteFile(outputFileName, data, 0644)
		if err != nil {
			return "", "", fmt.Errorf("save downloaded file: %v", err)
		}
		return outputFileName, extractVersion(pkgName), nil
	}
	return "", src, nil
}

func extractVersion(name string) string {
	pieces := strings.Split(name, "_")
	version := name
	if len(pieces) >= 2 {
		version = pieces[1]
	}
	pieces = strings.Split(version, "+")
	pieces = strings.Split(pieces[0], "-")
	pieces = strings.Split(pieces[0], "~")
	return pieces[0]
}

func versionN(pkgVersion string, n int) string {
	buildVersion := pkgVersion
	pieces := strings.Split(pkgVersion, ".")
	if len(pieces) >= n {
		buildVersion = strings.Join(pieces[0:n], ".")
	}
	return buildVersion
}

func majorVersion(pkgVersion string) string {
	return versionN(pkgVersion, 1)
}

func majorMinorVersion(pkgVersion string) string {
	return versionN(pkgVersion, 2)
}

func buildVersion(pkgVersion string) string {
	return versionN(pkgVersion, 3)
}

type Image struct {
	Dir        string
	BuildArgs  []string
	Labels     []string
	FileServer bool
	Requirements
}

func NewImage(srcDir string, destDir string, req Requirements) (*Image, error) {
	if !requireCommand("docker") {
		return nil, fmt.Errorf("docker is not installed")
	}

	dir, err := copySourceFiles(srcDir, destDir)
	if err != nil {
		return nil, fmt.Errorf("copy source files: %v", err)
	}

	if len(req.Tags) == 0 {
		return nil, errors.New("image tag is required")
	}
	return &Image{Dir: dir, Requirements: req}, nil
}

func requireCommand(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func tmpDir() (string, error) {
	dir, err := ioutil.TempDir("", "selenoid-images")
	if err != nil {
		return "", fmt.Errorf("create temporary dir: %v", err)
	}
	return dir, nil
}

func copySourceFiles(srcDir string, destDir string) (string, error) {

	const prefix = "/static"
	walkDir := filepath.Join(prefix, srcDir)
	err := pkger.Walk(walkDir, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		regex := regexp.MustCompile(`^.+:/static(.+)$`)
		relativePath := regex.FindStringSubmatch(path)[1]
		outputPath := filepath.Join(destDir, relativePath)
		if info.IsDir() {
			return os.MkdirAll(outputPath, info.Mode())
		}

		fileDir := filepath.Join(destDir, filepath.Dir(relativePath))
		if !fileExists(fileDir) {
			log.Printf("mkdir dir %s", fileDir)
			return os.MkdirAll(fileDir, info.Mode())
		}

		src, err := pkger.Open(path)
		if err != nil {
			return err
		}
		defer src.Close()

		dest, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer dest.Close()

		_, err = io.Copy(dest, src)
		if err != nil {
			return err
		}

		err = dest.Sync()
		if err != nil {
			return err
		}

		err = os.Chmod(outputPath, info.Mode())
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return filepath.Join(destDir, srcDir), nil
}

func (i *Image) Build() error {

	args := []string{"build"}
	for _, tag := range i.Tags {
		args = append(args, "-t", tag)
	}
	if len(i.BuildArgs) > 0 {
		for _, arg := range i.BuildArgs {
			args = append(args, "--build-arg", arg)
		}
	}
	if httpProxy := os.Getenv("HTTP_PROXY"); httpProxy != "" {
		args = append(args, "--build-arg", fmt.Sprintf("http_proxy=%s", httpProxy))
	}
	if httpsProxy := os.Getenv("HTTPS_PROXY"); httpsProxy != "" {
		args = append(args, "--build-arg", fmt.Sprintf("https_proxy=%s", httpsProxy))
	}
	if len(i.Labels) > 0 {
		for _, label := range i.Labels {
			args = append(args, "--label", label)
		}
	}

	if i.NoCache {
		args = append(args, "--no-cache")
	}

	if i.FileServer {
		server := &http.Server{
			Handler: http.FileServer(http.Dir(i.Dir)),
		}
		ln, err := net.Listen("tcp", ":8080")
		if err != nil {
			return fmt.Errorf("failed to allocate free port: %v", err)
		}

		e := make(chan error)
		go func() {
			e <- server.Serve(ln)
		}()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		defer server.Shutdown(ctx)
		if runtime.GOOS == "linux" {
			ip, err := dockerHostIP()
			if err != nil {
				return fmt.Errorf("failed to detect host machine IP: %v", err)
			}
			args = append(args, "--add-host", fmt.Sprintf("host.docker.internal:%s", ip))
		}
	}

	args = append(args, i.Dir)
	log.Printf("running command: docker %s", strings.Join(args, " "))
	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start command: %v", err)
	}
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("command execution error: %v", err)
	}
	return nil
}

func dockerHostIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Name != "docker0" {
			continue
		}
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("no live network interfaces detected")
}

func (i *Image) Test(testsDir string, browserName string, browserVersion string) error {
	if !i.RunTests {
		log.Println("not running tests")
		return nil
	}
	ref := i.Tags[0]
	err := doTest(ref, testsDir, browserName, browserVersion)
	if err != nil {
		if i.IgnoreTests {
			log.Printf("ignoring tests: %v", err)
			return nil
		}
		return fmt.Errorf("tests error: %v", err)
	}
	log.Println("tests passed")
	return nil
}

func doTest(ref string, testsDir string, browserName string, browserVersion string) error {
	if !fileExists(testsDir) {
		return fmt.Errorf("tests directory %s does not exist", testsDir)
	}

	seleniumUrl := "http://localhost:4445/"
	if browserName == "firefox" || (browserName == "opera" && browserVersion == "12.16") {
		seleniumUrl = "http://localhost:4445/wd/hub"
	}

	exec.Command("docker", "rm", "-f", "selenium").Output()
	defer func() {
		exec.Command("docker", "rm", "-f", "selenium").Output()
	}()

	output, err := exec.Command("docker", "run", "-d", "--name", "selenium", "--privileged", "-p", "4445:4444", ref).Output()
	if err != nil {
		return fmt.Errorf("failed to start docker image: %s %v", string(output), err)
	}

	if !requireCommand("mvn") {
		return fmt.Errorf("maven is not installed")
	}

	mvnCmd := exec.Command("mvn", "clean", "test",
		fmt.Sprintf("-Dgrid.connection.url=%s", seleniumUrl),
		fmt.Sprintf("-Dgrid.browser.name=%s", browserName),
		fmt.Sprintf("-Dgrid.browser.version=%s", browserVersion),
	)
	mvnCmd.Dir = testsDir
	mvnCmd.Stdout = os.Stdout
	mvnCmd.Stderr = os.Stderr
	err = mvnCmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start tests: %v", err)
	}
	err = mvnCmd.Wait()
	if err != nil {
		return fmt.Errorf("tests finished with error: %v", err)
	}

	return nil
}

func fileExists(p string) bool {
	_, err := os.Stat(p)
	return !os.IsNotExist(err)
}

func (i *Image) Push() error {
	if i.PushImage {
		for _, tag := range i.Tags {
			cmd := exec.Command("docker", "push", tag)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Start()
			if err != nil {
				return fmt.Errorf("invalid docker push %s: %v", tag, err)
			}
			err = cmd.Wait()
			if err != nil {
				return fmt.Errorf("pushing failed: %v", err)
			}
		}
	}
	return nil
}

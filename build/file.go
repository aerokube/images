package build

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gopkg.in/cheggaaa/pb.v1"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	zipMagicHeader  = "504b"
	gzipMagicHeader = "1f8b"
)

func downloadFile(url string) ([]byte, error) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	err := downloadFileWithProgressBar(url, w)
	if err != nil {
		return nil, err
	}
	w.Flush()
	return b.Bytes(), nil
}

func downloadFileWithProgressBar(url string, w io.Writer) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("file download error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected response code: %d", resp.StatusCode)
	}

	contentLength := int(resp.ContentLength)
	writer := w

	if contentLength > 0 {
		bar := pb.New(contentLength).SetUnits(pb.U_BYTES)
		bar.Output = os.Stderr
		bar.Start()
		defer bar.Finish()
		writer = io.MultiWriter(w, bar)
	}

	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %v", err)
	}
	return nil
}

func downloadDriver(url string, filename string, outputDir string) (string, error) {
	log.Printf("downloading driver from %s", url)
	data, err := downloadFile(url)
	if err != nil {
		return "", fmt.Errorf("failed to download driver archive: %v", err)
	}
	return extractFile(data, filename, outputDir)
}

func getMagicHeader(data []byte) string {
	if len(data) >= 2 {
		return hex.EncodeToString(data[:2])
	}
	return ""
}

func isZipFile(data []byte) bool {
	return getMagicHeader(data) == zipMagicHeader
}

func isTarGzFile(data []byte) bool {
	return getMagicHeader(data) == gzipMagicHeader
}

func extractFile(data []byte, filename string, outputDir string) (string, error) {
	if isZipFile(data) {
		return unzip(data, filename, outputDir)
	} else if isTarGzFile(data) {
		return untar(data, filename, outputDir)
	} else {
		outputPath := filepath.Join(outputDir, filename)
		err := os.WriteFile(outputPath, data, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("failed to save file %s: %v", outputPath, err)
		}
		return outputPath, nil
	}
}

// Based on http://stackoverflow.com/questions/20357223/easy-way-to-unzip-file-with-golang
func unzip(data []byte, fileName string, outputDir string) (string, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) (string, error) {
		rc, err := f.Open()
		if err != nil {
			return "", err
		}
		defer rc.Close()

		outputPath := filepath.Join(outputDir, f.Name)

		if f.FileInfo().IsDir() {
			return "", fmt.Errorf("can only unzip files but %s is a directory", f.Name)
		}

		err = outputFile(outputPath, f.Mode(), rc)
		if err != nil {
			return "", err
		}
		return outputPath, nil
	}

	if err == nil {
		for _, f := range zr.File {
			if f.Name == fileName {
				return extractAndWriteFile(f)
			}
		}
		err = fmt.Errorf("file %s does not exist in archive", fileName)
	}

	return "", err
}

// Based on https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07
func untar(data []byte, fileName string, outputDir string) (string, error) {

	gzr, err := gzip.NewReader(bytes.NewReader(data))
	defer gzr.Close()

	extractAndWriteFile := func(tr *tar.Reader, header *tar.Header) (string, error) {

		outputPath := filepath.Join(outputDir, header.Name)

		if header.Typeflag == tar.TypeDir {
			return "", fmt.Errorf("can only untar files but %s is a directory", header.Name)
		}

		err = outputFile(outputPath, os.FileMode(header.Mode), tr)
		if err != nil {
			return "", err
		}
		return outputPath, nil
	}

	if err == nil {
		tr := tar.NewReader(gzr)

	loop:
		for {
			header, err := tr.Next()
			switch {
			case err == io.EOF:
				break loop
			case err != nil:
				return "", err
			case header == nil:
				continue
			}
			return extractAndWriteFile(tr, header)
		}
		err = fmt.Errorf("file %s does not exist in archive", fileName)
	}

	return "", err
}

func outputFile(outputPath string, mode os.FileMode, r io.Reader) error {
	os.MkdirAll(filepath.Dir(outputPath), 0755)
	f, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	if err != nil {
		return err
	}
	return nil
}

func doSendGet(url string, token string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid request: %v", err)
	}
	if token != "" {
		req.Header.Add("Authorization", fmt.Sprintf("token %s", token))
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unsuccessful response: %d %s", resp.StatusCode, resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %v", err)
	}
	return data, nil
}

func sendGet(url string) ([]byte, error) {
	return doSendGet(url, "")
}

func sendGetWithAuth(url string, token string) ([]byte, error) {
	return doSendGet(url, token)
}

func latestGithubRelease(repo string) (string, error) {
	token := os.Getenv("GITHUB_TOKEN")
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	data, err := sendGetWithAuth(url, token)
	if err != nil {
		return "", fmt.Errorf("get latest github release data: %v", err)
	}
	type info struct {
		TagName string `json:"tag_name"`
	}
	var i info
	err = json.Unmarshal(data, &i)
	if err != nil {
		return "", fmt.Errorf("json unmarshal: %v", err)
	}
	return i.TagName, nil
}

func githubLinuxAssetURL(repo string, version string) (string, error) {
	token := os.Getenv("GITHUB_TOKEN")
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases", repo)
	data, err := sendGetWithAuth(url, token)
	if err != nil {
		return "", fmt.Errorf("get github releases data: %v", err)
	}
	type AssetInfo struct {
		BrowserDownloadURL string `json:"browser_download_url"`
	}
	type Release struct {
		Assets []AssetInfo `json:"assets"`
	}
	type Releases []Release
	var releases Releases
	err = json.Unmarshal(data, &releases)
	if err != nil {
		return "", fmt.Errorf("json unmarshal: %v", err)
	}
	for _, release := range releases {
		for _, asset := range release.Assets {
			if version != LatestVersion && !strings.Contains(asset.BrowserDownloadURL, version) {
				continue
			}
			if strings.Contains(asset.BrowserDownloadURL, "linux") {
				return asset.BrowserDownloadURL, nil
			}
		}
	}
	return "", fmt.Errorf("could not find github linux asset")
}

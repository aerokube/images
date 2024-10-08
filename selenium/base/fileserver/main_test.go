package main

import (
	. "github.com/aandryashin/matchers"
	. "github.com/aandryashin/matchers/httpresp"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

var (
	dir string
	srv *httptest.Server
)

func init() {
	dir, _ = os.MkdirTemp("", "fileserver")
	srv = httptest.NewServer(mux(dir))
}

func withUrl(path string) string {
	return srv.URL + path
}

func TestDownloadAndRemoveFile(t *testing.T) {
	tempFile, _ := os.CreateTemp(dir, "fileserver")
	_ = os.WriteFile(tempFile.Name(), []byte("test-data"), 0644)
	tempFileName := filepath.Base(tempFile.Name())
	tempFileStat, _ := tempFile.Stat()
	resp, err := http.Get(withUrl("/" + tempFileName))
	AssertThat(t, err, Is{nil})
	AssertThat(t, resp, Code{200})
	_, err = os.Stat(tempFile.Name())
	AssertThat(t, err, Is{nil})

	var files []FileInfo

	rsp, err := http.Get(withUrl("/?json"))
	AssertThat(t, err, Is{nil})
	AssertThat(t, rsp, Code{http.StatusOK})
	AssertThat(t, rsp, IsJson{&files})
	AssertThat(t, files, EqualTo{[]FileInfo{{
		Name:         tempFileName,
		Size:         tempFileStat.Size(),
		LastModified: tempFileStat.ModTime().Unix(),
	}}})

	hash, _ := getHash(tempFile.Name(), "md5")
	rsp, err = http.Get(withUrl("/?json&hash=md5"))
	AssertThat(t, err, Is{nil})
	AssertThat(t, rsp, Code{http.StatusOK})
	AssertThat(t, rsp, IsJson{&files})
	AssertThat(t, files, EqualTo{[]FileInfo{{
		Name:         tempFileName,
		Size:         tempFileStat.Size(),
		LastModified: tempFileStat.ModTime().Unix(),
		HashSum:      hash,
	}}})

	req, _ := http.NewRequest(http.MethodDelete, withUrl("/"+tempFileName), nil)
	resp, err = http.DefaultClient.Do(req)
	AssertThat(t, err, Is{nil})
	AssertThat(t, resp, Code{200})
	_, err = os.Stat(tempFile.Name())
	AssertThat(t, err, Not{nil})
}

func TestRemoveMissingFile(t *testing.T) {
	req, _ := http.NewRequest(http.MethodDelete, withUrl("/missing-file"), nil)
	resp, err := http.DefaultClient.Do(req)
	AssertThat(t, err, Is{nil})
	AssertThat(t, resp, Code{404})
}

package main

import (
	. "github.com/aandryashin/matchers"
	. "github.com/aandryashin/matchers/httpresp"
	"io/ioutil"
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
	dir, _ = ioutil.TempDir("", "fileserver")
	srv = httptest.NewServer(mux(dir))
}

func withUrl(path string) string {
	return srv.URL + path
}

func TestDownloadAndRemoveFile(t *testing.T) {
	tempFile, _ := ioutil.TempFile(dir, "fileserver")
	ioutil.WriteFile(tempFile.Name(), []byte("test-data"), 0644)
	tempFileName := filepath.Base(tempFile.Name())
	resp, err := http.Get(withUrl("/" + tempFileName))
	AssertThat(t, err, Is{nil})
	AssertThat(t, resp, Code{200})
	_, err = os.Stat(tempFile.Name())
	AssertThat(t, err, Is{nil})
	
	req, _ := http.NewRequest(http.MethodDelete, withUrl("/" + tempFileName), nil)
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

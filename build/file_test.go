package build

import (
	. "github.com/aandryashin/matchers"
	"io/ioutil"
	"os"
	"testing"
)

func TestUnzip(t *testing.T) {
	data := readFile(t, "testfile.zip")
	AssertThat(t, isZipFile(data), Is{true})
	AssertThat(t, isTarGzFile(data), Is{false})
	testUnpack(t, data, "zip-testfile", func(data []byte, filePath string, outputDir string) (string, error) {
		return unzip(data, filePath, outputDir)
	}, "zip\n")
}

func TestUntar(t *testing.T) {
	data := readFile(t, "testfile.tar.gz")
	AssertThat(t, isTarGzFile(data), Is{true})
	AssertThat(t, isZipFile(data), Is{false})
	testUnpack(t, data, "gzip-testfile", func(data []byte, filePath string, outputDir string) (string, error) {
		return untar(data, filePath, outputDir)
	}, "gzip\n")
}

func withTmpDir(t *testing.T, prefix string, fn func(*testing.T, string)) {
	dir, err := ioutil.TempDir("", prefix)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	fn(t, dir)

}

func testUnpack(t *testing.T, data []byte, fileName string, fn func([]byte, string, string) (string, error), correctContents string) {

	withTmpDir(t, "test-unpack", func(t *testing.T, dir string) {
		unpackedFile, err := fn(data, fileName, dir)
		if err != nil {
			t.Fatal(err)
		}

		if !fileExists(unpackedFile) {
			t.Fatalf("file %s does not exist\n", unpackedFile)
		}

		unpackedFileContents := string(readFile(t, unpackedFile))
		if unpackedFileContents != correctContents {
			t.Fatalf("incorrect unpacked file contents; expected: '%s', actual: '%s'\n", correctContents, unpackedFileContents)
		}
	})

}

func readFile(t *testing.T, fileName string) []byte {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Fatal(err)
	}
	return data
}


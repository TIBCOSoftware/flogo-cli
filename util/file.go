package fgutil

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"runtime"
	"path/filepath"
	"net/url"
	"net/http"
)

const FileURIPrefix = "file://"


type PathInfo struct {
	IsLocal  bool
	IsURL    bool
	FileURL  *url.URL
	FullPath string
	FileName string
}

func GetPathInfo(pathStr string) (*PathInfo, error) {

	fileURL, err := url.Parse(pathStr)

	pi := &PathInfo{}

	if err != nil {
		return nil, err
	}

	if len(fileURL.Scheme) > 0 {
		pi.FileURL = fileURL
		pi.IsURL = true

		filePath, local := URLToFilePath(fileURL)

		if local {
			pi.IsLocal = local
			pi.FullPath = filePath
		}
	} else {
		pi.FullPath = pathStr
	}

	idx := strings.LastIndex(pathStr, "/")
	pi.FileName = fileURL.Path[idx+1:]

	return pi, nil
}


//// ToFilePath convert fileURL to file path
//func ToFilePath(urlString string) (string, bool) {
//
//	itemURL, err := url.Parse(urlString)
//
//	if err != nil {
//		return
//	}
//
//	return URLToFilePath(itemURL)
//}

// ToFilePath convert fileURL to file path
func URLToFilePath(fileURL *url.URL) (string, bool) {

	if fileURL.Scheme == "file" {

		filePath :=fileURL.Path

		if runtime.GOOS == "windows" {
			if strings.HasPrefix(filePath, "/") {
				filePath = filePath[1:]
			}
			filePath = filepath.FromSlash(filePath)
		}

		filePath = strings.Replace(filePath, "%20", " ", -1)

		return filePath, true
	}

	return "", false
}

func ToAbsOsPath(filePath string) (string, error) {

	if runtime.GOOS == "windows" {
		filePath = filepath.FromSlash(filePath)
	}

	return filepath.Abs(filePath)
}

func PathToFileURL(filePath string) (string, error) {

	fixedPath, err := ToAbsOsPath(filePath)

	if (err != nil) {
		return "", err
	}

	fixedPath = strings.Replace(fixedPath, `\`, "/", -1)

	if runtime.GOOS == "windows" {
		return "file:///" + fixedPath, nil
	} else {
		return "file:///" + fixedPath, nil
	}
}


// WriteJSONtoFile encodes the data to json and saves it to a file
func WriteJSONtoFile(filePath string, data interface{}) error {

	f, _ := os.Create(filePath)
	defer f.Close()

	j, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	_, err = f.Write(j)
	if err != nil {
		return err
	}

	return nil
}

// CopyFile copies the file from the source to the destination file
func CopyFile(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			os.Chmod(dest, sourceinfo.Mode())
		}
	}

	return
}

func CopyRemoteFile(sourceURL string, dest string) (err error) {

	resp, err := http.Get(sourceURL)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destfile.Close()

	io.Copy(destfile, resp.Body)

	return nil
}

// CopyDir copies the specified directory and its contents to the specified destination
func CopyDir(source string, dest string) (err error) {

	// get properties of source dir
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// create dest dir

	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)

	objects, err := directory.Readdir(-1)

	for _, obj := range objects {

		sourcefilepointer := source + "/" + obj.Name()

		destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			// create sub-directories - recursively
			err = CopyDir(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// perform copy
			err = CopyFile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
	return
}

func DeleteFilesWithPrefix(dir string, filePrefix string) int {

	deleted := 0
	delFunc := func(path string, f os.FileInfo, err error) (e error) {

		if strings.HasPrefix(f.Name(), filePrefix) {
			os.Remove(path)
			deleted++
		}
		return
	}

	filepath.Walk(dir, delFunc)

	return deleted
}

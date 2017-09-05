package app

import (
	"os"

	"fmt"
	"github.com/TIBCOSoftware/flogo-cli/util"
	"io"
	"io/ioutil"
	"path"
	"path/filepath"
)

const (
	fileDescriptor    string = "flogo.json"
	fileMainGo        string = "main.go"
	fileImportsGo     string = "imports.go"
	fileEmbeddedAppGo string = "embeddedapp.go"
	lambdaHandlerGo   string = "handler.go"
	lambdaHandlerZip  string = "handler.zip"
	lambdaMake        string = "Makefile"

	pathFlogoLib string = "github.com/TIBCOSoftware/flogo-lib"
)

func createMainGoFile(codeSourcePath string, flogoJSON string) {

	data := struct {
		FlogoJSON string
	}{
		flogoJSON,
	}

	f, _ := os.Create(path.Join(codeSourcePath, fileMainGo))
	fgutil.RenderTemplate(f, tplNewMainGoFile, &data)
	f.Close()
}

var tplNewMainGoFile = `package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/TIBCOSoftware/flogo-lib/app"
	"github.com/TIBCOSoftware/flogo-lib/engine"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var (
	cp app.ConfigProvider
)

func main() {

	e := startEngine()
	code := waitSignal()
	stopEngine(e)
	exit(code)
}

func startEngine() engine.Engine{
	if cp == nil {
		// Use default config provider
		cp = app.DefaultConfigProvider()
	}

	app, err := cp.GetApp()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	e, err := engine.New(app)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	e.Start()

	return e
}

func stopEngine(e engine.IEngine){
	e.Stop()
}

func exit(code int){
	os.Exit(code)
}

func waitSignal() int{
	exitChan := setupSignalHandling()
	code := <-exitChan
	return code
}

func setupSignalHandling() chan int {

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	exitChan := make(chan int)
	go func() {
		for {
			s := <-signalChan
			switch s {
			// kill -SIGHUP
			case syscall.SIGHUP:
				exitChan <- 0
			// kill -SIGINT/Ctrl+c
			case syscall.SIGINT:
				exitChan <- 0
			// kill -SIGTERM
			case syscall.SIGTERM:
				exitChan <- 0
			// kill -SIGQUIT
			case syscall.SIGQUIT:
				exitChan <- 0
			default:
				logger.Debug("Unknown signal.")
				exitChan <- 1
			}
		}
	}()

	return exitChan
}
`

func createImportsGoFile(codeSourcePath string, deps []*Dependency) error {
	f, err := os.Create(path.Join(codeSourcePath, fileImportsGo))

	if err != nil {
		return err
	}

	fgutil.RenderTemplate(f, tplNewImportsGoFile, deps)
	f.Close()

	return nil
}

var tplNewImportsGoFile = `package main

import (

{{range $i, $dep := .}}	_ "{{ $dep.Ref }}"
{{end}}
)
`

func createEmbeddedAppGoFile(codeSourcePath string, flogoJSON string) {

	data := struct {
		FlogoJSON string
	}{
		flogoJSON,
	}

	f, _ := os.Create(path.Join(codeSourcePath, fileEmbeddedAppGo))
	fgutil.RenderTemplate(f, tplEmbeddedAppGoFile, &data)
	f.Close()
}

func removeEmbeddedAppGoFile(codeSourcePath string) {
	os.Remove(path.Join(codeSourcePath, fileEmbeddedAppGo))
}

var tplEmbeddedAppGoFile = `// Do not change this file, it has been generated using flogo-cli
// If you change it and rebuild the application your changes might get lost
package main

import (
	"encoding/json"

	"github.com/TIBCOSoftware/flogo-lib/app"
)

// embedded flogo app descriptor file
const flogoJSON string = ` + "`{{.FlogoJSON}}`" + `

func init () {
	cp = EmbeddedProvider()
}

// embeddedConfigProvider implementation of ConfigProvider
type embeddedProvider struct {
}

//EmbeddedProvider returns an app config from a compiled json file
func EmbeddedProvider() (app.ConfigProvider){
	return &embeddedProvider{}
}

// GetApp returns the app configuration
func (d *embeddedProvider) GetApp() (*app.Config, error){

	app := &app.Config{}
	err := json.Unmarshal([]byte(flogoJSON), app)
	if err != nil {
		return nil, err
	}
	return app, nil
}
`

func createLambdaHandlerGoFile(codeSourcePath string) error {

	f, err := os.Create(path.Join(codeSourcePath, lambdaHandlerGo))
	if err != nil {
		return err
	}
	fgutil.RenderTemplate(f, lambdaHandlerGoFile, nil)
	return f.Close()
}

func removeLambdaHandlerGoFile(codeSourcePath string) error {
	return os.Remove(path.Join(codeSourcePath, lambdaHandlerGo))
}

var lambdaHandlerGoFile = `// Do not change this file, it has been generated using flogo-cli
// If you change it and rebuild the application your changes might get lost
package main

import (
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
)

func Handle(evt interface{}, ctx *runtime.Context) (string, error) {
	e := startEngine()
	stopEngine(e)
	return "Flow finished succesfully", nil
}
`

func createLambdaMakeFile(codeSourcePath string) error {

	f, err := os.Create(path.Join(codeSourcePath, lambdaMake))
	if err != nil {
		return err
	}
	fgutil.RenderTemplate(f, lambdaMakeFile, nil)
	return f.Close()
}

func removeLambdaMakeFile(codeSourcePath string) error {
	return os.Remove(path.Join(codeSourcePath, lambdaMake))
}

var lambdaMakeFile = `#
# This is free and unencumbered software released into the public domain.
#
# Anyone is free to copy, modify, publish, use, compile, sell, or
# distribute this software, either in source code form or as a compiled
# binary, for any purpose, commercial or non-commercial, and by any
# means.
#
# In jurisdictions that recognize copyright laws, the author or authors
# of this software dedicate any and all copyright interest in the
# software to the public domain. We make this dedication for the benefit
# of the public at large and to the detriment of our heirs and
# successors. We intend this dedication to be an overt act of
# relinquishment in perpetuity of all present and future rights to this
# software under copyright law.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
# EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
# MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
# IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
# OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
# ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
# OTHER DEALINGS IN THE SOFTWARE.
#
# For more information, please refer to <http://unlicense.org/>
#

HANDLER ?= handler
PACKAGE ?= $(HANDLER)

ifeq ($(OS),Windows_NT)
	GOPATH ?= $(USERPROFILE)/go
	GOPATH := /$(subst ;,:/,$(subst \,/,$(subst :,,$(GOPATH))))
	CURDIR := /$(subst :,,$(CURDIR))
	RM := del /q
else
	GOPATH ?= $(HOME)/go
	RM := rm -f
endif

MAKEFILE = $(word $(words $(MAKEFILE_LIST)),$(MAKEFILE_LIST))

docker:
	docker run --rm\
		-e HANDLER=$(HANDLER)\
		-e PACKAGE=$(PACKAGE)\
		-e GOPATH=$(GOPATH)\
		-e LDFLAGS='$(LDFLAGS)'\
		-v $(CURDIR):$(CURDIR)\
		$(foreach GP,$(subst :, ,$(GOPATH)),-v $(GP):$(GP))\
		-w $(CURDIR)\
		eawsy/aws-lambda-go-shim:latest make -f $(MAKEFILE) all

.PHONY: docker

all: build pack perm

.PHONY: all

build:
	go build -buildmode=plugin -ldflags='-w -s $(LDFLAGS)' -o $(HANDLER).so

.PHONY: build

pack:
	pack $(HANDLER) $(HANDLER).so $(PACKAGE).zip

.PHONY: pack

perm:
	chown $(shell stat -c '%u:%g' .) $(HANDLER).so $(PACKAGE).zip
	$(RM) $(HANDLER).so

.PHONY: perm

clean:
	$(RM) $(HANDLER).so $(PACKAGE).zip

.PHONY: clean
`

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

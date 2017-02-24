package env

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/TIBCOSoftware/flogo-cli/util"
)

type GbProject struct {
	BinDir         string
	RootDir        string
	SourceDir      string
	VendorDir      string
	CodeSourcePath string
}

func NewGbProjectEnv() Project {

	env := &GbProject{}
	env.SourceDir = "src"
	env.VendorDir = fgutil.Path("vendor", "src")

	return env
}

func (e *GbProject) Init(path string) error {

	exists := fgutil.ExecutableExists("gb")

	if !exists {
		return errors.New("gb not installed")
	}

	e.RootDir = path
	e.SourceDir = fgutil.Path(path, "src")
	e.VendorDir = fgutil.Path(path, "vendor", "src")

	return nil
}

// Init creates directories for the gb project
func (e *GbProject) Create(createBin bool) error {

	if _, err := os.Stat(e.RootDir); err == nil {
		return fmt.Errorf("Cannot create project, directory '%s' already exists", e.RootDir)
	}

	os.MkdirAll(e.RootDir, os.ModePerm)
	os.MkdirAll(e.SourceDir, os.ModePerm)
	os.MkdirAll(e.VendorDir, os.ModePerm)

	if createBin {
		e.BinDir = fgutil.Path(e.RootDir, "bin")
		os.MkdirAll(e.BinDir, os.ModePerm)
	}

	return nil
}

// Open the project directory and validate its structure
func (e *GbProject) Open() error {

	info, err := os.Stat(e.RootDir)

	if err != nil || !info.IsDir() {
		return fmt.Errorf("Cannot open project, directory '%s' doesn't exists", e.RootDir)
	}

	info, err = os.Stat(e.SourceDir)

	if err != nil || !info.IsDir() {
		return errors.New("Invalid project, source directory doesn't exists")
	}

	info, err = os.Stat(e.VendorDir)

	if err != nil || !info.IsDir() {
		return errors.New("Invalid project, vendor directory doesn't exists")
	}

	binDir := fgutil.Path(e.RootDir, "bin")
	info, err = os.Stat(binDir)

	if err != nil || info.IsDir() {
		e.BinDir = binDir
	}

	return nil
}

func (e *GbProject) GetBinDir() string {
	return e.BinDir
}

func (e *GbProject) GetRootDir() string {
	return e.RootDir
}

func (e *GbProject) GetSourceDir() string {
	return e.SourceDir
}

func (e *GbProject) GetVendorDir() string {
	return e.VendorDir
}

func (e *GbProject) InstallDependency(path string, version string) error {
	var cmd *exec.Cmd

	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)

	if version == "" {
		cmd = exec.Command("gb", "vendor", "fetch", path)
	} else {
		var tag string

		if version[0] != 'v' {
			tag = "v" + version
		} else {
			tag = version
		}

		cmd = exec.Command("gb", "vendor", "fetch", "-tag", tag, path)
	}

	os.Chdir(e.RootDir)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (e *GbProject) Build() error {
	cmd := exec.Command("gb", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)

	os.Chdir(e.RootDir)
	fmt.Println(e.RootDir)

	return cmd.Run()
}

func IsGbProject(path string) bool {

	sourceDir := fgutil.Path(path, "src")
	vendorDir := fgutil.Path(path, "vendor", "src")

	info, err := os.Stat(sourceDir)

	if err != nil || !info.IsDir() {
		return false
	}

	info, err = os.Stat(vendorDir)

	if err != nil || !info.IsDir() {
		return false
	}

	return true
}

//Env checker?
//IsProject(fgutil.Path string) bool

// Gb structure that contains gb project paths
type Gb struct {
	BinPath        string
	SourcePath     string
	VendorPath     string
	CodeSourcePath string
}

// NewGb creates a new Gb struct
func NewGb(codePath string) *Gb {

	env := &Gb{}
	env.BinPath = "bin"
	env.SourcePath = "src"
	env.VendorPath = fgutil.Path("vendor", "src")
	env.CodeSourcePath = fgutil.Path("src", codePath)

	return env
}

// Init creates directories for the gb project
func (e *Gb) Init(createBin bool) {
	os.MkdirAll(e.SourcePath, 0777)
	os.MkdirAll(e.VendorPath, 0777)
	os.MkdirAll(e.CodeSourcePath, 0777)

	if createBin {
		os.MkdirAll(e.BinPath, 0777)
	}
}

// Installed indicates if gb is installed
func (e *Gb) Installed() bool {
	return fgutil.ExecutableExists("gb")
}

// NewBinFilefgutil.Path creates a new file fgutil.Path in the bin directory
func (e *Gb) NewBinFilePath(fileName string) string {
	return fgutil.Path(e.BinPath, fileName)
}

// VendorFetch performs a 'gb vendor fetch'
func (e *Gb) VendorFetch(path string, version string) error {

	var cmd *exec.Cmd

	if version == "" {
		cmd = exec.Command("gb", "vendor", "fetch", path)
	} else {

		var tag string

		if version[0] != 'v' {
			tag = "v" + version
		} else {
			tag = version
		}

		cmd = exec.Command("gb", "vendor", "fetch", "-tag", tag, path)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// VendorDeleteSilent performs a 'gb vendor delete' silently
func (e *Gb) VendorDeleteSilent(path string) error {
	cmd := exec.Command("gb", "vendor", "delete", path)

	return cmd.Run()
}

// VendorDelete performs a 'gb vendor delete'
func (e *Gb) VendorDelete(path string) error {
	cmd := exec.Command("gb", "vendor", "delete", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Build performs a 'gb build'
func (e *Gb) Build() error {
	cmd := exec.Command("gb", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

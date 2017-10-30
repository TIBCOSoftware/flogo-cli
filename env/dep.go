package env

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/TIBCOSoftware/flogo-cli/config"
	"github.com/TIBCOSoftware/flogo-cli/util"
	"io/ioutil"
	"path"
)

type DepProject struct {
	BinDir             string
	RootDir            string
	SourceDir          string
	VendorDir          string
	VendorSrcDir       string
	CodeSourcePath     string
	AppDir             string
	FileDescriptorPath string
}

func NewDepProject() Project {
	return &DepProject{}
}

func (e *DepProject) Init(rootDir string) error {

	exists := fgutil.ExecutableExists("dep")

	if !exists {
		return errors.New("dep not installed")
	}
	e.RootDir = rootDir
	e.SourceDir = path.Join(e.RootDir, "src")
	e.BinDir = path.Join(e.RootDir, "bin")
	return nil
}

// Create creates directories for the project
func (e *DepProject) Create(createBin bool, vendorDir string) error {

	if _, err := os.Stat(e.RootDir); err == nil {
		return fmt.Errorf("Cannot create project, directory '%s' already exists", e.RootDir)
	}

	os.MkdirAll(e.RootDir, os.ModePerm)
	os.MkdirAll(e.SourceDir, os.ModePerm)

	return nil
}

// Open the project directory and validate its structure
func (e *DepProject) Open() error {

	// Check root dir
	info, err := os.Stat(e.RootDir)

	if err != nil || !info.IsDir() {
		return fmt.Errorf("Cannot open project, directory '%s' doesn't exists", e.RootDir)
	}

	// Check source dir
	info, err = os.Stat(e.SourceDir)

	if err != nil || !info.IsDir() {
		return errors.New("Invalid project, source directory doesn't exists")
	}

	// Check file descriptor
	fd := path.Join(e.RootDir, config.FileDescriptor)
	_, err = os.Stat(fd)

	if err != nil {
		return fmt.Errorf("Invalid project, file descriptor '%s' doesn't exists", fd)
	}
	e.FileDescriptorPath = fd

	fdbytes, err := ioutil.ReadFile(fd)
	if err != nil {
		return fmt.Errorf("Invalid reading file descriptor '%s' error: %s", fd, err)
	}

	descriptor, err := ParseAppDescriptor(string(fdbytes))
	if err != nil {
		return fmt.Errorf("Invalid parsing file descriptor '%s' error: %s", fd, err)
	}

	appName := descriptor.Name

	// Validate that there is an app dir
	e.AppDir = path.Join(e.SourceDir, appName)
	info, err = os.Stat(e.AppDir)

	if err != nil || !info.IsDir() {
		return fmt.Errorf("Invalid project, app directory '%s' doesn't exists", e.AppDir)
	}

	e.VendorDir = path.Join(e.AppDir, "vendor")
	e.VendorSrcDir = e.VendorDir

	return nil
}

func (e *DepProject) GetBinDir() string {
	return e.BinDir
}

func (e *DepProject) GetRootDir() string {
	return e.RootDir
}

func (e *DepProject) GetSourceDir() string {
	return e.SourceDir
}

func (e *DepProject) GetVendorDir() string {
	return e.VendorDir
}

func (e *DepProject) GetVendorSrcDir() string {
	return e.VendorSrcDir
}

// GetAppDir returns the directory of the app
func (e *DepProject) GetAppDir() string {
	return e.AppDir
}

func (e *DepProject) InstallDependency(depPath string, version string) error {
	// Deprecated, dependency managements responsibility
	return nil
}

func (e *DepProject) UninstallDependency(depPath string) error {
	// Deprecated, dependency managements responsibility
	return nil
}

func (e *DepProject) Build() error {
	exists := fgutil.ExecutableExists("go")
	if !exists {
		return errors.New("go not installed")
	}

	cmd := exec.Command("go", "install", "./...")
	cmd.Dir = e.GetAppDir()
	newEnv := os.Environ()
	newEnv = append(newEnv, fmt.Sprintf("GOPATH=%s", e.GetRootDir()))
	cmd.Env = newEnv

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func IsDepProject(projectPath string) bool {

	sourceDir := path.Join(projectPath, "src")
	vendorDir := path.Join(projectPath, "vendor", "src")

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
//IsProject(path.Join string) bool

// Gb structure that contains gb project paths
type Dep struct {
	BinPath        string
	SourcePath     string
	VendorPath     string
	CodeSourcePath string
}

// NewGb creates a new Gb struct
func NewDep(codePath string) *Gb {

	env := &Gb{}
	env.BinPath = "bin"
	env.SourcePath = "src"
	env.VendorPath = path.Join("vendor", "src")
	env.CodeSourcePath = path.Join("src", codePath)

	return env
}

// Init creates directories for the gb project
func (e *Dep) Init(createBin bool) {
	os.MkdirAll(e.SourcePath, 0777)
	os.MkdirAll(e.VendorPath, 0777)
	os.MkdirAll(e.CodeSourcePath, 0777)

	if createBin {
		os.MkdirAll(e.BinPath, 0777)
	}
}

// Installed indicates if gb is installed
func (e *Dep) Installed() bool {
	return fgutil.ExecutableExists("gb")
}

// NewBinFilepath.Join creates a new file path.Join in the bin directory
func (e *Dep) NewBinFilePath(fileName string) string {
	return path.Join(e.BinPath, fileName)
}

// VendorFetch performs a 'gb vendor fetch'
func (e *Dep) VendorFetch(depPath string, version string) error {

	var cmd *exec.Cmd

	if version == "" {
		cmd = exec.Command("gb", "vendor", "fetch", depPath)
	} else {

		var tag string

		if version[0] != 'v' {
			tag = "v" + version
		} else {
			tag = version
		}

		cmd = exec.Command("gb", "vendor", "fetch", "-tag", tag, depPath)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// VendorDeleteSilent performs a 'gb vendor delete' silently
func (e *Dep) VendorDeleteSilent(depPath string) error {
	cmd := exec.Command("gb", "vendor", "delete", depPath)

	return cmd.Run()
}

// VendorDelete performs a 'gb vendor delete'
func (e *Dep) VendorDelete(depPath string) error {
	cmd := exec.Command("gb", "vendor", "delete", depPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Build performs a 'gb build'
func (e *Dep) Build() error {
	cmd := exec.Command("gb", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// ParseAppDescriptor parse the application descriptor
func ParseAppDescriptor(appJson string) (*config.FlogoAppDescriptor, error) {
	descriptor := &config.FlogoAppDescriptor{}

	err := json.Unmarshal([]byte(appJson), descriptor)

	if err != nil {
		return nil, err
	}

	return descriptor, nil
}

package env

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	//"strings"

	"github.com/TIBCOSoftware/flogo-cli/util"
	"path"
)

// TempEnv allows you to temporarily change the value of an env variable
type TempEnv struct {
	key   string
	newValue string
	oldValue string
	wasSet   bool
}


type DepProject struct {
	BinDir         string
	RootDir        string
	SourceDir      string
	VendorDir      string
	VendorSrcDir   string
	CodeSourcePath string
}

type DepManager struct {
	AppDir string
}

func NewTempEnv(key, newValue string) *TempEnv {
	return &TempEnv{key:key, newValue:newValue}
}

// change changes the environment keys to the new value
func (te *TempEnv) change() error{
	// Save values
	te.oldValue, te.wasSet = os.LookupEnv(te.key)
	// Change
	return os.Setenv(te.key, te.newValue)
}

// revert reverts any changes performed by change
func (te *TempEnv) revert() error {
	if !te.wasSet {
		os.Unsetenv(te.key)
		return nil
	}
	return os.Setenv(te.key, te.oldValue)
}

// Init initializes the dependency manager
func (b *DepManager) Init(rootDir, appDir string) error {
	exists := fgutil.ExecutableExists("dep")
	if !exists {
		return errors.New("dep not installed")
	}

	// Change GOPATH temporarily
	//tempEnv := NewTempEnv("GOPATH", rootDir)
	//tempEnv.change()
	//defer tempEnv.revert()

	cmd := exec.Command("dep", "init", appDir)
	newEnv := os.Environ()
	newEnv = append(newEnv, fmt.Sprintf("GOPATH=%s", rootDir))
	cmd.Env = newEnv


	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return err
	}

	// TODO remove this prune once it gets absorved into dep ensure https://github.com/golang/dep/issues/944
	cmd = exec.Command("dep", "prune")
	cmd.Dir = appDir
	cmd.Env = newEnv

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func NewDepProject() Project {
	return &DepProject{}
}

func (e *DepProject) Init(basePath string) error {

	exists := fgutil.ExecutableExists("dep")

	if !exists {
		return errors.New("dep not installed")
	}
	e.RootDir = basePath
	e.SourceDir = path.Join(basePath, "src")
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

	info, err := os.Stat(e.RootDir)

	if err != nil || !info.IsDir() {
		return fmt.Errorf("Cannot open project, directory '%s' doesn't exists", e.RootDir)
	}

	info, err = os.Stat(e.SourceDir)

	if err != nil || !info.IsDir() {
		return errors.New("Invalid project, source directory doesn't exists")
	}

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

func (e *DepProject) InstallDependency(depPath string, version string) error {
	var cmd *exec.Cmd

	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)

	//check if dependency is installed
	if _, err := os.Stat(path.Join(e.VendorSrcDir, depPath)); err == nil {
		//todo ignore installed dependencies for now
		//exists, return
		return nil
	}

	if version == "" {
		//if strings.HasPrefix(depPath,"github.com/TIBCOSoftware/flogo-") {
		//	cmd = exec.Command("gb", "vendor", "fetch", "-branch", "entrypoint", depPath)
		//} else {
		cmd = exec.Command("gb", "vendor", "fetch", depPath)
		//}
	} else {
		var tag string

		if version[0] != 'v' {
			tag = "v" + version
		} else {
			tag = version
		}

		cmd = exec.Command("gb", "vendor", "fetch", "-tag", tag, depPath)
	}

	os.Chdir(e.RootDir)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (e *DepProject) UninstallDependency(depPath string) error {

	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)

	//check if dependency is installed
	if _, err := os.Stat(path.Join(e.VendorSrcDir, depPath)); err != nil {
		//todo ignore dependencies that are not installed for now
		//exists, return
		return nil
	}

	os.Chdir(e.RootDir)

	cmd := exec.Command("gb", "vendor", "delete", depPath)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (e *DepProject) Build() error {
	cmd := exec.Command("gb", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)

	os.Chdir(e.RootDir)
	fmt.Println(e.RootDir)

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

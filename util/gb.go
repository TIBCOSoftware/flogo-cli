package fgutil

import (
	"os"
	"os/exec"
)

//type GbProjectEnv struct {
//	BinDir         string
//	RootDir        string
//	SourceDir      string
//	VendorDir      string
//	CodeSourcePath string
//}
//
//func NewGbProjectEnv() ProjectEnv {
//
//	env := &GbProjectEnv{}
//	env.SourceDir = "src"
//	env.VendorDir = Path("vendor", "src")
//
//	return env
//}

//func (e *GbProjectEnv) Init(path string) error {
//
//	exists := ExecutableExists("gb")
//
//	if !exists {
//		return errors.New("gb not installed")
//	}
//
//	e.RootDir = path
//	e.SourceDir = Path(path,"src")
//	e.VendorDir = Path(path,"vendor", "src")
//
//	return nil
//}
//
//// Init creates directories for the gb project
//func (e *GbProjectEnv) Create(createBin bool) error {
//
//	if _, err := os.Stat(e.RootDir); err == nil {
//		return fmt.Errorf("Cannot create project, directory '%s' already exists", e.RootDir)
//	}
//
//	os.MkdirAll(e.RootDir, os.ModePerm)
//	os.MkdirAll(e.SourceDir, os.ModePerm)
//	os.MkdirAll(e.VendorDir, os.ModePerm)
//
//	if createBin {
//		e.BinDir = Path(e.RootDir,"bin")
//		os.MkdirAll(e.BinDir, os.ModePerm)
//	}
//
//	return nil
//}
//
//// Open the project directory and validate its structure
//func (e *GbProjectEnv) Open() error {
//
//	info, err := os.Stat(e.RootDir)
//
//	if err != nil || !info.IsDir() {
//		return fmt.Errorf("Cannot open project, directory '%s' doesn't exists", e.RootDir)
//	}
//
//	info, err = os.Stat(e.SourceDir)
//
//	if err != nil || !info.IsDir() {
//		return errors.New("Invalid project, source directory doesn't exists")
//	}
//
//	info, err = os.Stat(e.VendorDir)
//
//	if err != nil || !info.IsDir() {
//		return errors.New("Invalid project, vendor directory doesn't exists")
//	}
//
//	binDir := Path(e.RootDir,"bin")
//	info, err = os.Stat(binDir)
//
//	if err != nil || info.IsDir() {
//		e.BinDir = binDir
//	}
//
//	return nil
//}
//
//func (e *GbProjectEnv) GetBinDir() string {
//	return e.BinDir
//}
//
//func (e *GbProjectEnv) GetRootDir() string {
//	return e.RootDir
//}
//
//func (e *GbProjectEnv) GetSourceDir() string {
//	return e.SourceDir
//}
//
//func (e *GbProjectEnv) GetVendorDir() string {
//	return e.VendorDir
//}
//
//func (e *GbProjectEnv) InstallDependency(path string, version string) error {
//	var cmd *exec.Cmd
//
//	if version == "" {
//		cmd = exec.Command("gb", "vendor", "fetch", path)
//	} else {
//
//		var tag string
//
//		if version[0] != 'v' {
//			tag = "v" + version
//		} else {
//			tag = version
//		}
//
//		cmd = exec.Command("gb", "vendor", "fetch", "-tag", tag, path)
//	}
//
//	cmd.Stdout = os.Stdout
//	cmd.Stderr = os.Stderr
//
//	return cmd.Run()
//}
//
//func (e *GbProjectEnv) Build() error {
//	cmd := exec.Command("gb", "build")
//	cmd.Stdout = os.Stdout
//	cmd.Stderr = os.Stderr
//
//	return cmd.Run()
//}

func IsGbProject(path string) bool {

	sourceDir := Path(path,"src")
	vendorDir := Path(path,"vendor", "src")

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
//IsProject(Path string) bool

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
	env.VendorPath = Path("vendor", "src")
	env.CodeSourcePath = Path("src", codePath)

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
	return ExecutableExists("gb")
}

// NewBinFilePath creates a new file Path in the bin directory
func (e *Gb) NewBinFilePath(fileName string) string {
	return Path(e.BinPath, fileName)
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


package fgutil

import (
	"os"
	"os/exec"
	"strings"
)

// Gb stucture that contains gb project paths
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
	env.VendorPath = path("vendor", "src")
	env.CodeSourcePath = path("src", codePath)

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

// NewBinFilePath creates a new file path in the bin directory
func (e *Gb) NewBinFilePath(fileName string) string {
	return path(e.BinPath, fileName)
}

// VendorFetch performs a 'gb vendor fetch'
func (e *Gb) VendorFetch(path string) error {
	cmd := exec.Command("gb", "vendor", "fetch", path)
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

func path(parts ...string) string {
	return strings.Join(parts[:], string(os.PathSeparator))
}

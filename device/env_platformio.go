package device

import (
	"os/exec"
	"os"
	"strconv"
	"errors"
	"fmt"
	"path"
)

type PioProject struct {
	RootDir   string
	LibDir    string
	SourceDir string
}

func NewPlatformIoProject() Project {

	project := &PioProject{}

	return project
}

func (p *PioProject) GetRootDir() string {
	return p.RootDir
}

func (p *PioProject) GetSourceDir() string {
	return p.SourceDir
}

func (p *PioProject) GetLibDir() string {
	return p.LibDir
}

func (p *PioProject) Init(basePath string) error {

	_, err := exec.LookPath("platformio")

	if err != nil {
		return errors.New("platformio not installed")
	}

	p.RootDir = basePath
	p.SourceDir = path.Join(basePath,"src")
	p.LibDir = path.Join(basePath, "lib")

	return nil
}

func (p *PioProject) Create(board string) error {

	if _, err := os.Stat(p.RootDir); err == nil {
		return fmt.Errorf("Cannot create project, directory '%s' already exists", p.RootDir)
	}

	os.MkdirAll(p.RootDir, os.ModePerm)

	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(currentDir)

	os.Chdir(p.RootDir)

	cmd := exec.Command("platformio", "init", "--board", board)
	//cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (p *PioProject) Open() error {

	info, err := os.Stat(p.RootDir)

	if err != nil || !info.IsDir() {
		return fmt.Errorf("Cannot open project, directory '%s' doesn't exists", p.RootDir)
	}

	if _, err := os.Stat(path.Join(p.RootDir,"platformio.ini")); os.IsNotExist(err) {
		return errors.New("Invalid device project, platformio.ini doesn't exists")
	}

	return nil
}

func (*PioProject) InstallLib(name string, id int) error {
	cmd := exec.Command("platformio", "lib", "install", strconv.Itoa(id))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (*PioProject) Build() error {
	cmd := exec.Command("platformio", "run")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (*PioProject) Upload() error {
	cmd := exec.Command("platformio", "run", "--target", "upload")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (*PioProject) Clean() error {
	cmd := exec.Command("platformio", "run", "--target", "clean")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
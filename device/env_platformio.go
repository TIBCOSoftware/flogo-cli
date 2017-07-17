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
	RootDir         string
	LibDir          string
	SourceDir       string
	ContributionDir string
}

func NewPlatformIoProject() Project {

	project := &PioProject{}

	return project
}

func (p *PioProject) Init(basePath string) error {

	_, err := exec.LookPath("platformio")

	if err != nil {
		return errors.New("platformio not installed")
	}

	p.RootDir = basePath
	p.SourceDir = path.Join(basePath,"src")
	p.LibDir = path.Join(basePath, "lib")
	p.ContributionDir = path.Join(basePath, "vendor", "src")
	return nil
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

func (p *PioProject) GetContributionDir() string {
	return p.ContributionDir
}

func (p *PioProject) Create() error {

	if _, err := os.Stat(p.RootDir); err == nil {
		return fmt.Errorf("Cannot create project, directory '%s' already exists", p.RootDir)
	}

	os.MkdirAll(p.RootDir, os.ModePerm)
	os.MkdirAll(p.SourceDir, os.ModePerm)

	//currentDir, err := os.Getwd()
	//if err != nil {
	//	return err
	//}
	//defer os.Chdir(currentDir)
	//
	//os.Chdir(p.RootDir)
	//
	//cmd := exec.Command("platformio", "init", "--board", board)
	////cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	//
	//return cmd.Run()

	return nil
}

func (p *PioProject) Setup(board string) error {

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

func (p *PioProject) InstallLib(name string, id int) error {

	currentDir, _ := os.Getwd()
	defer os.Chdir(currentDir)

	os.Chdir(p.RootDir)

	cmd := exec.Command("platformio", "lib", "install", strconv.Itoa(id))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (p *PioProject) InstallContribution(depPath string, version string) error {
	var cmd *exec.Cmd

	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)

	//check if dependency is installed
	if _, err := os.Stat(path.Join(p.ContributionDir, depPath)); err == nil {
		//todo ignore installed dependencies for now
		//exists, return
		return nil
	}

	if version == "" {
		//cmd = exec.Command("gb", "vendor", "fetch", "-branch", "device_contribs", depPath)
		cmd = exec.Command("gb", "vendor", "fetch", depPath)
	} else {
		var tag string

		if version[0] != 'v' {
			tag = "v" + version
		} else {
			tag = version
		}

		//cmd = exec.Command("gb", "vendor", "fetch", "-branch", "device_contribs", "-tag", tag, depPath)
		cmd = exec.Command("gb", "vendor", "fetch", "-tag", tag, depPath)
	}

	os.Chdir(p.RootDir)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (p *PioProject) UninstallContribution(depPath string) error {

	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)

	//check if dependency is installed
	if _, err := os.Stat(path.Join(p.ContributionDir, depPath)); err != nil {
		//todo ignore dependencies that are not installed for now
		//exists, return
		return nil
	}

	os.Chdir(p.RootDir)

	cmd := exec.Command("gb", "vendor", "delete", depPath)

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
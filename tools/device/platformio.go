package device

import (
	"os/exec"
	"os"
	"strconv"
)

func PioInstalled() bool {
	_, err := exec.LookPath("platformio")
	return err == nil;
}

func PioInit(board string) error {
	cmd := exec.Command("platformio", "init", "--board", board)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func PioIsProject() bool {
	if _, err := os.Stat("platformio.ini"); os.IsNotExist(err) {
		return false
	}

	return true
}

func PioDirIsProject(dir string) bool {
	if _, err := os.Stat(dir + "/platformio.ini"); os.IsNotExist(err) {
		return false
	}

	return true
}

func PioInstallLib(libId int) error {
	cmd := exec.Command("platformio", "lib", "install", strconv.Itoa(libId))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func PioBuild() error {
	cmd := exec.Command("platformio", "run")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func PioUpload() error {
	cmd := exec.Command("platformio", "run", "--target", "upload")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func PioClean() error {
	cmd := exec.Command("platformio", "run", "--target", "clean")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
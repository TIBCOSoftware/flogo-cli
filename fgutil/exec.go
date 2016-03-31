package fgutil

import (
	"os/exec"
	"os"
)

// ExecutableExists checks if the specified executable
// exists in the users PATH
func ExecutableExists(execName string) bool {
	_, err := exec.LookPath(execName)
	if err != nil {
		return false
	}
	return true
}

func FileExists(filePath string) bool {

	f, err := os.Open(filePath)
	f.Close()
	if err != nil {
		return false
	}
	return true
}
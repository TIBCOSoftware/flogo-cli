package fgutil

import (
	"os"
	"os/exec"
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

// FileExists determines if the named file exists
func FileExists(filePath string) bool {

	f, err := os.Open(filePath)
	f.Close()
	if err != nil {
		return false
	}
	return true
}

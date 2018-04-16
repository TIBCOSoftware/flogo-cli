//+build ignore
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

func main() {
	ldf := flags()
	cmd := exec.Command("go", "install", "-ldflags="+ldf, "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "TIBCOSoftware", "flogo-cli")
	cmd.Env = append(os.Environ())

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// tag returns the git tag for the current branch or "" if none.
func tag() string {
	out, err := exec.Command("git", "describe", "--tags").Output()
	if err != nil {
		log.Fatal(err)
	}
	re := regexp.MustCompile("\\n")
	tag := re.ReplaceAllString(string(out), "")
	return tag
}

// hash returns the git hash for the current repo or "" if none.
func hash() string {
	out, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	if err != nil {
		log.Fatal(err)
	}
	re := regexp.MustCompile("\\n")
	hash := re.ReplaceAllString(string(out), "")
	return hash
}

func flags() string {
	hash := hash()
	tag := tag()
	return fmt.Sprintf(`-X "github.com/TIBCOSoftware/flogo-cli/app.tag=%s" -X "github.com/TIBCOSoftware/flogo-cli/app.hash=%s"`, tag, hash)
}

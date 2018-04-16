package app

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/TIBCOSoftware/flogo-cli/cli"
)

var optVersion = &cli.OptionInfo{
	Name:      "version",
	UsageLine: "version",
	Short:     "displays the version of flogo cli and flogo-contrib",
	Long: `Get the current version number of the cli and contrib.

`,
}

var tag = ""
var hash = ""

func init() {
	CommandRegistry.RegisterCommand(&cmdVersion{option: optVersion})
}

type cmdVersion struct {
	option *cli.OptionInfo
}

// HasOptionInfo implementation of cli.HasOptionInfo.OptionInfo
func (c *cmdVersion) OptionInfo() *cli.OptionInfo {
	return c.option
}

// Exec implementation of cli.Command.Exec
func (c *cmdVersion) AddFlags(fs *flag.FlagSet) {
	//op op
}

// Exec implementation of cli.Command.Exec
func (c *cmdVersion) Exec(args []string) error {

	line := fmt.Sprintf("flogo cli version [%s] and commithash [%s]\n", tag, hash)
	fmt.Fprint(os.Stdout, line)

	cmd := exec.Command("git", "describe", "--tags")
	cmd.Dir = filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "TIBCOSoftware", "flogo-contrib")
	cmd.Env = append(os.Environ())

	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	re := regexp.MustCompile("\\n")
	fc := re.ReplaceAllString(string(out), "")

	line = fmt.Sprintf("flogo-contrib version [%s]\n\n", fc)
	fmt.Fprint(os.Stdout, line)

	return nil
}

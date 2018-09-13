package app

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/TIBCOSoftware/flogo-cli/cli"
	toml "github.com/pelletier/go-toml"
)

var optVersion = &cli.OptionInfo{
	Name:      "version",
	UsageLine: "version",
	Short:     "displays the version of flogo cli",
	Long: `Get the current version number of the cli.

`,
}

type rawLock struct {
	SolveMeta solveMeta          `toml:"solve-meta"`
	Projects  []rawLockedProject `toml:"projects"`
}

type solveMeta struct {
	InputsDigest    string `toml:"inputs-digest"`
	AnalyzerName    string `toml:"analyzer-name"`
	AnalyzerVersion int    `toml:"analyzer-version"`
	SolverName      string `toml:"solver-name"`
	SolverVersion   int    `toml:"solver-version"`
}

type rawLockedProject struct {
	Name     string   `toml:"name"`
	Branch   string   `toml:"branch,omitempty"`
	Revision string   `toml:"revision"`
	Version  string   `toml:"version,omitempty"`
	Source   string   `toml:"source,omitempty"`
	Packages []string `toml:"packages"`
}

const lockName = "Gopkg.lock"

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

	cmd := exec.Command("git", "describe", "--tags")
	gopath, set := os.LookupEnv("GOPATH")
	if !set {
		out, err := exec.Command("go", "env", "GOPATH").Output()
		if err != nil {
			log.Fatal(err)
		}
		gopath = strings.TrimSuffix(string(out), "\n")
	}
	cmd.Dir = filepath.Join(gopath, "src", "github.com", "TIBCOSoftware", "flogo-cli")
	cmd.Env = append(os.Environ())

	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	re := regexp.MustCompile("\\n")
	fc := re.ReplaceAllString(string(out), "")

	line := fmt.Sprintf("flogo cli version [%s]\n", fc)
	fmt.Fprint(os.Stdout, line)

	appDir, err := os.Getwd()

	if err != nil {
		fmt.Fprint(os.Stderr, "Error: Unable to determine working directory\n\n")
		os.Exit(2)
	}

	project, _ := SetupExistingProjectEnvWithOrWithoutExitOnFailure(appDir, false)

	if project != nil {
		config, err := toml.LoadFile(filepath.Join(project.GetAppDir(), lockName))

		if err != nil {
			fmt.Println("Error ", err.Error())
		} else {
			raw := rawLock{}
			err := config.Unmarshal(&raw)
			if err != nil {
				fmt.Printf("Unable to parse the lock as TOML")
			}

			for _, v := range raw.Projects {
				if caseInsensitiveContains(v.Name, "flogo") {
					if v.Version == "" {
						line = fmt.Sprintf("Your project uses %s branch %s and revision %s\n", v.Name, v.Branch, v.Revision)
					} else {
						line = fmt.Sprintf("Your project uses %s version %s\n", v.Name, v.Version)
					}
					fmt.Fprint(os.Stdout, line)
				}
			}
		}
	}

	return nil
}

// This isn't the most performant way, but this will be able to check if the string exists while ignoring any case sensitivity.
func caseInsensitiveContains(s, substr string) bool {
	s, substr = strings.ToUpper(s), strings.ToUpper(substr)
	return strings.Contains(s, substr)
}

package model

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/TIBCOSoftware/flogo-cli/cli"
	"github.com/TIBCOSoftware/flogo-cli/util"
)

var optCreate = &cli.OptionInfo{
	Name:      "create",
	UsageLine: "create [-no_gb] modelName",
	Short:     "creae a model project",
	Long: `Creates a flogo model project.

Options:
    -no_gb       generate without gb structure

`,
}

func init() {
	Tool().CommandRegistry().RegisterCommand(&cmdCreate{option: optCreate})
}

type cmdCreate struct {
	option *cli.OptionInfo
	noGB   bool
}

func (c *cmdCreate) OptionInfo() *cli.OptionInfo {
	return c.option
}

func (c *cmdCreate) AddFlags(fs *flag.FlagSet) {
	fs.BoolVar(&(c.noGB), "nogb", false, "generate without gb structure")
}

func (c *cmdCreate) Exec(args []string) error {

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: Model name not specified\n\n")
		Tool().CmdUsage(c)
	}

	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Error: Too many arguments given\n\n")
		Tool().CmdUsage(c)
	}

	modelName := args[0]

	if _, err := os.Stat(modelName); err == nil {
		fmt.Fprintf(os.Stderr, "Error: Cannot create model project, directory '%s' already exists\n\n", modelName)
		os.Exit(2)
	}

	fmt.Fprintf(os.Stdout, "Creating flogo model '%s'...\n", modelName)
	os.MkdirAll(modelName, 0777)
	os.Chdir(modelName)

	//var srcPath string
	var codeSrcPath string

	if !c.noGB {

		gb := fgutil.NewGb(modelName)

		if !gb.Installed() {
			fmt.Fprintf(os.Stderr, "Error: Cannot create model project [%s] using 'gb', gb is not installed\n\n", modelName)
			os.Exit(2)
		}

		gb.Init(false)

		err := gb.VendorFetch("github.com/TIBCOSoftware/flogo-lib")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(2)
		}

		//srcPath = gb.SourcePath
		codeSrcPath = gb.CodeSourcePath

	} else {
		//srcPath = ""
		codeSrcPath = ""
	}

	data := struct {
		Name string
	}{
		modelName,
	}

	//todo revisit model project layout

	createProjectDescriptor(codeSrcPath, data)
	createModelGoFile(codeSrcPath, data)
	createModelTestGoFile(codeSrcPath, data)

	return nil
}

func path(parts ...string) string {
	return strings.Join(parts[:], string(os.PathSeparator))
}

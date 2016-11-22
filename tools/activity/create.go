package activity

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
	UsageLine: "create [-no_gb] activityName",
	Short:     "create an activity project",
	Long: `Creates a flogo activity project.

Options:
    -no_gb    generate without gb structure
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
		fmt.Fprintf(os.Stderr, "Error: Activity name not specified\n\n")
		Tool().CmdUsage(c)
	}

	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Error: Too many arguments given\n\n")
		Tool().CmdUsage(c)
	}

	activityName := args[0]

	if _, err := os.Stat(activityName); err == nil {
		fmt.Fprintf(os.Stderr, "Error: Cannot create activity project, directory '%s' already exists\n\n", activityName)
		os.Exit(2)
	}

	fmt.Fprintf(os.Stdout, "Creating flogo activity '%s'...\n", activityName)
	os.MkdirAll(activityName, 0777)
	os.Chdir(activityName)

	var srcPath string
	var codeSrcPath string

	if !c.noGB {

		gb := fgutil.NewGb(dirRT)

		if !gb.Installed() {
			fmt.Fprintf(os.Stderr, "Error: Cannot create activity project [%s], gb is not installed\n\n", activityName)
			os.Exit(2)
		}

		gb.Init(false)

		//todo should we add the ability to specify the flogo-lib version
		err := gb.VendorFetch("github.com/TIBCOSoftware/flogo-lib", "")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(2)
		}

		srcPath = gb.SourcePath
		codeSrcPath = gb.CodeSourcePath

	} else {
		os.MkdirAll(dirRT, 0777)
		srcPath = ""
		codeSrcPath = dirRT

	}

	data := struct {
		Name string
	}{
		activityName,
	}

	createProjectDescriptor(srcPath, data)
	createActivityGoFile(codeSrcPath, data)
	createActivityTestGoFile(codeSrcPath, data)
	createMetadataGoFile(codeSrcPath, data)

	return nil
}

func path(parts ...string) string {
	return strings.Join(parts[:], string(os.PathSeparator))
}

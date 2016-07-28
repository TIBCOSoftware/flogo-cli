package trigger

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
	UsageLine: "create [-no_ui] [-no_gb] triggerName",
	Short:     "create an trigger project",
	Long: `Creates a flogo trigger project.

Options:
    -no_ui    generate trigger ui
    -no_gb       generate without gb structure

`,
}

func init() {
	Tool().CommandRegistry().RegisterCommand(&cmdCreate{option: optCreate})
}

type cmdCreate struct {
	option *cli.OptionInfo
	noUI   bool
	noGB   bool
}

func (c *cmdCreate) OptionInfo() *cli.OptionInfo {
	return c.option
}

func (c *cmdCreate) AddFlags(fs *flag.FlagSet) {
	fs.BoolVar(&(c.noUI), "no-ui", false, "generate ui components")
	fs.BoolVar(&(c.noGB), "nogb", false, "generate without gb structure")
}

func (c *cmdCreate) Exec(args []string) error {

	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Error: Trigger name not specified\n\n")
		Tool().CmdUsage(c)
	}

	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Error: Too many arguments given\n\n")
		Tool().CmdUsage(c)
	}

	triggerName := args[0]

	if _, err := os.Stat(triggerName); err == nil {
		fmt.Fprintf(os.Stderr, "Error: Cannot create trigger project, directory '%s' already exists\n\n", triggerName)
		os.Exit(2)
	}

	fmt.Fprintf(os.Stdout, "Creating flogo trigger '%s'...\n", triggerName)
	os.MkdirAll(triggerName, 0777)
	os.Chdir(triggerName)

	var srcPath string
	var codeSrcPath string

	if !c.noGB {

		gb := fgutil.NewGb(dirRT)

		if !gb.Installed() {
			fmt.Fprintf(os.Stderr, "Error: Cannot create trigger project [%s] using 'gb', gb is not installed\n\n", triggerName)
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

		if !c.noUI {
			// create ui directory
			os.Mkdir(path(gb.SourcePath, dirUI), 0777)
		}

	} else {
		os.MkdirAll(dirRT, 0777)
		srcPath = ""
		codeSrcPath = dirRT

		if !c.noUI {
			// create ui directory
			os.Mkdir(dirUI, 0777)
		}
	}

	data := struct {
		Name string
	}{
		triggerName,
	}

	createProjectDescriptor(srcPath, data)
	createTriggerGoFile(codeSrcPath, data)
	createTriggerTestGoFile(codeSrcPath, data)
	createMetadataGoFile(codeSrcPath, data)

	return nil
}

func path(parts ...string) string {
	return strings.Join(parts[:], string(os.PathSeparator))
}

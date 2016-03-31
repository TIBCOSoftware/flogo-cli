package trigger

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/TIBCOSoftware/flogo-tools/fg"
	"github.com/TIBCOSoftware/flogo-tools/fgutil"
)

const dirDT string = "dt"
const dirRT string = "rt"
const fileDescriptor string = "trigger.json"
const fileTriggerGo string = "trigger.go"
const fileTriggerTestGo string = "trigger_test.go"
const fileTriggerMdGo string = "trigger_metadata.go"

var optCreate = &flogo.OptionInfo{
	Name:      "create",
	UsageLine: "create [-with-ui] [-gb] triggerName",
	Short:     "create an trigger project",
	Long: `Creates a flogo trigger project.

Options:
    -with-ui    generate trigger ui
    -gb         generate within gb structure
`,
}

func init() {
	Tool().CommandRegistry().RegisterCommand(&cmdCreate{option: optCreate})
}

type cmdCreate struct {
	option *flogo.OptionInfo
	withUI bool
	useGB  bool
}

func (c *cmdCreate) OptionInfo() *flogo.OptionInfo {
	return c.option
}

func (c *cmdCreate) AddFlags(fs *flag.FlagSet) {
	fs.BoolVar(&(c.withUI), "with-ui", false, "generate ui components")
	fs.BoolVar(&(c.useGB), "gb", false, "generate within gb structure")
}

func (c *cmdCreate) Exec(ctx *flogo.Context, args []string) error {

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

	basePath := triggerName

	if c.useGB {

		if !fgutil.ExecutableExists("gb") {
			fmt.Fprintf(os.Stderr, "Error: Cannot create trigger project [%s] using 'gb', gb is not installed\n\n", triggerName)
			os.Exit(2)
		}

		fmt.Fprintf(os.Stdout, "Creating flogo trigger '%s'...\n", triggerName)

		basePath = path(triggerName, "src")
		os.MkdirAll(path(basePath), 0777)
		os.MkdirAll(path(triggerName, "vendor", "src"), 0777)

		os.Chdir(triggerName)
		fmt.Fprint(os.Stdout, "Installing flogo lib...\n")
		cmd := exec.Command("gb", "vendor", "fetch", "github.com/TIBCOSoftware/flogo-lib")
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(2)
		}
		os.Chdir("..")
	} else {
		fmt.Fprintf(os.Stdout, "Creating flogo trigger '%s'...\n", triggerName)
		os.Mkdir(triggerName, 0777)
	}

	data := struct {
		Name string
	}{
		triggerName,
	}

	// create trigger.json file
	f, _ := os.Create(path(basePath, fileDescriptor))
	fgutil.RenderTemplate(f, tplTriggerJSON, data)
	f.Close()

	// create rt directory
	rtDir := path(basePath, dirRT)
	os.Mkdir(rtDir, 0777)

	// create trigger Go file
	f, _ = os.Create(path(rtDir, fileTriggerGo))
	fgutil.RenderTemplate(f, tplTriggerGoFile, data)
	f.Close()

	// create trigger test Go file
	f, _ = os.Create(path(rtDir, fileTriggerTestGo))
	fgutil.RenderTemplate(f, tplTriggerTestGoFile, data)
	f.Close()

	// create trigger metadata Go file
	f, _ = os.Create(path(rtDir, fileTriggerMdGo))
	fgutil.RenderTemplate(f, tplMetadataGoFile, data)
	f.Close()

	if c.withUI {
		// create dt directory
		os.Mkdir(path(basePath, dirDT), 0777)
	}

	return nil
}

func path(parts ...string) string {
	return strings.Join(parts[:], string(os.PathSeparator))
}

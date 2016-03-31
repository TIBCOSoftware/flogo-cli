package activity

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
const fileDescriptor string = "activity.json"
const fileActivityGo string = "activity.go"
const fileActivityTestGo string = "activity_test.go"
const fileActivityMdGo string = "activity_metadata.go"

var optCreate = &flogo.OptionInfo{
	Name:      "create",
	UsageLine: "create [-with-ui] [-gb] activityName",
	Short:     "create an activity project",
	Long: `Creates a flogo activity project.

Options:
    -with-ui    generate activity ui
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

	basePath := activityName

	if c.useGB {

		if !fgutil.ExecutableExists("gb") {
			fmt.Fprintf(os.Stderr, "Error: Cannot create activity project [%s] using 'gb', gb is not installed\n\n", activityName)
			os.Exit(2)
		}

		fmt.Fprintf(os.Stdout, "Creating flogo activity '%s'...\n", activityName)

		basePath = path(activityName, "src")
		os.MkdirAll(path(basePath), 0777)
		os.MkdirAll(path(activityName, "vendor", "src"), 0777)

		os.Chdir(activityName)
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
		fmt.Fprintf(os.Stdout, "Creating flogo activity '%s'...\n", activityName)
		os.Mkdir(activityName, 0777)
	}

	data := struct {
		Name string
	}{
		activityName,
	}

	// create activity.json file
	f, _ := os.Create(path(basePath, fileDescriptor))
	fgutil.RenderTemplate(f, tplActivityJSON, data)
	f.Close()

	// create rt directory
	rtDir := path(basePath, dirRT)
	os.Mkdir(rtDir, 0777)

	// create activity Go file
	f, _ = os.Create(path(rtDir, fileActivityGo))
	fgutil.RenderTemplate(f, tplActivityGoFile, data)
	f.Close()

	// create activity test Go file
	f, _ = os.Create(path(rtDir, fileActivityTestGo))
	fgutil.RenderTemplate(f, tplActivityTestGoFile, data)
	f.Close()

	// create activity metadata Go file
	f, _ = os.Create(path(rtDir, fileActivityMdGo))
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

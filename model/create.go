package model

import (
	"fg"
	"fgutil"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const fileDescriptor string = "model.json"
const fileModelGo string = "model.go"
const fileModelTestGo string = "model_test.go"

var optCreate = &flogo.OptionInfo{
	Name:      "create",
	UsageLine: "create [-gb] modelName",
	Short:     "creae a model project",
	Long: `Creates a flogo model project.

Options:
    -gb         generate within gb structure

`,
}

func init() {
	Tool().CommandRegistry().RegisterCommand(&cmdCreate{option: optCreate})
}

type cmdCreate struct {
	option *flogo.OptionInfo
	useGB  bool
}

func (c *cmdCreate) OptionInfo() *flogo.OptionInfo {
	return c.option
}

func (c *cmdCreate) AddFlags(fs *flag.FlagSet) {
	fs.BoolVar(&(c.useGB), "gb", false, "generate within gb structure")
}

func (c *cmdCreate) Exec(ctx *flogo.Context, args []string) error {

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

	basePath := modelName

	if c.useGB {

		if !fgutil.ExecutableExists("gb") {
			fmt.Fprintf(os.Stderr, "Error: Cannot create model project [%s] using 'gb', gb is not installed\n\n", modelName)
			os.Exit(2)
		}

		fmt.Fprintf(os.Stdout, "Creating flogo model '%s'...\n", modelName)

		basePath = path(modelName, "src")
		os.MkdirAll(path(basePath), 0777)
		os.MkdirAll(path(modelName, "vendor", "src"), 0777)

		os.Chdir(modelName)
		fmt.Fprint(os.Stdout, "Installing flogo lib...\n")
		cmd := exec.Command("gb", "vendor", "fetch", "github.com/TIBCOSoftware/flogo/golib")
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(2)
		}
		os.Chdir("..")
	} else {
		fmt.Fprintf(os.Stdout, "Creating flogo model '%s'...\n", modelName)
		os.Mkdir(modelName, 0777)
	}

	data := struct {
		Name string
	}{
		modelName,
	}

	// create model.json file
	f, _ := os.Create(path(basePath, fileDescriptor))
	fgutil.RenderTemplate(f, tplModelJSON, data)
	f.Close()

	// create model Go file
	f, _ = os.Create(path(basePath, fileModelGo))
	fgutil.RenderTemplate(f, tplModelGoFile, data)
	f.Close()

	// create model test Go file
	f, _ = os.Create(path(basePath, fileModelTestGo))
	fgutil.RenderTemplate(f, tplModelTestGoFile, data)
	f.Close()

	return nil
}

func path(parts ...string) string {
	return strings.Join(parts[:], string(os.PathSeparator))
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/TIBCOSoftware/flogo/cli"
	"github.com/TIBCOSoftware/flogo/util"
)

var optList = &cli.OptionInfo{
	Name:      "list",
	UsageLine: "list [activity|model|trigger]",
	Short:     "list objects configured on the flogo project",
	Long: `List the objects the flogo project has been configured with.
`,
}

func init() {
	commandRegistry.RegisterCommand(&cmdList{option: optList})
}

type cmdList struct {
	option *cli.OptionInfo
}

func (c *cmdList) OptionInfo() *cli.OptionInfo {
	return c.option
}

func (c *cmdList) AddFlags(fs *flag.FlagSet) {
	//op op
}

func (c *cmdList) Exec(args []string) error {

	projectDescriptor := loadProjectDescriptor()

	if len(args) > 1 {
		fmt.Fprintf(os.Stderr, "Error: Too many arguments given\n\n")
		cmdUsage(c)
	}

	var tpl string

	if len(args) == 0 {

		tpl = tplListAll

	} else {

		itemType := args[0]

		switch itemType {
		case itActivity:
			tpl = tplListActivities
		case itModel:
			tpl = tplListModels
		case itTrigger:
			tpl = tplListTriggers
		default:
			fmt.Fprintf(os.Stderr, "Error: Unknown item type '%s'\n\n", itemType)
			os.Exit(2)
		}
	}

	bw := bufio.NewWriter(os.Stdout)
	fgutil.RenderTemplate(bw, tpl, projectDescriptor)
	bw.Flush()

	return nil
}

var tplListAll = `

Activities:
{{range .Activities}}
    - {{.Name}} [{{.Path}}]{{end}}

Triggers:
{{range .Triggers}}
    - {{.Name}} [{{.Path}}]{{end}}

Models:
{{range .Models}}
    - {{.Name}} [{{.Path}}]{{end}}

`

var tplListActivities = `

Activities:
{{range .Activities}}
    - {{.Name}} [{{.Path}}]{{end}}
`

var tplListTriggers = `

Triggers:
{{range .Triggers}}
    - {{.Name}} [{{.Path}}]{{end}}

`

var tplListModels = `

Models:
{{range .Models}}
    - {{.Name}} [{{.Path}}]{{end}}

`

package engine

import (
	"flag"
	"fmt"
	"os"

	"github.com/TIBCOSoftware/flogo-tools/fg"
	"github.com/TIBCOSoftware/flogo-tools/fgutil"
	"encoding/json"
	"bufio"
)

var optList = &flogo.OptionInfo{
	Name:      "list",
	UsageLine: "list [activity|model|trigger]",
	Short:     "list objects configured on the engine project",
	Long: `List the object the engine project has been configured with.
`,
}

func init() {
	Tool().CommandRegistry().RegisterCommand(&cmdList{option: optList})
}

type cmdList struct {
	option *flogo.OptionInfo
}

func (c *cmdList) OptionInfo() *flogo.OptionInfo {
	return c.option
}

func (c *cmdList) AddFlags(fs *flag.FlagSet) {
	//op op
}

func (c *cmdList) Exec(ctx *flogo.Context, args []string) error {

	configFile, err := os.Open(fileProjectConfig)

	if err != nil {
		fmt.Fprint(os.Stderr, "Error: Current working directory is not a flogo-based engine project.\n\n")
		os.Exit(2)
	}

	if len(args) > 1 {
		fmt.Fprintf(os.Stderr, "Error: Too many arguments given\n\n")
		Tool().CmdUsage(c)
	}

	engineConfig := &EngineProjectConfig{}
	jsonParser := json.NewDecoder(configFile)

	if err = jsonParser.Decode(engineConfig); err != nil {
		fmt.Fprint(os.Stderr, "Error: Unable to parse engine.json, file may be corrupted.\n\n")
		os.Exit(2)
	}

	configFile.Close()

	var tpl string

	if len(args) == 0 {

		tpl = tplListAll

	} else {

		itemType := args[0]

		switch itemType {
		case "activity":
			tpl = tplListActivities
		case "model":
			tpl = tplListModels
		case "trigger":
			tpl = tplListTriggers
		default:
			fmt.Fprintf(os.Stderr, "Error: Unknown item type '%s'\n\n", itemType)
			os.Exit(2)
		}
	}

	bw := bufio.NewWriter(os.Stdout)
	fgutil.RenderTemplate(bw, tpl, engineConfig)
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
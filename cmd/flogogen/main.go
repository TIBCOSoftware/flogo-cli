package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/TIBCOSoftware/flogo-cli/gen"
	"github.com/TIBCOSoftware/flogo-cli/util"
	"path"
)

var (
	generators = make(map[string]gen.CodeGenerator)
)

func init() {
	generators["action"] = &gen.ActionGenerator{}
	generators["trigger"] = &gen.TriggerGenerator{}
	generators["activity"] = &gen.ActivityGenerator{}
	generators["flowmodel"] = &gen.FlowModelGenerator{}
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "FATAL: "+format+"\n", args...)
	os.Exit(1)
}

func main() {
	args := os.Args
	if len(args) != 3 || args[1] == "-h" {
		usage()
	}

	contribution := args[1]
	name := args[2]

	generator, exists := generators[contribution]

	if exists {

		data := struct {
			Name string
		}{
			name,
		}

		currentDir, _ := os.Getwd()
		basePath := path.Join(currentDir, name)

		if _, err := os.Stat(basePath); err == nil {
			fmt.Fprintf(os.Stderr, "Error: Cannot create project, directory '%s' already exists\n\n", name)
			os.Exit(2)
		}

		os.MkdirAll(basePath, 0777)

		err := generator.Generate(basePath, data)

		if err != nil{
			fmt.Fprintf(os.Stderr, "Error generating contribution: %s\n\n", err.Error())
			os.Exit(2)
		}

	} else {
		fmt.Fprintf(os.Stderr, "FATAL: unknown contribution type %q\n\n", contribution)
		usage()
	}

	os.Exit(0)
}

func usage() {
	printUsage(os.Stderr)
	os.Exit(2)
}

func printUsage(w io.Writer) {
	bw := bufio.NewWriter(w)

	fgutil.RenderTemplate(bw, usageTpl, generators)
	bw.Flush()
}

var usageTpl = `Usage:

    flogogen <contribution> name

  contributions:
{{ range $key, $value := . }}
    {{$key | printf "%-12s"}} {{$value.Description}}{{end}}

`

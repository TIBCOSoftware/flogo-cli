package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/TIBCOSoftware/flogo/cli"
	"github.com/TIBCOSoftware/flogo/util"
)

var optBuild = &cli.OptionInfo{
	Name:      "build",
	UsageLine: "build [-validate]",
	Short:     "build the flogo application",
	Long: `Build the flogo application.

Options:
    -validate   validate that the project is buildable
`,
}

const fileFlowsGo string = "flows.go"

func init() {
	commandRegistry.RegisterCommand(&cmdBuild{option: optBuild})
}

type cmdBuild struct {
	option   *cli.OptionInfo
	validate bool
}

func (c *cmdBuild) OptionInfo() *cli.OptionInfo {
	return c.option
}

func (c *cmdBuild) AddFlags(fs *flag.FlagSet) {
	fs.BoolVar(&(c.validate), "validate", false, "only validate if buildable")
}

func (c *cmdBuild) Exec(args []string) error {

	projectConfig := loadProjectConfig()

	if len(args) > 1 {
		fmt.Fprintf(os.Stderr, "Error: Too many arguments given\n\n")
		cmdUsage(c)
	}

	gb := fgutil.NewGb(projectConfig.Name)

	flows := importFlows(dirFlows)

	createFlowsGoFile(gb.CodeSourcePath, flows)

	if len(projectConfig.Models) == 0 {
		fmt.Fprint(os.Stderr, "Error: Project must have a least one model.\n\n")
		os.Exit(2)
	}

	if len(projectConfig.Triggers) == 0 {
		fmt.Fprint(os.Stderr, "Error: Project must have a least one trigger.\n\n")
		os.Exit(2)
	}

	if c.validate {
		return nil
	}

	err := gb.Build()
	if err != nil {
		os.Exit(2)
	}

	return nil
}

func importFlows(flowDir string) map[string]string {

	flows := make(map[string]string)

	fileInfos, err := ioutil.ReadDir(flowDir)

	if err == nil {

		for _, fileInfo := range fileInfos {

			if !fileInfo.IsDir() {

				fileName := fileInfo.Name()

				// validate flow json
				flowFilePath := path(flowDir, fileName)
				b64 := gzipAndB64(flowFilePath) //todo: is gzip necessary
				idx := strings.Index(fileName, ".")

				flows[fileName[:idx]] = b64
			}
		}
	}

	return flows
}

func gzipAndB64(flowFilePath string) string {

	in, err := os.Open(flowFilePath)
	if err != nil {
		//log.Fatal(err)
	}

	bufin := bufio.NewReader(in)

	var b bytes.Buffer
	gz, err := gzip.NewWriterLevel(&b, gzip.BestCompression)
	_, err = bufin.WriteTo(gz)

	if err != nil {
		panic(err)
	}

	if err := gz.Flush(); err != nil {
		panic(err)
	}
	if err := gz.Close(); err != nil {
		panic(err)
	}

	in.Close()

	return base64.StdEncoding.EncodeToString(b.Bytes())
}

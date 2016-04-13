package engine

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/TIBCOSoftware/flogo/fg"
)

var optBuild = &flogo.OptionInfo{
	Name:      "build",
	UsageLine: "build [-validate]",
	Short:     "build the engine using gb",
	Long: `Build the engine project using gb.

Options:
    -validate   validate if the engine is buildable
`,
}

const fileFlowsGo string = "flows.go"

func init() {
	Tool().CommandRegistry().RegisterCommand(&cmdBuild{option: optBuild})
}

type cmdBuild struct {
	option   *flogo.OptionInfo
	validate bool
}

func (c *cmdBuild) OptionInfo() *flogo.OptionInfo {
	return c.option
}

func (c *cmdBuild) AddFlags(fs *flag.FlagSet) {
	fs.BoolVar(&(c.validate), "validate", false, "only validate if buildable")
}

func (c *cmdBuild) Exec(ctx *flogo.Context, args []string) error {

	configFile, err := os.Open(fileProjectConfig)

	if err != nil {
		fmt.Fprint(os.Stderr, "Error: Current working directory is not a flogo-based engine project.\n\n")
		os.Exit(2)
	}

	if len(args) > 1 {
		fmt.Fprintf(os.Stderr, "Error: Too many arguments given\n\n")
		Tool().CmdUsage(c)
	}

	projectConfig := &EngineProjectConfig{}
	jsonParser := json.NewDecoder(configFile)

	if err = jsonParser.Decode(projectConfig); err != nil {
		fmt.Fprint(os.Stderr, "Error: Unable to parse engine.json, file may be corrupted.\n\n")
		os.Exit(2)
	}

	configFile.Close()

	flowsPath := "./flows"

	flows := importFlows(flowsPath)

	createFlowsGoFile(path("src", projectConfig.Name), flows)

	if len(projectConfig.Models) == 0 {
		fmt.Fprint(os.Stderr, "Error: Engine must have a least one model.\n\n")
		os.Exit(2)
	}

	if len(projectConfig.Triggers) == 0 {
		fmt.Fprint(os.Stderr, "Error: Engine must have a least one trigger.\n\n")
		os.Exit(2)
	}

	if c.validate {
		return nil
	}

	cmd := exec.Command("gb", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		os.Exit(2)
	}

	return nil
}

func importFlows(flowDir string) map[string]string {

	fmt.Println("flows dir: ", flowDir)

	flows := make(map[string]string)

	fileInfos, err := ioutil.ReadDir(flowDir)

	if err == nil {

		for _, fileInfo := range fileInfos {

			fmt.Println("file/dir: ", fileInfo.Name())

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

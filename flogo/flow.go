package main

import (
	"github.com/xeipuuv/gojsonschema"
	"os"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"bufio"
	"bytes"
	"encoding/base64"
	"strings"
	"compress/gzip"
)

func ImportFlows(projectDescriptor *FlogoProjectDescriptor, flowDir string) map[string]string {

	flows := make(map[string]string)

	fileInfos, err := ioutil.ReadDir(flowDir)

	if err == nil {

		for _, fileInfo := range fileInfos {

			if !fileInfo.IsDir() {

				fileName := fileInfo.Name()

				// validate flow json
				flowFilePath := path(flowDir, fileName)

				ValidateFlow(projectDescriptor, flowFilePath)

				b64 := gzipAndB64(flowFilePath) //todo: is gzip necessary

				flows[genFlowURI(fileName)] = b64
			}
		}
	}

	return flows
}

func genFlowURI(fileName string) string {

	idx := strings.LastIndex(fileName, ".")
	return "local://" + fileName[:idx]
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

func ValidateFlow(projectDescriptor *FlogoProjectDescriptor, flowPath string) {

	// first validate the flow json
	validateFlowSchema(flowPath)

	//next check if all activities used int he flow are installed in engine

	file, _ := ioutil.ReadFile(flowPath)

	var flowObj interface{}
	json.Unmarshal(file, &flowObj)

	activityTypes := make(map[string]bool)

	getActivityTypes(flowObj, activityTypes)

	for _, desc := range projectDescriptor.Activities {
		delete(activityTypes, desc.Name)
	}

	if len(activityTypes) > 0 {
		fmt.Fprintf(os.Stderr, "Error: cannot embed '%s', the activites required to run the flow have not been added to your project\n", flowPath)

		for k, _ := range activityTypes {
			fmt.Fprintf(os.Stderr, "    MISSING: %s\n", k)
		}

		os.Exit(2)
	}
}

func validateFlowSchema(flowPath string) {

	workingDir, _ := os.Getwd()

	schemaURI := "file://" + workingDir + "/vendor/src/github.com/TIBCOSoftware/flogo-lib/schemas/flow_schema.json"
	flowURI := "file://" + workingDir + "/" + flowPath

	schemaLoader := gojsonschema.NewReferenceLoader(schemaURI)
	flowLoader := gojsonschema.NewReferenceLoader(flowURI)

	result, err := gojsonschema.Validate(schemaLoader, flowLoader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: cannot embed '%s', could not validate the flow's json\n  - %s\n", flowPath, err.Error())
		os.Exit(2)
	}

	if !result.Valid() {

		fmt.Fprintf(os.Stderr, "Error: cannot embed '%s', validation of the flow's json failed:\n", flowPath)
		for _, desc := range result.Errors() {
			fmt.Fprintf(os.Stderr, "  - %s\n", desc)
		}

		os.Exit(2)
	}
}

func getActivityTypes(flowObj interface{}, activityTypes map[string]bool) {

	switch obj := flowObj.(type) {
	case map[string]interface{}:
		for k, v := range obj {

			if k == "activityType" {
				activityType := v.(string)
				if len(activityType) != 0 {
					activityTypes[v.(string)] = true
				}
			} else {
				getActivityTypes(v, activityTypes)
			}
		}
	case []interface{}:
		for _, v := range obj {
			getActivityTypes(v, activityTypes)
		}
	}
}
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
	"github.com/TIBCOSoftware/flogo-cli/util"
	"net/url"
	"net/http"
	"path/filepath"
)

const flowSchemaFilePath string = "/vendor/src/github.com/TIBCOSoftware/flogo-lib/flow/flowdef/schema.json"

func ImportFlows(projectDescriptor *FlogoProjectDescriptor, flowDir string) map[string]string {

	flows := make(map[string]string)

	fileInfos, err := ioutil.ReadDir(flowDir)

	if err == nil {

		for _, fileInfo := range fileInfos {

			if !fileInfo.IsDir() {

				fileName := fileInfo.Name()
				flowFilePath := path(flowDir, fileName)

				// validate flow json
				ValidateFlow(projectDescriptor, flowFilePath, false)

				b64 := gzipAndB64(flowFilePath) //todo: is gzip necessary

				flows[genFlowURI(fileName)] = b64
			}
		}
	}

	return flows
}

func genFlowURI(fileName string) string {

	idx := strings.LastIndex(fileName, ".")

	if idx == -1 {
		return "embedded://" + fileName
	}

	return "embedded://" + fileName[:idx]
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

func ValidateFlow(projectDescriptor *FlogoProjectDescriptor, flowPath string, isURL bool) {

	// first validate the flow json
	validateFlowSchema(flowPath, isURL)

	//next check if all activities used int he flow are installed in engine

	var file []byte

	if isURL {

		flowURL, _ := url.Parse(flowPath)
		flowFilePath, local := fgutil.URLToFilePath(flowURL)

		if !local {
			resp, err := http.Get(flowURL.String())
			defer resp.Body.Close()

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: Unable to access '%s'\n  - %s\n", flowURL.String(), err.Error())
				os.Exit(2)
			}

			file, _ = ioutil.ReadAll(resp.Body)

		} else {
			file, _ = ioutil.ReadFile(flowFilePath)
		}

	} else {
		file, _ = ioutil.ReadFile(flowPath)
	}

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

func validateFlowSchema(flowPath string, isURL bool) {

	workingDir, _ := os.Getwd()

	//schemaURL,_ := fgutil.PathToFileURL(workingDir + flowSchemaFilePath)

	var flowURL string
	if isURL {
		flowURL = flowPath
	} else {
		var err error
		flowURL, err = fgutil.PathToFileURL(flowPath)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: cannot embed '%s', could not parse path\n", flowPath)
			os.Exit(2)
		}
	}

	//schemaLoader := gojsonschema.NewReferenceLoader(schemaURL)

	//work around for issue using file loader with gojsonschema
	schemaBytes, err := ioutil.ReadFile(workingDir + flowSchemaFilePath);
	schemaStr := string(schemaBytes)
	schemaLoader := gojsonschema.NewStringLoader(schemaStr)

	flowLoader := gojsonschema.NewReferenceLoader(flowURL)

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

func getAllActivityTypes(flowDir string) map[string]bool {

	fileInfos, err := ioutil.ReadDir(flowDir)

	activityTypes := make(map[string]bool)

	if err == nil {

		for _, fileInfo := range fileInfos {

			if !fileInfo.IsDir() {

				fileName := fileInfo.Name()
				flowFilePath := filepath.Join(flowDir, fileName)

				file, _ := ioutil.ReadFile(flowFilePath)

				var flowObj interface{}
				json.Unmarshal(file, &flowObj)
				getActivityTypes(flowObj, activityTypes)
			}
		}
	}

	return activityTypes
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

func getAllFlowExprs(flowDir string) map[string]map[int]string {

	fileInfos, err := ioutil.ReadDir(flowDir)

	allFlowExprs := make(map[string]map[int]string)

	if err == nil {

		for _, fileInfo := range fileInfos {

			if !fileInfo.IsDir() {

				fileName := fileInfo.Name()
				flowFilePath := filepath.Join(flowDir, fileName)

				file, _ := ioutil.ReadFile(flowFilePath)

				var flowObj interface{}
				json.Unmarshal(file, &flowObj)

				flowURI := genFlowURI(fileName)

				linkExprs := make(map[int]string)
				getLinkExprs(flowObj, linkExprs)

				if len(linkExprs) > 0 {
					allFlowExprs[flowURI] = linkExprs
				}
			}
		}
	}

	return allFlowExprs
}

func getLinkExprs(flowObj interface{}, linkExprs map[int]string) {

	switch obj := flowObj.(type) {
	case map[string]interface{}:
		for k, v := range obj {

			if k == "value" {
				value := v.(string)
				idFloat, hasID := obj["id"]
				_, hasFrom := obj["from"]

				if len(value) != 0 && hasID && hasFrom {
					id := int(idFloat.(float64))
					linkExprs[id] = value
				}
			} else {
				getLinkExprs(v, linkExprs)
			}
		}
	case []interface{}:
		for _, v := range obj {
			getLinkExprs(v, linkExprs)
		}
	}
}

func convertExprsToGo(linkExprs map[int]string) map[int]string {

	transExprs := make(map[int]string, len(linkExprs))

	for id, expr := range linkExprs {

		_, translated := transExpr(expr)
		transExprs[id] = translated
	}

	return transExprs;
}

func transExpr(s string) ([]string, string) {

	var attrs []string
	var rattrs []string

	strLen := len(s)

	for i := 0; i < strLen; i++ {
		if s[i] == '$' {
			var j int
			for j = i + 1; j < strLen; j++ {
				if s[j] == ' ' {
					break
				}
			}
			attrs = append(attrs, s[i+1:j])
			rattrs = append(rattrs, s[i:j])
			rattrs = append(rattrs, `v["`+s[i+1:j]+`"]`)
			i = j
		}
	}

	replacer := strings.NewReplacer(rattrs...)

	return attrs, replacer.Replace(s)
}

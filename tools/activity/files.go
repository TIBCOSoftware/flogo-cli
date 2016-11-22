package activity

import (
	"os"

	"github.com/TIBCOSoftware/flogo-cli/util"
)

const (
	fileDescriptor     string = "activity.json"
	fileActivityGo     string = "activity.go"
	fileActivityTestGo string = "activity_test.go"
	fileActivityMdGo   string = "activity_metadata.go"

	dirRT string = "runtime"
)

func createProjectDescriptor(sourcePath string, data interface{}) {

	filePath := fileDescriptor
	if len(sourcePath) > 0 {
		filePath = path(sourcePath, fileDescriptor)
	}

	f, _ := os.Create(filePath)
	fgutil.RenderTemplate(f, tplActivityDescriptorJSON, data)
	f.Close()
}

var tplActivityDescriptorJSON = `{
  "name": "{{.Name}}",
  "version": "0.0.1",
  "description": "activity description",
  "author": "Your Name <you.name@example.org>",
  "inputs":[
    {
      "name": "input",
      "type": "string"
    }
  ],
  "outputs": [
    {
      "name": "output",
      "type": "string"
    }
  ]
}`

func createActivityGoFile(codeSourcePath string, data interface{}) {
	f, _ := os.Create(path(codeSourcePath, fileActivityGo))
	fgutil.RenderTemplate(f, tplActivityGoFile, data)
	f.Close()
}

var tplActivityGoFile = `package {{.Name}}

import (
	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
)

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// init create & register activity
func init() {
	md := activity.NewMetadata(jsonMetadata)
	activity.Register(&MyActivity{metadata: md})
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, err error)  {

	// do eval

	return true, nil
}
`

func createActivityTestGoFile(codeSourcePath string, data interface{}) {
	f, _ := os.Create(path(codeSourcePath, fileActivityTestGo))
	fgutil.RenderTemplate(f, tplActivityTestGoFile, data)
	f.Close()
}

var tplActivityTestGoFile = `package {{.Name}}

import (
	"testing"
	"github.com/TIBCOSoftware/flogo-lib/flow/activity"
	"github.com/TIBCOSoftware/flogo-lib/flow/test"
)

func TestRegistered(t *testing.T) {
	act := activity.Get("{{.Name}}")

	if act == nil {
		t.Error("Activity Not Registered")
		t.Fail()
		return
	}
}

func TestEval(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	md := activity.NewMetadata(jsonMetadata)
	act := &MyActivity{metadata: md}

	tc := test.NewTestActivityContext(md)
	//setup attrs

	act.Eval(tc)

	//check result attr
}
`

func createMetadataGoFile(codeSourcePath string, data interface{}) {
	f, _ := os.Create(path(codeSourcePath, fileActivityMdGo))
	fgutil.RenderTemplate(f, tplMetadataGoFile, data)
	f.Close()
}

var tplMetadataGoFile = `package {{.Name}}

var jsonMetadata = ` + "`" + tplActivityDescriptorJSON + "`"

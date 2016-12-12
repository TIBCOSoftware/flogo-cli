package trigger

import (
	"os"

	"github.com/TIBCOSoftware/flogo-cli/util"
)

const (
	fileDescriptor    string = "trigger.json"
	fileTriggerGo     string = "trigger.go"
	fileTriggerTestGo string = "trigger_test.go"
	fileTriggerMdGo   string = "trigger_metadata.go"

	dirRT string = "runtime"
)

func createProjectDescriptor(sourcePath string, data interface{}) {

	filePath := fileDescriptor
	if len(sourcePath) > 0 {
		filePath = path(sourcePath, fileDescriptor)
	}

	f, _ := os.Create(filePath)
	fgutil.RenderTemplate(f, tplTriggerDescriptorJSON, data)
	f.Close()
}

var tplTriggerDescriptorJSON = `{
  "name": "{{.Name}}",
  "version": "0.0.1",
  "description": "trigger description",
  "author": "Your Name <you.name@example.org>",
  "settings":[
    {
      "name": "input",
      "type": "string",
      "value": "default"
    }
  ],
  "outputs": [
    {
      "name": "output",
      "type": "string"
    }
  ]
}`

func createTriggerGoFile(codeSourcePath string, data interface{}) {
	f, _ := os.Create(path(codeSourcePath, fileTriggerGo))
	fgutil.RenderTemplate(f, tplTriggerGoFile, data)
	f.Close()
}

var tplTriggerGoFile = `package {{.Name}}

import (
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

// MyTrigger is a stub for your Trigger implementation
type MyTrigger struct {
	metadata *trigger.Metadata
	runner   action.Runner
	config   *trigger.Config
}

func init() {
	md := trigger.NewMetadata(jsonMetadata)
	trigger.Register(&MyTrigger{metadata: md})
}

// Init implements trigger.Trigger.Init
func (t *MyTrigger) Init(config *trigger.Config, runner action.Runner) {
	t.config = config
	t.runner = runner
}

// Metadata implements trigger.Trigger.Metadata
func (t *MyTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Start implements trigger.Trigger.Start
func (t *MyTrigger) Start() error {
	// start the trigger
	return nil
}

// Stop implements trigger.Trigger.Start
func (t *MyTrigger) Stop() error {
	// stop the trigger
	return nil
}
`

func createTriggerTestGoFile(codeSourcePath string, data interface{}) {
	f, _ := os.Create(path(codeSourcePath, fileTriggerTestGo))
	fgutil.RenderTemplate(f, tplTriggerTestGoFile, data)
	f.Close()
}

var tplTriggerTestGoFile = `package {{.Name}}

import (
	"context"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
)

type TestRunner struct {
}

// Run implements action.Runner.Run
func (tr *TestRunner) Run(context context.Context, action action.Action, uri string, options interface{}) (code int, data interface{}, err error) {
	return 0, nil, nil
}

func TestRegistered(t *testing.T) {
	act := trigger.Get("{{.Name}}")

	if act == nil {
		t.Error("Trigger Not Registered")
		t.Fail()
		return
	}
}
`

func createMetadataGoFile(codeSourcePath string, data interface{}) {
	f, _ := os.Create(path(codeSourcePath, fileTriggerMdGo))
	fgutil.RenderTemplate(f, tplMetadataGoFile, data)
	f.Close()
}

var tplMetadataGoFile = `package {{.Name}}

var jsonMetadata = ` + "`" + tplTriggerDescriptorJSON + "`"

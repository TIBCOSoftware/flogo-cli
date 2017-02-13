package gen

import (
	"github.com/TIBCOSoftware/flogo-cli/util"
)

const (
	fileActionDescriptor string = "action.json"
	fileActionGo         string = "action.go"
	fileActionGoTest     string = "action_test.go"
)

type ActionGenerator struct {

}

func (g *ActionGenerator) Description() string {
	return "generates an action project"
}

func (g *ActionGenerator) Generate(basePath string, data interface{}) error {

	err := fgutil.CreateFileFromTemplate(basePath, fileActionDescriptor, tplActionDescriptor, data)
	if err != nil {
		return err
	}
	err = fgutil.CreateFileFromTemplate(basePath, fileActionGo, tplActionGo, data)
	if err != nil {
		return err
	}

	err = fgutil.CreateFileFromTemplate(basePath, fileActionGoTest, tplActionGoTest, data)
	if err != nil {
		return err
	}

	return nil
}

var tplActionDescriptor = `{
  "name": "{{.Name}}",
  "version": "0.0.1",
  "type": "flogo:action",
  "description": "action description",
  "author": "Your Name <you.name@example.org>",
}`


var tplActionGo = `package {{.Name}}

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
var tplActionGoTest = `package {{.Name}}

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

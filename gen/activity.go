package gen

import (
	"github.com/TIBCOSoftware/flogo-cli/util"
)

const (
	fileActivityDescriptor string = "activity.json"
	fileActivityGo         string = "activity.go"
	fileActivityGoTest     string = "activity_test.go"
)

type ActivityGenerator struct {

}

func (g *ActivityGenerator) Description() string {
	return "generates an activity project"
}

func (g *ActivityGenerator) Generate(basePath string, data interface{}) error {

	err := fgutil.CreateFileFromTemplate(basePath, fileActivityDescriptor, tplActivityDescriptor, data)
	if err != nil {
		return err
	}

	err = fgutil.CreateFileFromTemplate(basePath, fileActivityGo, tplActivityGo, data)
	if err != nil {
		return err
	}

	err = fgutil.CreateFileFromTemplate(basePath, fileActivityGoTest, tplActivityGoTest, data)
	if err != nil {
		return err
	}

	return nil
}

var tplActivityDescriptor = `{
  "name": "{{.Name}}",
  "version": "0.0.1",
  "type": "flogo:activity",
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

var tplActivityGo = `package {{.Name}}

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

var tplActivityGoTest = `package {{.Name}}

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
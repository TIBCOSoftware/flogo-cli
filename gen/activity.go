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
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MyActivity{metadata: metadata}
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
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
)

var activityMetadata *activity.Metadata

func getActivityMetadata() *activity.Metadata {

	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("activity.json")
		if err != nil{
			panic("No Json Metadata found for activity.json path")
		}

		activityMetadata = activity.NewMetadata(string(jsonMetadataBytes))
	}

	return activityMetadata
}

func TestCreate(t *testing.T) {

	act := NewActivity(getActivityMetadata())

	if act == nil {
		t.Error("Activity Not Created")
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

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs

	act.Eval(tc)

	//check result attr
}
`
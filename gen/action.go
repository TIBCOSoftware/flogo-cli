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
	"context"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
)

// MyAction is a stub for your Activity implementation
type MyAction struct {

}

// init create & register action
func init() {
	action.Register("{{.Name}}", &MyAction{})) {
}

// Eval implements action.Action.Run
func (a *MyAction) Run(context context.Context, uri string, options interface{}, handler ResultHandler) error  {

	// perform action
	return nil
}
`
var tplActionGoTest = `package {{.Name}}

import (
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
)

func TestRegistered(t *testing.T) {
	act := action.Get("{{.Name}}")

	if act == nil {
		t.Error("Action Not Registered")
		t.Fail()
		return
	}
}

func TestRun(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := &MyAction{}

	//setup context

	act.Run(ctx, "", nil, nil)
}
`

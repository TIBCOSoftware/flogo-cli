package app

import (
	"testing"
	"flag"
	"github.com/TIBCOSoftware/flogo-cli/cli"
	"os"
	"path"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"bytes"
)

var GOLD_FLOGO_JSON = `{
  "name": "testApp",
  "type": "flogo:app",
  "version": "0.0.1",
  "description": "My flogo application description",
  "triggers": [
    {
      "id": "my_rest_trigger",
      "ref": "github.com/TIBCOSoftware/flogo-contrib/trigger/rest",
      "settings": {
        "port": "9233"
      },
      "handlers": [
        {
          "actionId": "my_simple_flow",
          "settings": {
            "method": "GET",
            "path": "/test"
          }
        }
      ]
    }
  ],
  "actions": [
    {
      "id": "my_simple_flow",
      "name": "my simple flow",
      "ref": "github.com/TIBCOSoftware/flogo-contrib/action/flow",
      "data": {
        "flow": {
          "name": "my simple flow",
          "attributes": [],
          "rootTask": {
            "id": 1,
            "type": 1,
            "tasks": [
              {
                "id": 2,
                "type": 1,
                "activityRef": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
                "name": "log",
                "attributes": [
                  {
                    "name": "message",
                    "value": "Simple Log",
                    "type": "string"
                  }
                ]
              }
            ],
            "links": [
            ]
          }
        }
      }
    }
  ]
}`

type TestEnv struct {
	currentDir string
}

func (t *TestEnv) getTestwd() (dir string, err error){
	return t.currentDir, nil
}

func (t *TestEnv) cleanup(){
	os.RemoveAll(t.currentDir)
}


// TestCmdCreate_Exec test the default cmd create, create new app
func TestCmdCreate_Exec (t *testing.T) {
	// TODO remote this after merging
	err := os.Setenv("FLOGO_BUILD_EXPERIMENTAL", "true")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Unsetenv("FLOGO_BUILD_EXPERIMENTAL")
	tempDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	testEnv := &TestEnv{currentDir:tempDir}
	defer testEnv.cleanup()
	t.Logf("Current dir '%s'", testEnv.currentDir)
	cmd := &cmdCreate{option: optCreate, currentDir: testEnv.getTestwd}
	fs := flag.NewFlagSet("test", flag.ExitOnError)
	appName := "testApp"
	args := []string{appName}

	err = cli.ExecCommand(fs, cmd, args)
	if err != nil {
		t.Fatal(err)
	}

	// Validate the structure
	currentDir := testEnv.currentDir
	_, err = os.Stat(currentDir)
	assert.Nil(t, err, "There should be a folder for temp dir '%s'", currentDir)

	// There should be a folder with app name
	appFolderPath := path.Join(currentDir, appName)
	fi, err := os.Stat(appFolderPath)
	assert.Nil(t, err, "There should be a folder for app name '%s'", appFolderPath)
	assert.True(t, fi.IsDir(), "There should be a folder for app name '%s'", appFolderPath)

	// There should be a flogo.json file
	flogojsonPath := path.Join(appFolderPath, fileDescriptor)
	fi, err = os.Stat(flogojsonPath)
	assert.Nil(t, err, "There should be a file for flogo json '%s'", flogojsonPath)

	// Compare the flogo.json with GOLD file
	jsonBytes, err := ioutil.ReadFile(flogojsonPath)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, bytes.Equal(jsonBytes, []byte(GOLD_FLOGO_JSON)), "Unexpected flogo.json format")

	// There should be a src folder
	srcPath := path.Join(appFolderPath, "src")
	fi, err = os.Stat(srcPath)
	assert.Nil(t, err, "There should be a folder for src '%s'", srcPath)

	// There should be an app folder
	appSrcFolderPath := path.Join(srcPath, appName)
	fi, err = os.Stat(appSrcFolderPath)
	assert.Nil(t, err, "There should be a folder for app name '%s'", appSrcFolderPath)
	assert.True(t, fi.IsDir(), "There should be a folder for app name '%s'", appSrcFolderPath)

	// There should be a main file
	mainPath := path.Join(appSrcFolderPath, fileMainGo)
	fi, err = os.Stat(mainPath)
	assert.Nil(t, err, "There should be a file for main.go '%s'", mainPath)

	// There should be an imports file
	importsPath := path.Join(appSrcFolderPath, fileImportsGo)
	fi, err = os.Stat(importsPath)
	assert.Nil(t, err, "There should be a file for imports.go '%s'", importsPath)






}

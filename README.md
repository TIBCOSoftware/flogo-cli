# flogo-cli
> Command line tool for building **Flogo**-based applications.

**Flogo** is an IoT Integration framework written in Go. It was designed from the ground up to be robust enough for cloud applications and at the same time sufficiently lean for IoT devices.


## Installation
### Prerequisites
* The Go programming language 1.7 or later should be [installed](https://golang.org/doc/install).
* Set GOPATH environment variable on your system
* In order to simplify development and building in Go, we are using the **gb** build tool.  It can be downloaded from [here](https://getgb.io).  

### Install flogo
    go get github.com/TIBCOSoftware/flogo-cli/...

### Update flogo
    go get -u github.com/TIBCOSoftware/flogo-cli/...
    
### IMPORTANT UPDATE ##

The original **flogo** CLI tool has been deprecated and will be going away in the near future.  It has been temporarily renamed to **flogo_old** and its documentation can still be accessed [here](README_OLD.md).
    
## Getting Started
A flogo application is created using the **flogo** CLI tool.  The tool can be used to create an application from an existing *flogo.json* or to create a simple base application to get you started.  In this example we will walk you through creating the base/sample application.

To create the base application, which consists of a REST trigger and a simple flow with a log activity, you use the following commands.


```bash
flogo create myApp
cd myApp

flogo build
```

- Cd bin folder 
- Start flogo engine by running ./myApp
- Flogo will start a REST server
- Send GET request to run the flow. eg: http://localhost: 9233/test

The built in sample application is based of the following flogo.json.  This file can be manually modified to add additional triggers and flow actions.  This file can also be generated using the flogo-web UI.

```json
{
  "name": "myApp",
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
      "ref": "github.com/TIBCOSoftware/flogo-contrib/action/flow",
      "data": {
        "flow": {
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
}
```


For more details about the REST Trigger configuration go [here](https://github.com/TIBCOSoftware/flogo-contrib/tree/master/trigger/rest#example-configurations)

## Documentation
Additional documentation on flogo and the CLI tool

  - **flogo tool**
    - creating an [application](docs/app.md)
  - **flogogen tool**
    - creating a [trigger](docs/trigger.md)
    - creating a [action](docs/action.md)
    - creating an [activity](docs/activity.md)
    - creating a [flow model](docs/flow_model.md)

## Contributing and support

### Contributing

We welcome all bug fixes and issue reports.

Pull requests are also welcome. If you would like to submit one, please follow these guidelines:

* Code must be [gofmt](https://golang.org/cmd/gofmt/) compliant.
* Execute [golint](https://github.com/golang/lint) on your code.
* Document all funcs, structs and types.
* Ensure that 'go test' succeeds.


Please submit a github issue if you would like to propose a significant change or request a new feature.

### Build flogo from source
```
$go get github.com/TIBCOSoftware/flogo-cli/...

$cd $GOPATH/src/github.com/TIBCOSoftware/flogo-cli

[optional, only if building from branch] 
$git checkout my_branch

[need to manually go get all dependencies for example:] 
$go get github.com/xeipuuv/gojsonschema

$go install ./... 
```

##License
flogo-cli is licensed under a BSD-type license. See TIBCO LICENSE.txt for license text.


### Support
You can post your questions via [GitHub issues](https://github.com/TIBCOSoftware/flogo/issues)

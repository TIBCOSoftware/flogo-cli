# flogo_old
> Deprecated Command line tool for building **Flogo**-based applications.

This version of the CLI tool has been deprecated, please refer to the latest [README](README.md) on how to use the current CLI tool.  Note that the deprecated version of the CLI tool was renamed from **flogo** to **flogo_old**.

## Installation
### Prerequisites
* The Go programming language 1.7 or later should be [installed](https://golang.org/doc/install).
* Set GOPATH environment variable on your system
* In order to simplify development and building in Go, we are using the **gb** build tool.  It can be downloaded from [here](https://getgb.io).  

### Install flogo CLI
    go get github.com/TIBCOSoftware/flogo-cli/...

### Update flogo CLI
    go get -u github.com/TIBCOSoftware/flogo-cli/...
    
## Getting Started
This simple example demonstrates how to create a simple flogo application that has a log activity and REST trigger.


- Download flow [myflow.json](https://github.com/TIBCOSoftware/flogo-cli/blob/master/samples/gettingstarted/cli/myflow.json) to build your first application. You can also download more samples from the [samples folder](https://github.com/TIBCOSoftware/flogo/tree/master/samples) in the flogo repo. 

```bash
flogo_old create myApp
cd myApp

flogo_old add activity github.com/TIBCOSoftware/flogo-contrib/activity/log
flogo_old add trigger github.com/TIBCOSoftware/flogo-contrib/trigger/rest
#Make sure myflow.json file under current location
flogo_old add flow myflow.json
flogo_old build
```

- Cd bin folder and open trigger.json in a text editor
- Replace content of trigger.json with the following

```json
{
  "triggers": [
    {
      "name": "tibco-rest",
      "settings": {
        "port": "9999"
      },
      "endpoints": [
        {
          "actionType": "flow",
          "actionURI": "embedded://myflow",
          "settings": {
            "autoIdReply": "true",
            "method": "GET",
            "path": "/flow",
            "useReplyHandler": "true"
          }
        }
      ]
    }
  ]
}
```

- Start flogo engine by running ./myApp
- Flogo will start a REST server
- Send GET request to run the flow. eg: http://localhost:9999/flow

For more details about the REST Trigger configuration go [here](https://github.com/TIBCOSoftware/flogo-contrib/tree/master/trigger/rest#example-configurations)

## Documentation
Additional documentation on flogo_old CLI tool

  - **Flogo tool**
    - creating an [application](docs/app_old.md)
    - creating an [activity](docs/activity_old.md)
    - creating a [trigger](docs/trigger_old.md)
    - creating a [model](docs/model_old.md)

## Contributing and support

### Contributing

We welcome all bug fixes and issue reports.

Pull requests are also welcome. If you would like to submit one, please follow these guidelines:

* Code must be [gofmt](https://golang.org/cmd/gofmt/) compliant.
* Execute [golint](https://github.com/golang/lint) on your code.
* Document all funcs, structs and types.
* Ensure that 'go test' succeeds.


Please submit a github issue if you would like to propose a significant change or request a new feature.

### Build flogo CLI from source
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

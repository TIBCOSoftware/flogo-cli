# flogo-cli
> Command line tool for building **Flogo**-based applications.

**Flogo** is an IoT Integration framework written in Go. It was designed from the ground up to be robust enough for cloud applications and at the same time sufficiently lean for IoT devices.


## Installation
### Prerequisites
* The Go programming language should be [installed](https://golang.org/doc/install).
* Set GOPATH environment variable on your system
* In order to simplify development and building in Go, we are using the **gb** build tool.  It can be downloaded from [here](https://getgb.io).  

### Install flogo
    go get github.com/TIBCOSoftware/flogo-cli/...

### Update flogo
    go get -u github.com/TIBCOSoftware/flogo-cli/...
    
## Getting Started
This simple example demonstrates how to create a simple flogo application that has a log activity and REST trigger.


- Download flow [myflow.json](https://github.com/TIBCOSoftware/flogo-cli/blob/master/samples/gettingstarted/cli/myflow.json) to build your first application. You can also download more samples from the [samples folder](https://github.com/TIBCOSoftware/flogo/tree/master/samples) in the flogo repo. 

```bash
flogo create myApp
cd myApp

flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/log
flogo add trigger github.com/TIBCOSoftware/flogo-contrib/trigger/rest
#Make sure myflow.json file under current location
flogo add flow myflow.json
flogo build
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
Additional documentation on flogo and the CLI tool

  - **Flogo tool**
    - creating an [application](docs/app.md)
    - creating an [activity](docs/activity.md)
    - creating a [trigger](docs/trigger.md)
    - creating a [model](docs/model.md)

## Contributing and support

### Contributing

We welcome all bug fixes and issue reports.

Pull requests are also welcome. If you would like to submit one, please follow these guidelines:

* Code must be [gofmt](https://golang.org/cmd/gofmt/) compliant.
* Execute [golint](https://github.com/golang/lint) on your code.
* Document all funcs, structs and types.
* Ensure that 'go test' succeeds.


Please submit a github issue if you would like to propose a significant change or request a new feature.

##License
flogo-cli is licensed under a BSD-type license. See TIBCO LICENSE.txt for license text.


### Support
You can post your questions via [GitHub issues](https://github.com/TIBCOSoftware/flogo/issues)

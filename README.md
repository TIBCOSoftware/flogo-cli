# flogo
> Command line tool for building **Flogo**-based applications.

![Flogo icon](floyd.png)

**Flogo** is micro-flow engine written in Go. It was designed from the ground up to be robust enough for cloud applications and at the same time sufficiently lean for IOT devices.


## Installation
### Prerequisites
* The Go programming language should be [installed](https://golang.org/doc/install).
* In order to simplify development and building in Go, we are using the **gb** build tool.  It can be downloaded from [here](https://getgb.io).  

### Install flogo
    go get github.com/TIBCOSoftware/flogo/...

## Creating a new Flogo project
This simple example demonstrates how to create a simple flogo application that has a log activity and REST trigger.

```bash
flogo create myApp
cd myApp

flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/log
flogo add trigger github.com/TIBCOSoftware/flogo-contrib/trigger/rest
flogo build
```

## Contributing and support

### Contributing

We welcome all bug fixes and issue reports.

Pull requests are also welcome. If you would like to submit one, please follow these guidelines:

* Code must be [gofmt](https://golang.org/cmd/gofmt/) compliant.
* Execute [golint](https://github.com/golang/lint) on your code.
* Document all funcs, structs and types.
* Ensure that 'go test' succeeds.


Please submit a github issue if you would like to propose a significant change or request a new feature.

### Support
For Q&A you can post your questions on [Slack](https://tibco-cloud.slack.com/messages/flogo-general/)

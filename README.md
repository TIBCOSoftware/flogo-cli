<p align="center">
  <img src ="https://raw.githubusercontent.com/TIBCOSoftware/flogo/master/images/projectflogo.png" />
</p>

<p align="center" >
  <b>Serverless functions and edge microservices made painless</b>
</p>

<p align="center">
  <img src="https://travis-ci.org/TIBCOSoftware/flogo-cli.svg"/>
  <img src="https://img.shields.io/badge/dependencies-up%20to%20date-green.svg"/>
  <img src="https://img.shields.io/badge/license-BSD%20style-blue.svg"/>
  <a href="https://gitter.im/project-flogo/Lobby?utm_source=share-link&utm_medium=link&utm_campaign=share-link"><img src="https://badges.gitter.im/Join%20Chat.svg"/></a>
</p>

<p align="center">
  <a href="#Installation">Installation</a> | <a href="#getting-started">Getting Started</a> | <a href="#repos">Repos</a> | <a href="#contributing">Contributing</a> | <a href="#license">License</a>
</p>

<br/>
Project Flogo is an open source framework to simplify building efficient & modern serverless functions and edge microservices and _this_ is the cli that makes it all happen. 

## Installation
### Prerequisites
To get started with the Project Flogo cli you'll need to have a few things
* The Go programming language version 1.8 or later should be [installed](https://golang.org/doc/install).
* The **GOPATH** environment variable on your system must be set properly
* In order to simplify dependency management, we're using **go dep**. You can install that by following the instructions [here](https://github.com/golang/dep#setup).

### Install the cli
To install the cli, simply open a terminal and enter the below command
```
$ go get -u github.com/TIBCOSoftware/flogo-cli/...
```
_Note that the -u parameter automatically updates the cli if it exists_

### Build the cli from source
You can build the cli from source code as well, which is convenient if you're developing new features for it! To do that, follow these easy steps
```bash
# Get the flogo-cli from GitHub
$ go get github.com/TIBCOSoftware/flogo-cli/...

# Go to the right directory
$ cd $GOPATH/src/github.com/TIBCOSoftware/flogo-cli

# Optionally check out the branch you want to use 
$ git checkout my_branch

# Run the install command
$ go install ./... 
```

## Getting started
Getting started should be easy and fun, and so is getting started with the Flogo cli. 

First, create a file called `flogo.json` and with the below content (which is a simple app with an [HTTP trigger](https://tibcosoftware.github.io/flogo/development/webui/triggers/rest/))
```json
{
  "name": "SampleApp",
  "type": "flogo:app",
  "version": "0.0.1",
  "appModel": "1.0.0",
  "triggers": [
    {
      "id": "receive_http_message",
      "ref": "github.com/TIBCOSoftware/flogo-contrib/trigger/rest",
      "name": "Receive HTTP Message",
      "description": "Simple REST Trigger",
      "settings": {
        "port": 9233
      },
      "handlers": [
        {
          "action": {
            "ref": "github.com/TIBCOSoftware/flogo-contrib/action/flow",
            "data": {
              "flowURI": "res://flow:sample_flow"
            }
          },
          "settings": {
            "method": "GET",
            "path": "/test"
          }
        }
      ]
    }
  ],
  "resources": [
    {
      "id": "flow:sample_flow",
      "data": {
        "name": "SampleFlow",
        "tasks": [
          {
            "id": "log_2",
            "name": "Log Message",
            "description": "Simple Log Activity",
            "activity": {
              "ref": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
              "input": {
                "message": "Simple Log",
                "flowInfo": "false",
                "addToFlow": "false"
              }
            }
          }
        ]
      }
    }
  ]
}
```

Based on this file we'll create a new flogo app
```bash
$ flogo create -f flogo.json myApp
```

From the app folder we can build the executable
```bash
$ cd myApp
$ flogo build -e
```

Now that there is an executable we can run it!
```bash
$ cd bin
$ ./myApp
```

The above commands will start the REST server and wait for messages to be sent to `http://localhost:9233/test`. To send a message you can use your browser, or a new terminal window and run
```bash
$ curl http://localhost:9233/test
```

_For more tutorials check out the [Labs](https://tibcosoftware.github.io/flogo/labs/) section in our documentation_

## Repos

[Project Flogo](https://github.com/TIBCOSoftware/flogo) consists of the following sub-projects available as separate repos:
* [flogo-cli](https://github.com/TIBCOSoftware/flogo-cli): Command line tools for building Flogo apps & extensions (you're here now)
* [flogo-lib](https://github.com/TIBCOSoftware/flogo-lib): The core Flogo library
* [flogo-services](https://github.com/TIBCOSoftware/flogo-services): Backing services required by Flogo 
* [flogo-contrib](https://github.com/TIBCOSoftware/flogo-contrib): Flogo contributions/extensions

## Contributing
Want to contribute to Project Flogo? We've made it easy, all you need to do is fork the repository you intend to contribute to, make your changes and create a Pull Request! Once the pull request has been created, you'll be prompted to sign the CLA (Contributor License Agreement) online.

Not sure where to start? No problem, here are a few suggestions:

* [flogo-contrib](https://github.com/TIBCOSoftware/flogo-contrib): This repository contains all of the contributions, such as activities, triggers, etc. Perhaps there is something missing? Create a new activity or trigger or fix a bug in an existing activity or trigger.
* Browse all of the Project Flogo repositories and look for issues tagged `kind/help-wanted` or `good first issue`

If you have any questions, feel free to post an issue and tag it as a question, email flogo-oss@tibco.com or chat with the team and community:

* The [project-flogo/Lobby](https://gitter.im/project-flogo/Lobby?utm_source=share-link&utm_medium=link&utm_campaign=share-link) Gitter channel should be used for general discussions, start here for all things Flogo!
* The [project-flogo/developers](https://gitter.im/project-flogo/developers?utm_source=share-link&utm_medium=link&utm_campaign=share-link) Gitter channel should be used for developer/contributor focused conversations. 

For additional details, refer to the [Contribution Guidelines](https://github.com/TIBCOSoftware/flogo/blob/master/CONTRIBUTING.md).

## License 
Flogo source code in [this](https://github.com/TIBCOSoftware/flogo-cli) repository is under a BSD-style license, refer to [LICENSE](https://github.com/TIBCOSoftware/flogo-cli/blob/master/LICENSE) 

# application - DEPRECATED
> Details on flogo application projects and associated CLI commands.

## Commands
#### create
This command creates a flogo application project.
	
	flogo_old create my_app
	
### add
This command is used to add a activity, trigger, flow or model to the application.

*activity*

	flogo_old add activity github.com/TIBCOSoftware/flogo-contrib/activity/log
	
*trigger*

	flogo_old add trigger github.com/TIBCOSoftware/flogo-contrib/trigger/rest
	
*model*

	flogo_old add model github.com/TIBCOSoftware/flogo-contrib/model/simple	  
	
*flow*

	flogo_old add flow file:///tmp/myflow.json
		
Note: tibco-simple model is added to an application by default 	
### del
This command is used to remove a activity, trigger or model from the application.

*activity*

	flogo_old del activity tibco-log
	
*trigger*

	flogo_old del trigger tibco-rest
	
*model*

	flogo_old del model tibco-simple	
	
### list
This command is used to list the activities, triggers, flows and models added to the application.  
	 
	flogo_old list
	
	Activities:
    	- tibco-log [github.com/TIBCOSoftware/flogo-contrib/activity/log]

	Triggers:
   		- tibco-rest [github.com/TIBCOSoftware/flogo-contrib/trigger/rest]

	Models:
   		- tibco-simple [github.com/TIBCOSoftware/flogo-contrib/model/simple]

	Flows:
		- myflow

### build
This command is used to build the application.

 	flogo_old build
 	
**options**
	
- [ -o ] : optimize compilation, application will only contain activities and triggers used by its flows
- [ -i ] : incorporates the configuration into the compiled application	 	 
- [ -c configDir] : specifies the directory to use for configuration when using the -i flag

### help
This command is used to display help on a particular command
	
	flogo_old help build 

##Application Project

###Structure

The create command creates a basic structure and files for an application.


	my_app/
		bin/
			config.json
			triggers.json
		flogo.json
		flows/
		src/
			my_app/
				config.go
				env.go
				flows.go
				imports.go
				main.go
		vendor/
		
**files**

- *flogo.json* : flogo project metadata json file
- *config.json* : configuration for the application
- *triggers.json* : trigger configuration for the application
- *config.go* : contains embedded configuration or reference to config.json
- *env.go* : basic engine environment configuration
- *flows.go* : contains embedded flows, gzipped and stored in base64
- *imports.go* : contains go imports for activities, triggers and models used by the application
- *main.go* : basic/initial test file for the model

**directories**	
	
- *bin* :	contains the application binary and configuration
- *flows* : contains the flows to embed
- *vendor* : go libraries

###Metadata

The *flogo.json* file is the metadata describing the application project.  It includes the name, version and the components (activities, triggers and models) that have been installed.

	{
	  "name": "my_app",
	  "version": "0.0.1",
	  "description": "My flogo application description",
	  "models": [
	    {
	      "name": "tibco-simple",
	      "path": "github.com/TIBCOSoftware/flogo-contrib/model/simple",
	      "version": "latest"
	    }
	  ],
	  "activities": [
	    {
	      "name": "tibco-log",
	      "path": "github.com/TIBCOSoftware/flogo-contrib/activity/log",
	      "version": "latest"
	    }
	  ],
	  "triggers": [
	    {
	      "name": "tibco-rest",
	      "path": "github.com/TIBCOSoftware/flogo-contrib/trigger/rest",
	      "version": "latest"
	    }
	  ]
	}

## Application Configuration

### Application

The *config.json* file contains the configration for application.  It is used to configure the internal process engine and services.

	{
	  "loglevel": "INFO",
	  "actionRunner": {
	    "type": "pooled"
	  },
	  "services": [
	    {
	      "name": "stateRecorder",
	      "enabled": false,
	      "settings": {
	        "host": "",
	        "port": ""
	      }
	    },
	    {
	      "name": "engineTester",
	      "enabled": true,
	      "settings": {
	        "port": "8080"
	      }
	    }
	  ]
	}

***Settings***

- loglevel: set the loglevel for the application

- *actionRunner* runs the action
	- pooled: uses a worker pool to execute actions
		- numworkers: the number of action runners
		- workQueueSize: the max number of queued actionss to execute
	- direct: actions are executed on the same goroutine/thread of the trigger

***Services***

- *stateRecorder* recordes the full/incremental state of a flow
	- enabled: true/false
	- host: the host of the stateRecorder service
	- port: the port of the stateRecorder service
- *engineTester* is a simple REST interface used to directly start a flow and bypass a trigger, used the the UI to directly execute a flow
	- enabled: true/false
	- port: the local port to expost the engineTester service

### Triggers
The *triggers.json* contains the configuration for the triggers used by the application.

	{
      "triggers": [
        {
          "name": "tibco-rest",
          "settings": {
            "port": "9090"
          },
          "endpoints": [
            {
              "actionType": "flow",
              "actionURI": "embedded://myflow",
              "settings": {
                "autoIdReply": "true",
                "method": "POST",
                "path": "/device/update"
              }
            }
          ]
        }
      ]
    }

***Trigger Configuration***

- name: the name of the trigger
- settings: global settings for the trigger
- *endpoints* the endpoints configured for the trigger
	- actionType: the type of action the endpoint runs
	- actionURI: the uri for the action
	- settings: the endpoint specific settings

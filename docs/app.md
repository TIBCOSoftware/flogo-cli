# application
> Details on flogo application related commands.

## Commands
#### create
This command creates a flogo application project.
	
	flogo create my_app
	
### add
This command is used to add a activity, trigger, flow or model to the application.

*activity*

	flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/log
	
*trigger*

	flogo add trigger github.com/TIBCOSoftware/flogo-contrib/trigger/rest
	
*model*

	flogo add model github.com/TIBCOSoftware/flogo-contrib/model/simple	  
	
*flow*

	flogo add flow file:///tmp/myflow.json
		
Note: tibco-simple model is added to an application by default 	
### del
This command is used to remove a activity, trigger or model from the application.

*activity*

	flogo del activity tibco-log
	
*trigger*

	flogo del trigger tibco-rest
	
*model*

	flogo del model tibco-simple	
	
### list
This command is used to list the activities, triggers, flows and models added to the application.  
	 
	flogo list
	
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

 	flogo build
 	
**options**
	
- [ -o ] : optimize compilation, application will only contain activities and triggers used by its flows
- [ -i ] : incorporates the configuration into the compiled application	 	 

### help
This command is used to display help on a particular command
	
	flogo help build 

##Project Structure

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
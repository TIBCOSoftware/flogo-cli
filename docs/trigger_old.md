# trigger - DEPRECATED
> Details on flogo trigger related commands.  Used to create a custom trigger.

## Commands
#### create
This command creates a flogo trigger project.
	
	flogo_old trigger create mytrigger
	 	 

### help
This command is used to display help on a particular command
	
	flogo_old trigger help create

##Project Structure

The create command creates a basic structure and files for a trigger.


	mytrigger\
		src\
			trigger.json
			runtime\
				trigger.go
				trigger_metadata.go
				trigger_test.go
		vendor\

**files**

- *trigger.json* : trigger project metadata json file
- *trigger.go*   : rudimentary trigger implementation in go
- *trigger_metadata.go* : trigger metadata go file
- *trigger_test.go* : basic/initial test file for the trigger

		
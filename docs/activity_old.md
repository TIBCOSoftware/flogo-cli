# activity - DEPRECATED
> Details on flogo activity related commands.  Used to create a custom activity.

## Commands
#### create
This command creates a flogo activity project.
	
	flogo_old activity create myactivity
	 	 

### help
This command is used to display help on a particular command
	
	flogo_old activity help create

##Project Structure

The create command creates a basic structure and files for an activity.


	myactivity\
		src\
			activity.json
			runtime\
				activity.go
				activity_metadata.go
				activity_test.go
		vendor\

**files**

- *activity.json* : activity project metadata json file
- *activity.go*   : rudimentary activity implementation in go
- *activity_metadata.go* : activity metadata go file
- *activity_test.go* : basic/initial test file for the activity

		
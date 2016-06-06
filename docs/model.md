# model
> Details on flogo model related commands.  Used to create a custom model.

## Commands
#### create
This command creates a flogo model project.
	
	flogo model create mymodel
	 	 

### help
This command is used to display help on a particular command
	
	flogo model help create

##Project Structure

The create command creates a basic structure and files for an activity.


	mymodel/
		src/
			mymodel/
				model.go
				model.json
				model_test.go
		vendor/

**files**

- *model.json* : model project metadata json file
- *model.go*   : rudimentary model implementation in go
- *model_test.go* : basic/initial test file for the model

**directories**	
	
- *vendor*: go libraries
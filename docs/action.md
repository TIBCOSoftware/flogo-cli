# action
> Details on how to create a custom flogo action.  In order to facilitate this you can use the **flogogen** tool.

## flogogen Command

This command generates scafolding for a flogo action project.
	
	flogogen action myaction
	 	 

##Project Structure

The create command creates a basic structure and files for an action.


	myaction\
			action.json
			action.go
			action_test.go

**files**

- *action.json* : action project metadata json file
- *action.go*   : rudimentary action implementation in go
- *action_test.go* : basic/initial test file for the action

		
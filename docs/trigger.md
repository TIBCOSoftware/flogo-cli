# trigger
> Details on how to create a custom flogo trigger.  In order to facilitate this you can use the **flogogen** tool.

## flogogen Command

This command generates scafolding for a flogo trigger project.
	
	flogogen trigger mytrigger
	 	 

##Project Structure

The create command creates a basic structure and files for an trigger.


	mytrigger\
			trigger.json
			trigger.go
			trigger_test.go

**files**

- *trigger.json* : trigger project metadata json file
- *trigger.go*   : rudimentary trigger implementation in go
- *trigger_test.go* : basic/initial test file for the trigger

		
# activity
> Details on how to create a custom flogo activity.  In order to facilitate this you can use the **flogogen** tool.

## flogogen Command

This command generates scafolding for a flogo activity project.
	
	flogogen activity myactivity
	 	 

##Project Structure

The create command creates a basic structure and files for an activity.


	myactivity\
			activity.json
			activity.go
			activity_test.go

**files**

- *activity.json* : activity project metadata json file
- *activity.go*   : rudimentary activity implementation in go
- *activity_test.go* : basic/initial test file for the activity

		
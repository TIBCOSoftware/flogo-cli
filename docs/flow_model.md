# flowmodel
> Details on how to create a custom flogo flowmodel.  In order to facilitate this you can use the **flogogen** tool.

## flogogen Command

This command generates scafolding for a flogo flowmodel project.
	
	flogogen flowmodel myflowmodel
	 	 

##Project Structure

The create command creates a basic structure and files for an flowmodel.


	myflowmodel\
			flowmodel.json
			flowmodel.go
			flowmodel_test.go

**files**

- *flowmodel.json* : flowmodel project metadata json file
- *flowmodel.go*   : rudimentary flowmodel implementation in go
- *flowmodel_test.go* : basic/initial test file for the flowmodel

		
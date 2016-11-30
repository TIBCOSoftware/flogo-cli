# device
> Details on how to use flogo device trigger in your flogo application projects and associated CLI commands.

## Installation
### Prerequisites

The flogo cli tool uses [platformio](https://http://platformio.org/) to build and deploy to IoT devices.  In order to use the device commands you most have the [platformio CLI](http://platformio.org/get-started/cli) tool installed.  

Follow the [installation](http://docs.platformio.org/en/stable/installation.html) instructions to install the platformio CLI tool in your enviroment.

If you have a Mac and homebrew installed you can just do the following

    brew install platformio


## Commands
#### prepare
This command creates a generates the source files for the IoT device associated with your device-trigger.  It places the source under the "devices" directory
	
	flogo device prepare

#### build
This command builds the firmware for the device from the source files
	
	flogo device build
	
#### upload
This command uploads the firmware to the device.  The command must be run in the corresponding trigger directory under devices.
	
	flogo device upload
	

## Application with Devices

### Device Trigger

A **device trigger** must be added to your project in order to use the device commands.  There is one available in our flogo-contrib repository.  This is a preview release, so the trigger must be installed from the '*device-trigger*' branch.

```flogo add -b device-trigger trigger github.com/TIBCOSoftware/flogo-contrib/trigger/device```

### Configuration
 
The device trigger must be configured in the *triggers.json* prior to using the device commands. The following is an example configuration for an adafruit feather m0 that kicks off the the "test" flow when pin A3 is turned on.

```json
{
  "triggers": [
    {
      "name": "tibco-device",
      "type": "device",
      "settings": {
        "mqtt_server":"192.168.1.50",
        "mqtt_user":"",
        "mqtt_pass":"",
        "device:name":"myarduino",
        "device:type":"feather_m0_wifi",
        "device:ssid":"mynetwork",
        "device:wifi_password": "mypass"
      },
      "endpoints": [
        {
          "actionType": "flow",
          "actionURI": "embedded://test",
          "settings": {
            "device:pin": "D:A3",
            "device:condition": "== HIGH",
            "device:response_pin": "D:A4"
          }
        }
      ]
    }
  ]
}

```

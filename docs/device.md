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
## Getting Started
This simple example demonstrates how to create a simple flogo application that has a log activity and device trigger.  The device used for this example is an [Adafruit Feather M0 WiFi](https://learn.adafruit.com/adafruit-feather-m0-wifi-atwinc1500).  It also assumes that you have a digital button attached to pin A3 which will trigger the flow.

- Download flow [myflow.json](https://github.com/TIBCOSoftware/flogo-cli/blob/master/samples/gettingstarted/cli/myflow.json) to build your application. You can also download more samples from the [samples folder](https://github.com/TIBCOSoftware/flogo/tree/master/samples) in the flogo repo. 

```bash
flogo create myDeviceApp
cd myDeviceApp

flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/log
flogo add -b device-trigger trigger github.com/TIBCOSoftware/flogo-contrib/trigger/device
#Make sure myflow.json file under current location
flogo add flow myflow.json

```

- Cd bin folder and open triggers.json in a text editor
- Replace content of triggers.json with the following

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
          "actionURI": "embedded://myflow",
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
- Prepare and build the device
- Build the flogo application
- Upload the device firmware

```bash
flogo device prepare
flogo device build
flogo build

cd devices/tibco-device
flogo device upload
```
- Note: Your MQTT server needs to be running.
- Start flogo engine by running ./myDeviceApp
- Press the button to trigger the flow

For more details about the Device Trigger go [here](https://github.com/TIBCOSoftware/flogo-contrib/tree/device-trigger/trigger/device)

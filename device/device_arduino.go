package device

import (
	"path"

	"github.com/TIBCOSoftware/flogo-cli/util"
)

var triggerContribs = make(map[string]*TriggerContrib)
var actionContribs = make(map[string]*ActionContrib)
var activityContribs = make(map[string]*ActivityContrib)

type TriggerContrib struct {
	Template   string
	Descriptor *TriggerDescriptor
}

func (tc *TriggerContrib) Libs() []*Lib {
	return tc.Descriptor.Libs
}

type ActionContrib struct {
	Ref      string
	Template string
	libs     []*Lib
}

func (ac *ActionContrib) GetActivities(tree *FlowTree) []*ActivityContrib {

	var activities []*ActivityContrib

	for _, task := range tree.AllTasks {
		contrib := activityContribs[task.ActivityRef]
		activities = append(activities, contrib)
	}

	return activities
}

type ActivityContrib struct {
	Template   string
	Descriptor *ActivityDescriptor
}

func (ac *ActivityContrib) Libs() []*Lib {
	return ac.Descriptor.Libs
}

//todo move to contrib contributions
func RegisterActionContrib(ref string, tpl string) *ActionContrib {

	action := &ActionContrib{Ref: ref, Template: tpl}
	actionContribs[action.Ref] = action

	return action
}

func init() {
	//feather_m0 := &DeviceDetails{Type: "feather_m0_wifi", Board: "adafruit_feather_m0_usb"}
	//
	//feather_m0.MainFile = tplArduinoMain
	//
	//files := map[string]string{
	//	"mqtt.ino": tplArduinoMqtt,
	//	"wifi.ino": tplArduinoWifi,
	//}
	//libs := map[string]int{
	//	"PubSubClient": 89,
	//	"WiFi101":      299,
	//}
	//feather_m0.MqttFiles = files
	//feather_m0.Libs = libs

	//Register(feather_m0)

	RegisterActionContrib("github.com/TIBCOSoftware/flogo-contrib/device/action/activity", tplActionActivity)
	RegisterActionContrib("github.com/TIBCOSoftware/flogo-contrib/device/action/flow", tplActionDeviceFlow)
}

func LoadTriggerContrib(proj Project, ref string) (*TriggerContrib, error) {
	contrib, exists := triggerContribs[ref]

	if !exists {
		proj.InstallContribution(ref, "")

		descFile := path.Join(proj.GetContributionDir(), ref, "trigger.json")
		descJson, err := fgutil.LoadLocalFile(descFile)

		if err != nil {
			return nil, err
		}

		desc, err := ParseTriggerDescriptor(descJson)
		//validate that device trigger

		//validate that supports device

		//load template that corresponds for the device

		var details *DeviceSupportDetails

		//fix hardcoded framework, should be based on project
		for _, value := range desc.DeviceSupport {
			if value.Framework == "arduino" {
				details = value
				break
			}
		}

		if details != nil {
			tmplFile := path.Join(proj.GetContributionDir(), ref, details.TemplateFile)
			tmpl, err := fgutil.LoadLocalFile(tmplFile)
			if err != nil {
				return nil, err
			}

			contrib = &TriggerContrib{Descriptor: desc, Template: tmpl}
			triggerContribs[ref] = contrib
		}
	}

	return contrib, nil
}

func LoadActivityContrib(proj Project, ref string) (*ActivityContrib, error) {
	contrib, exists := activityContribs[ref]

	if !exists {
		proj.InstallContribution(ref, "")

		descFile := path.Join(proj.GetContributionDir(), ref, "activity.json")
		actJson, err := fgutil.LoadLocalFile(descFile)

		if err != nil {
			return nil, err
		}

		desc, err := ParseActivityDescriptor(actJson)
		//validate that device trigger

		//validate that supports device

		//load template that corresponds for the device

		var details *DeviceSupportDetails

		//fix hardcoded framework, should be based on project
		for _, value := range desc.DeviceSupport {
			if value.Framework == "arduino" {
				details = value
				break
			}
		}


		if details != nil {
			tmplFile := path.Join(proj.GetContributionDir(), ref, details.TemplateFile)
			tmpl, err := fgutil.LoadLocalFile(tmplFile)
			if err != nil {
				return nil, err
			}

			contrib = &ActivityContrib{Descriptor: desc, Template: tmpl}
			activityContribs[ref] = contrib
		}
		//else return error

	}

	return contrib, nil
}

//var tplArduinoMain = `#include <SPI.h>
//{{if .MqttEnabled}}#include <WiFi101.h>
//#include <PubSubClient.h>
//
//WiFiClient wifiClient;
//PubSubClient client(wifiClient);
//
//{{end}}
//
//
//void setup() {
//    Serial.begin(115200);
//
//    while (!Serial) {
//        delay(10);
//    }
//
//    {{if .MqttEnabled}}
//    setup_wifi();
//    setup_mqtt();
//    {{end}}
//
//	//init triggers
//	{{range .Triggers}}t_{{.}}_init();
//	{{end}}
//
//	//init actions
//	{{range .Actions}}a_{{.}}_init();
//	{{end}}
//}
//
//{{if .MqttEnabled}}
//void init_mqtt_triggers() {
//  //init mqtt triggers
//  {{ range $name, $topic := .MqttTriggers }}t_{{$name}}_init();
//  {{end}}
//}{{end}}
//
//void loop() {
//    {{if .MqttEnabled}}
//    if (!client.connected()) {
//        mqtt_reconnect();
//    }
//
//    // MQTT client loop processing
//    client.loop();
//    {{end}}
//
//	//triggers
//	{{range .Triggers}}t_{{.}}();
//	{{end}}
//}
//
//{{if .MqttEnabled}}
//void callback(char *topic, byte *payload, unsigned int length) {
//
//    Serial.print("Message arrived [");
//    Serial.print(topic);
//    Serial.print("] ");
//    for (int i=0; i < length; i++) {
//        Serial.print((char) payload[i]);
//    }
//    Serial.println();
//
//	//mqtt triggers
//	{{ range $name, $topic := .MqttTriggers }}
//    if (strcmp(topic,"{{$topic}}") == 0) {
//	  t_{{$name}}(topic, payload, length);
//	}{{end}}
//}
//{{end}}
//`
//
//var tplArduinoWifi = `#include <SPI.h>
//#include <WiFi101.h>
//
//char ssid[] = "{{setting . "wifi:ssid"}}";
//const char *password = "{{setting . "wifi:password"}}";
//
//void setup_wifi() {
//
//    //Configure pins for Adafruit ATWINC1500 Feather
//    WiFi.setPins(8,7,4,2);
//
//    // check for the presence of the shield:
//    if (WiFi.status() == WL_NO_SHIELD) {
//        Serial.println("WiFi shield not present");
//        // don't continue:
//        while (true);
//    }
//
//    delay(10);
//
//    Serial.println();
//    Serial.print("Connecting to ");
//    Serial.println(ssid);
//
//    WiFi.begin(ssid, password);
//
//    while (WiFi.status() != WL_CONNECTED) {
//        delay(500);
//        Serial.print(".");
//    }
//
//    //WiFi.lowPowerMode();
//
//    randomSeed(micros());
//
//    Serial.println("");
//    Serial.println("WiFi connected");
//    Serial.println("IP address: ");
//    Serial.println(WiFi.localIP());
//}
//`
//
//var tplArduinoMqtt = `#include <SPI.h>
//#include <WiFi101.h>
//#include <PubSubClient.h>
//
//const char *mqtt_server = "{{setting . "mqtt:server"}}";
//const char *mqtt_user = "{{setting . "mqtt:user"}}";
//const char *mqtt_pass = "{{setting . "mqtt:pass"}}";
//const char *mqtt_pubTopic = "flogo/{{setting . "device:name"}}/out";
//const char *mqtt_subTopic = "flogo/{{setting . "device:name"}}/in";
//
//const char *mqtt_readyMsg = "{\"status\": \"READY\"}";
//
//char out_msg_buff[100];
//
////////////////////////
//
//void setup_mqtt() {
//    client.setServer(mqtt_server, 1883);
//    client.setCallback(callback);
//}
//
//void mqtt_reconnect() {
//    // Loop until we're reconnected
//    while (!client.connected()) {
//        Serial.print("Attempting MQTT connection...");
//        // Create a random client ID
//        String clientId = "device-{{setting . "device:name"}}-";
//        clientId += String(random(0xffff), HEX);
//        // Attempt to connect
//        if (client.connect(clientId.c_str(), mqtt_user, mqtt_pass)) {
//            Serial.println("connected");
//            client.publish(mqtt_pubTopic, mqtt_readyMsg);
//            //client.subscribe(mqtt_subTopic);
//
//            init_mqtt_triggers();
//
//        } else {
//            Serial.print("failed, rc=");
//            Serial.print(client.state());
//            Serial.println(" try again in 5 seconds");
//            // Wait 5 seconds before retrying
//            delay(5000);
//        }
//    }
//}
//
//void publishMQTT(String value, String payload) {
//	payload.toCharArray(out_msg_buff, payload.length() + 1);
//	client.publish(mqtt_pubTopic, out_msg_buff);
//}
//`

var tplActionActivity = `

void a_{{.Id}}_init() {
 	{{ template "activity_init" .Data }}
}

void a_{{.Id}}(int value) {
 	{{ template "activity_code" .Data }}
}
`

var tplActionDeviceFlow = `

void a_{{.Id}}_init() {
	{{range $task := .AllTasks -}}
	ac_{{$task.FlowId}}_{{$task.Id}}_init();
	{{ end }}
}

void a_{{.Id}}(int value) {
 	{{ template "T" .FirstTask }}
}
`

//var tplActionDeviceFlowEval = `
//{{ define "T"}}
//	{{ template "activity_code" .Data }}
//
//	{{if .Condition }}
//	if (value {{.Condition}}) {
//		{{if .Next}}{{with .Next}}{{template "T" .}}{{end}}{{end}}
//	} else {
//		{{if .NextElse}}{{with .NextElse}}{{template "T" .}}{{end}}{{end}}
//	}
//	{{else}}
//		{{if .Next}}{{with .Next}}{{template "T" .}}{{end}}{{end}}
//	{{end}}
//{{ end }}
//`

var tplActionDeviceFlowEval = `{{ define "T" -}}
	{{if .Precondition }}
	if ({{.Precondition}}) {
	    ac_{{.FlowId}}_{{.Id}}(value);
	    {{range $task := .NextTasks}}{{with $task}}{{template "T" .}}{{end}}
	    {{end}}
	}
	{{- else -}}
	ac_{{.FlowId}}_{{.Id}}(value);
	{{range $task := .NextTasks}}{{with $task}}{{template "T" .}}{{end}}
	{{end}}{{end}}{{ end }}
`

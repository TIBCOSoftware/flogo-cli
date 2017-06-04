package device

var triggerContribs = make(map[string]*TriggerContrib)
var actionContribs = make(map[string]*ActionContrib)
var activityContribs = make(map[string]*ActivityContrib)

type Lib struct {
	Name    string
	LibType string
	Ref     string
}
type TriggerContrib struct {
	Ref      string
	Template string
	libs     []*Lib
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
	Ref      string
	Template string
	libs     []*Lib
}

//todo move to contrib contributions
func RegisterTriggerContrib(ref string, tpl string) *TriggerContrib {

	trg := &TriggerContrib{Ref:ref, Template:tpl}
	triggerContribs[trg.Ref] = trg

	return trg
}

//todo move to contrib contributions
func RegisterActionContrib(ref string, tpl string) *ActionContrib {

	action := &ActionContrib{Ref: ref, Template: tpl}
	actionContribs[action.Ref] = action

	return action
}

//todo move to contrib contributions
func RegisterActivityContrib(ref string, tpl string) *ActivityContrib {

	activity := &ActivityContrib{Ref: ref, Template: tpl}
	activityContribs[activity.Ref] = activity

	return activity
}

func init() {
	feather_m0 := &DeviceDetails{Type:"feather_m0_wifi2", Board:"adafruit_feather_m0_usb"}

	feather_m0.MainFile = tplArduinoMainNew

	files := map[string]string{
		"mqtt.ino": tplArduinoMqttNew,
		"wifi.ino": tplArduinoWifiNew,
	}
	libs := map[string]int{
		"PubSubClient": 89,
		"WiFi101": 299,
	}
	feather_m0.MqttFiles = files
	feather_m0.Libs = libs

	Register(feather_m0)

	RegisterTriggerContrib("github.com/TIBCOSoftware/flogo-contrib/device/trigger/pin", tplTriggerPin)
	RegisterTriggerContrib("github.com/TIBCOSoftware/flogo-contrib/device/trigger/pinstream", tplTriggerPinStream)
	RegisterTriggerContrib("github.com/TIBCOSoftware/flogo-contrib/device/trigger/mqtt", tplTriggerMqtt)
	RegisterTriggerContrib("github.com/TIBCOSoftware/flogo-contrib/device/trigger/button", tplTriggerButton)
	bme := RegisterTriggerContrib("github.com/TIBCOSoftware/flogo-contrib/device/trigger/bme280stream", tplTriggerBME280Stream)
	bme.libs = append(bme.libs, &Lib{LibType:"platformio", Ref:"166"})
	bme.libs = append(bme.libs, &Lib{LibType:"platformio", Name:"Adafruit Unified Sensor", Ref:"31"})

	l0x := RegisterTriggerContrib("github.com/TIBCOSoftware/flogo-contrib/device/trigger/vl53l0x_stream", tplTriggerVL53L0XStream)
	l0x.libs = append(bme.libs, &Lib{LibType:"platformio", Ref:"1494"})

	RegisterActionContrib("github.com/TIBCOSoftware/flogo-contrib/device/action/activity", tplActionActivity)
	RegisterActionContrib("github.com/TIBCOSoftware/flogo-contrib/device/action/flow", tplActionDeviceFlow)

	RegisterActivityContrib("github.com/TIBCOSoftware/flogo-contrib/device/activity/setpin", tplActivityPin)
	RegisterActivityContrib("github.com/TIBCOSoftware/flogo-contrib/device/activity/mqtt", tplActivityMqtt)
	RegisterActivityContrib("github.com/TIBCOSoftware/flogo-contrib/device/activity/serial", tplActivitySerial)
}


var tplArduinoMainNew = `#include <SPI.h>
{{if .MqttEnabled}}#include <WiFi101.h>
#include <PubSubClient.h>

WiFiClient wifiClient;
PubSubClient client(wifiClient);

{{end}}


void setup() {
    Serial.begin(115200);

    while (!Serial) {
        delay(10);
    }

    {{if .MqttEnabled}}
    setup_wifi();
    setup_mqtt();
    {{end}}

	//init triggers
	{{range .Triggers}}t_{{.}}_init();
	{{end}}

	//init actions
	{{range .Actions}}a_{{.}}_init();
	{{end}}
}

{{if .MqttEnabled}}
void init_mqtt_triggers() {
  //init mqtt triggers
  {{ range $name, $topic := .MqttTriggers }}t_{{$name}}_init();
  {{end}}
}{{end}}

void loop() {
    {{if .MqttEnabled}}
    if (!client.connected()) {
        mqtt_reconnect();
    }

    // MQTT client loop processing
    client.loop();
    {{end}}

	//triggers
	{{range .Triggers}}t_{{.}}();
	{{end}}
}

{{if .MqttEnabled}}
void callback(char *topic, byte *payload, unsigned int length) {

    Serial.print("Message arrived [");
    Serial.print(topic);
    Serial.print("] ");
    for (int i=0; i < length; i++) {
        Serial.print((char) payload[i]);
    }
    Serial.println();

	//mqtt triggers
	{{ range $name, $topic := .MqttTriggers }}
    if (strcmp(topic,"{{$topic}}") == 0) {
	  t_{{$name}}(topic, payload, length);
	}{{end}}
}
{{end}}
`

var tplArduinoWifiNew = `#include <SPI.h>
#include <WiFi101.h>

char ssid[] = "{{setting . "wifi:ssid"}}";
const char *password = "{{setting . "wifi:password"}}";

void setup_wifi() {

    //Configure pins for Adafruit ATWINC1500 Feather
    WiFi.setPins(8,7,4,2);

    // check for the presence of the shield:
    if (WiFi.status() == WL_NO_SHIELD) {
        Serial.println("WiFi shield not present");
        // don't continue:
        while (true);
    }

    delay(10);

    Serial.println();
    Serial.print("Connecting to ");
    Serial.println(ssid);

    WiFi.begin(ssid, password);

    while (WiFi.status() != WL_CONNECTED) {
        delay(500);
        Serial.print(".");
    }

    //WiFi.lowPowerMode();

    randomSeed(micros());

    Serial.println("");
    Serial.println("WiFi connected");
    Serial.println("IP address: ");
    Serial.println(WiFi.localIP());
}
`

var tplArduinoMqttNew = `#include <SPI.h>
#include <WiFi101.h>
#include <PubSubClient.h>

const char *mqtt_server = "{{setting . "mqtt:server"}}";
const char *mqtt_user = "{{setting . "mqtt:user"}}";
const char *mqtt_pass = "{{setting . "mqtt:pass"}}";
const char *mqtt_pubTopic = "flogo/{{setting . "device:name"}}/out";
const char *mqtt_subTopic = "flogo/{{setting . "device:name"}}/in";

const char *mqtt_readyMsg = "{\"status\": \"READY\"}";

char out_msg_buff[100];

//////////////////////

void setup_mqtt() {
    client.setServer(mqtt_server, 1883);
    client.setCallback(callback);
}

void mqtt_reconnect() {
    // Loop until we're reconnected
    while (!client.connected()) {
        Serial.print("Attempting MQTT connection...");
        // Create a random client ID
        String clientId = "device-{{setting . "device:name"}}-";
        clientId += String(random(0xffff), HEX);
        // Attempt to connect
        if (client.connect(clientId.c_str(), mqtt_user, mqtt_pass)) {
            Serial.println("connected");
            client.publish(mqtt_pubTopic, mqtt_readyMsg);
            //client.subscribe(mqtt_subTopic);

            init_mqtt_triggers();

        } else {
            Serial.print("failed, rc=");
            Serial.print(client.state());
            Serial.println(" try again in 5 seconds");
            // Wait 5 seconds before retrying
            delay(5000);
        }
    }
}

void publishMQTT(String value, String payload) {
	payload.toCharArray(out_msg_buff, payload.length() + 1);
	client.publish(mqtt_pubTopic, out_msg_buff);
}
`

var tplTriggerPin = `
uint8_t t_{{.Id}}_pin = {{setting . "pin"}};    // set input pin
bool t_{{.Id}}_lc = false;    // last condition value

void t_{{.Id}}_init() {
	pinMode(t_{{.Id}}_pin, INPUT);
}

void t_{{.Id}}() {

    int value = {{if settingb . "digital"}}digitalRead(t_{{.Id}}_pin){{else}}analogRead(t_{{.Id}}_pin){{end}};

    // create custom condition
    bool condition = value {{setting . "condition"}};

    if (condition && !t_{{.Id}}_lc) {

        a_{{.ActionId}}(value);
    }

    t_{{.Id}}_lc = condition;
}
`

var tplTriggerPinStream = `
unsigned long t_{{.Id}}_lt = 0; // lastTrigger
uint8_t t_{{.Id}}_pin = {{setting . "pin"}};    // set input pin

void t_{{.Id}}_init() {
	pinMode(t_{{.Id}}_pin, INPUT);
}

void t_{{.Id}}() {

  int value = {{if settingb . "digital"}}digitalRead(t_{{.Id}}_pin){{else}}analogRead(t_{{.Id}}_pin){{end}};

  if ((millis() - t_{{.Id}}_lt) > {{setting . "interval"}}) {
    a_{{.ActionId}}(value);

	t_{{.Id}}_lt = millis();
  }
}
`

//triggers will say if they are part of loop or mqtt callback
var tplTriggerMqtt = `

void t_{{.Id}}_init() {
  client.subscribe("{{setting . "topic"}}");
}

void t_{{.Id}}(char *topic, byte *payload, unsigned int length) {

	char buf[8];
	int i=0;

	for(i=0; i<length; i++) {
		buf[i] = payload[i];
	}
	buf[i] = '\0';

	int value = atoi(buf);

	a_{{.ActionId}}(value);
}
`

var tplTriggerButton = `
unsigned long t_{{.Id}}_ldt = 0; // lastDebounceTime

int t_{{.Id}}_bs;          // the current reading from the input pin
int t_{{.Id}}_lbs = LOW; // the previous reading from the input pin

uint8_t t_{{.Id}}_pin = {{setting . "pin"}};  //set input pin

void t_{{.Id}}_init() {
	pinMode(t_{{.Id}}_pin, INPUT);
}

void t_{{.Id}}() {

  int reading = digitalRead(t_{{.Id}}_pin);

  if (reading != t_{{.Id}}_lbs) {
    // reset the debouncing timer
    t_{{.Id}}_ldt = millis();
  }

  if ((millis() - t_{{.Id}}_ldt) > 50) {

    if (reading != t_{{.Id}}_bs) {
      t_{{.Id}}_bs = reading;

      if (t_{{.Id}}_bs == HIGH) {
        a_{{.ActionId}}(HIGH);
      }
    }
  }

  t_{{.Id}}_lbs = reading;
}
`
var tplTriggerBME280Stream = `
#include <Wire.h>
#include <Adafruit_Sensor.h>
#include <Adafruit_BME280.h>

const int Bme280_cs_pin__i = 5;
Adafruit_BME280 bme(Bme280_cs_pin__i);

unsigned long t_{{.Id}}_lt = 0; // lastTrigger

void t_{{.Id}}_init() {

    bool status;

    status = bme.begin();
    if (!status) {
        Serial.println("Could not find a valid BME280 sensor, check wiring!");
        while (1);
    }
}

void t_{{.Id}}() {

  if ((millis() - t_{{.Id}}_lt) > {{setting . "interval"}}) {
    Serial.println(bme.readTemperature());
    a_{{.ActionId}}(bme.readTemperature());

	t_{{.Id}}_lt = millis();
  }
}
`

var tplTriggerVL53L0XStream = `
#include <Wire.h>
#include "Adafruit_VL53L0X.h"

Adafruit_VL53L0X lox = Adafruit_VL53L0X();

unsigned long t_{{.Id}}_lt = 0; // lastTrigger

void t_{{.Id}}_init() {

    bool status;

  if (!lox.begin()) {
    Serial.println(F("Failed to boot VL53L0X"));
    while(1);
  }
}

void t_{{.Id}}() {

  VL53L0X_RangingMeasurementData_t measure;

  if ((millis() - t_{{.Id}}_lt) > {{setting . "interval"}}) {

    lox.rangingTest(&measure, false); // pass in 'true' to get debug data
    Serial.print("Distance (mm): "); Serial.println(measure.RangeMilliMeter);
    a_{{.ActionId}}(measure.RangeMilliMeter);

	t_{{.Id}}_lt = millis();
  }
}
`

var tplActionActivity = `

void a_{{.Id}}_init() {
 	{{ template "activity_init" .Data }}
}

void a_{{.Id}}(int value) {
 	{{ template "activity_code" .Data }}
}
`

var tplActivityPin = `
void ac_{{.Id}}_{{.Activity.Id}}_init() {
  pinMode({{setting .Activity "pin"}}, OUTPUT);
}

void ac_{{.Id}}_{{.Activity.Id}}(int value) {
  int val = {{if .UseTriggerVal}}value{{else}}{{setting .Activity "value"}}{{end}};
  {{if settingb .Activity "digital"}}digitalWrite({{setting .Activity "pin"}}, val){{else}}analogWrite({{setting .Activity "pin"}}, val){{end}};
}
`

var tplActivityMqtt = `
void ac_{{.Id}}_{{.Activity.Id}}_init() {
}

void ac_{{.Id}}_{{.Activity.Id}}(int value) {
	{{$payload := setting .Activity "payload"}}
	String payload = {{if eq $payload "${value}"}}String(value){{else}}"{{$payload}}"{{end}};
	publishMQTT("{{setting .Activity "topic"}}", payload);
}
`
var tplActivitySerial = `
void ac_{{.Id}}_{{.Activity.Id}}_init() {
}

void ac_{{.Id}}_{{.Activity.Id}}(int value) {
	{{$message := setting .Activity "message"}}
	Serial.println({{if eq $message "${value}"}}String(value){{else}}"{{$message}}"{{end}});
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
